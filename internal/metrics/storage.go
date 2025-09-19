package metrics

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gproc/pkg/types"
)

type MetricsStorage struct {
	db *sql.DB
}

func NewMetricsStorage(dbPath string) (*MetricsStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	
	storage := &MetricsStorage{db: db}
	if err := storage.initTables(); err != nil {
		return nil, err
	}
	
	return storage, nil
}

func (m *MetricsStorage) initTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS process_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		process_id TEXT NOT NULL,
		timestamp DATETIME NOT NULL,
		cpu_usage REAL NOT NULL,
		memory_usage INTEGER NOT NULL,
		uptime INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_process_timestamp ON process_metrics(process_id, timestamp);
	`
	
	_, err := m.db.Exec(query)
	return err
}

func (m *MetricsStorage) StoreMetrics(processID string, metrics *types.ProcessMetrics) error {
	query := `
	INSERT INTO process_metrics (process_id, timestamp, cpu_usage, memory_usage, uptime)
	VALUES (?, ?, ?, ?, ?)
	`
	
	_, err := m.db.Exec(query, processID, time.Now(), 
		metrics.CPUUsage, metrics.MemoryUsage, int64(metrics.Uptime))
	return err
}

func (m *MetricsStorage) GetMetricsHistory(processID string, hours int) ([]types.MetricPoint, error) {
	query := `
	SELECT timestamp, cpu_usage, memory_usage 
	FROM process_metrics 
	WHERE process_id = ? AND timestamp > datetime('now', '-' || ? || ' hours')
	ORDER BY timestamp ASC
	`
	
	rows, err := m.db.Query(query, processID, hours)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var points []types.MetricPoint
	for rows.Next() {
		var point types.MetricPoint
		var timestamp string
		
		err := rows.Scan(&timestamp, &point.CPU, &point.Memory)
		if err != nil {
			continue
		}
		
		point.Timestamp, _ = time.Parse("2006-01-02 15:04:05", timestamp)
		points = append(points, point)
	}
	
	return points, nil
}

func (m *MetricsStorage) GetAggregatedMetrics(processID string) (*AggregatedMetrics, error) {
	query := `
	SELECT 
		AVG(cpu_usage) as avg_cpu,
		MAX(cpu_usage) as max_cpu,
		AVG(memory_usage) as avg_memory,
		MAX(memory_usage) as max_memory,
		COUNT(*) as data_points
	FROM process_metrics 
	WHERE process_id = ? AND timestamp > datetime('now', '-24 hours')
	`
	
	var metrics AggregatedMetrics
	err := m.db.QueryRow(query, processID).Scan(
		&metrics.AvgCPU, &metrics.MaxCPU,
		&metrics.AvgMemory, &metrics.MaxMemory,
		&metrics.DataPoints)
	
	return &metrics, err
}

func (m *MetricsStorage) CleanupOldMetrics(days int) error {
	query := `DELETE FROM process_metrics WHERE timestamp < datetime('now', '-' || ? || ' days')`
	_, err := m.db.Exec(query, days)
	return err
}

func (m *MetricsStorage) Close() error {
	return m.db.Close()
}

type AggregatedMetrics struct {
	AvgCPU     float64
	MaxCPU     float64
	AvgMemory  float64
	MaxMemory  float64
	DataPoints int
}