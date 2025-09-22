package enterprise

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gproc/pkg/types"
)

type BackupManager struct {
	config  *types.BackupConfig
	storage StorageProvider
}

type BackupMetadata struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Size        int64     `json:"size"`
	Files       []string  `json:"files"`
	Checksum    string    `json:"checksum"`
	Description string    `json:"description"`
}

type StorageProvider interface {
	Upload(ctx context.Context, key string, data io.Reader) error
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	List(ctx context.Context, prefix string) ([]string, error)
	Delete(ctx context.Context, key string) error
}

func NewBackupManager(config *types.BackupConfig) (*BackupManager, error) {
	if !config.Enabled {
		return &BackupManager{config: config}, nil
	}

	bm := &BackupManager{config: config}

	// Initialize storage provider
	if config.Storage != nil {
		storage, err := NewStorageProvider(config.Storage)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize storage provider: %v", err)
		}
		bm.storage = storage
	} else {
		// Default to local file storage
		bm.storage = &LocalStorage{basePath: "./backups"}
	}

	return bm, nil
}

func (bm *BackupManager) Start(ctx context.Context) error {
	if !bm.config.Enabled {
		return nil
	}

	// Start periodic backup
	go bm.periodicBackup(ctx)

	// Cleanup old backups
	go bm.cleanupOldBackups(ctx)

	fmt.Printf("Backup manager started (interval: %v, retention: %d)\n", 
		bm.config.Interval, bm.config.Retention)
	return nil
}

func (bm *BackupManager) CreateBackup(ctx context.Context, description string) (*BackupMetadata, error) {
	if !bm.config.Enabled {
		return nil, fmt.Errorf("backup not enabled")
	}

	backupID := fmt.Sprintf("backup-%d", time.Now().Unix())
	timestamp := time.Now()

	// Create temporary backup file
	tempFile, err := os.CreateTemp("", "gproc-backup-*.tar.gz")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Create compressed archive
	gzWriter := gzip.NewWriter(tempFile)
	tarWriter := tar.NewWriter(gzWriter)

	// Files to backup
	filesToBackup := []string{
		"gproc.json",     // Process state
		"gproc.yaml",     // Configuration
		"logs/",          // Log directory
		"certs/",         // Certificates (if exists)
	}

	var backedUpFiles []string
	for _, file := range filesToBackup {
		if err := bm.addToArchive(tarWriter, file); err != nil {
			fmt.Printf("Warning: failed to backup %s: %v\n", file, err)
			continue
		}
		backedUpFiles = append(backedUpFiles, file)
	}

	tarWriter.Close()
	gzWriter.Close()

	// Get file size
	stat, err := tempFile.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get backup size: %v", err)
	}

	// Calculate checksum (simplified)
	checksum := fmt.Sprintf("sha256:%x", timestamp.Unix())

	// Create metadata
	metadata := &BackupMetadata{
		ID:          backupID,
		Timestamp:   timestamp,
		Version:     "1.0.0",
		Size:        stat.Size(),
		Files:       backedUpFiles,
		Checksum:    checksum,
		Description: description,
	}

	// Upload backup to storage
	tempFile.Seek(0, 0) // Reset file pointer
	if err := bm.storage.Upload(ctx, backupID+".tar.gz", tempFile); err != nil {
		return nil, fmt.Errorf("failed to upload backup: %v", err)
	}

	// Upload metadata
	metadataJSON, _ := json.Marshal(metadata)
	metadataReader := strings.NewReader(string(metadataJSON))
	if err := bm.storage.Upload(ctx, backupID+".json", metadataReader); err != nil {
		return nil, fmt.Errorf("failed to upload metadata: %v", err)
	}

	fmt.Printf("Backup created: %s (size: %d bytes)\n", backupID, stat.Size())
	return metadata, nil
}

func (bm *BackupManager) RestoreBackup(ctx context.Context, backupID string) error {
	if !bm.config.Enabled {
		return fmt.Errorf("backup not enabled")
	}

	// Download metadata
	metadataReader, err := bm.storage.Download(ctx, backupID+".json")
	if err != nil {
		return fmt.Errorf("failed to download metadata: %v", err)
	}
	defer metadataReader.Close()

	var metadata BackupMetadata
	if err := json.NewDecoder(metadataReader).Decode(&metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %v", err)
	}

	// Download backup file
	backupReader, err := bm.storage.Download(ctx, backupID+".tar.gz")
	if err != nil {
		return fmt.Errorf("failed to download backup: %v", err)
	}
	defer backupReader.Close()

	// Create temporary file for extraction
	tempFile, err := os.CreateTemp("", "gproc-restore-*.tar.gz")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy backup data to temp file
	if _, err := io.Copy(tempFile, backupReader); err != nil {
		return fmt.Errorf("failed to copy backup data: %v", err)
	}

	// Extract backup
	tempFile.Seek(0, 0)
	if err := bm.extractBackup(tempFile); err != nil {
		return fmt.Errorf("failed to extract backup: %v", err)
	}

	fmt.Printf("Backup restored: %s (timestamp: %v)\n", backupID, metadata.Timestamp)
	return nil
}

func (bm *BackupManager) ListBackups(ctx context.Context) ([]*BackupMetadata, error) {
	if !bm.config.Enabled {
		return nil, fmt.Errorf("backup not enabled")
	}

	// List metadata files
	files, err := bm.storage.List(ctx, "backup-")
	if err != nil {
		return nil, fmt.Errorf("failed to list backups: %v", err)
	}

	var backups []*BackupMetadata
	for _, file := range files {
		if !strings.HasSuffix(file, ".json") {
			continue
		}

		// Download and parse metadata
		reader, err := bm.storage.Download(ctx, file)
		if err != nil {
			continue
		}

		var metadata BackupMetadata
		if err := json.NewDecoder(reader).Decode(&metadata); err != nil {
			reader.Close()
			continue
		}
		reader.Close()

		backups = append(backups, &metadata)
	}

	return backups, nil
}

func (bm *BackupManager) DeleteBackup(ctx context.Context, backupID string) error {
	if !bm.config.Enabled {
		return fmt.Errorf("backup not enabled")
	}

	// Delete backup file and metadata
	if err := bm.storage.Delete(ctx, backupID+".tar.gz"); err != nil {
		return fmt.Errorf("failed to delete backup file: %v", err)
	}

	if err := bm.storage.Delete(ctx, backupID+".json"); err != nil {
		return fmt.Errorf("failed to delete metadata: %v", err)
	}

	fmt.Printf("Backup deleted: %s\n", backupID)
	return nil
}

func (bm *BackupManager) periodicBackup(ctx context.Context) {
	ticker := time.NewTicker(bm.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			description := fmt.Sprintf("Automatic backup - %s", time.Now().Format("2006-01-02 15:04:05"))
			if _, err := bm.CreateBackup(ctx, description); err != nil {
				fmt.Printf("Automatic backup failed: %v\n", err)
			}
		}
	}
}

func (bm *BackupManager) cleanupOldBackups(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour) // Check daily
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bm.performCleanup(ctx)
		}
	}
}

func (bm *BackupManager) performCleanup(ctx context.Context) {
	backups, err := bm.ListBackups(ctx)
	if err != nil {
		fmt.Printf("Failed to list backups for cleanup: %v\n", err)
		return
	}

	if len(backups) <= bm.config.Retention {
		return // No cleanup needed
	}

	// Sort backups by timestamp (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.Before(backups[j].Timestamp)
	})

	// Delete oldest backups
	toDelete := len(backups) - bm.config.Retention
	for i := 0; i < toDelete; i++ {
		if err := bm.DeleteBackup(ctx, backups[i].ID); err != nil {
			fmt.Printf("Failed to delete old backup %s: %v\n", backups[i].ID, err)
		}
	}

	fmt.Printf("Cleaned up %d old backups\n", toDelete)
}

func (bm *BackupManager) addToArchive(tarWriter *tar.Writer, path string) error {
	return filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create tar header
		header, err := tar.FileInfoHeader(fi, file)
		if err != nil {
			return err
		}
		header.Name = file

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Write file content if it's a regular file
		if fi.Mode().IsRegular() {
			f, err := os.Open(file)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(tarWriter, f); err != nil {
				return err
			}
		}

		return nil
	})
}

func (bm *BackupManager) extractBackup(reader io.Reader) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Create directory if needed
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(header.Name, 0755); err != nil {
				return err
			}
			continue
		}

		// Create file
		file, err := os.Create(header.Name)
		if err != nil {
			return err
		}

		if _, err := io.Copy(file, tarReader); err != nil {
			file.Close()
			return err
		}
		file.Close()
	}

	return nil
}

// Local Storage Provider
type LocalStorage struct {
	basePath string
}

func (ls *LocalStorage) Upload(ctx context.Context, key string, data io.Reader) error {
	os.MkdirAll(ls.basePath, 0755)
	
	file, err := os.Create(filepath.Join(ls.basePath, key))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, data)
	return err
}

func (ls *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(ls.basePath, key))
}

func (ls *LocalStorage) List(ctx context.Context, prefix string) ([]string, error) {
	var files []string
	
	err := filepath.Walk(ls.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			name := filepath.Base(path)
			if strings.HasPrefix(name, prefix) {
				files = append(files, name)
			}
		}
		return nil
	})
	
	return files, err
}

func (ls *LocalStorage) Delete(ctx context.Context, key string) error {
	return os.Remove(filepath.Join(ls.basePath, key))
}

func NewStorageProvider(config *types.StorageConfig) (StorageProvider, error) {
	switch config.Provider {
	case "local":
		return &LocalStorage{basePath: config.Config["path"]}, nil
	case "s3":
		// Would implement S3 storage provider
		return nil, fmt.Errorf("S3 storage not implemented")
	case "gcs":
		// Would implement Google Cloud Storage provider
		return nil, fmt.Errorf("GCS storage not implemented")
	case "azure":
		// Would implement Azure Blob Storage provider
		return nil, fmt.Errorf("Azure storage not implemented")
	default:
		return nil, fmt.Errorf("unsupported storage provider: %s", config.Provider)
	}
}