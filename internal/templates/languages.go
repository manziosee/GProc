package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LanguageTemplate struct {
	Name        string
	Extensions  []string
	Command     string
	Args        []string
	HealthCheck string
	EnvVars     map[string]string
	Probes      []string
}

var LanguageTemplates = map[string]LanguageTemplate{
	"node": {
		Name:        "Node.js",
		Extensions:  []string{".js", ".mjs", ".ts"},
		Command:     "node",
		Args:        []string{},
		HealthCheck: "http://localhost:3000/health",
		EnvVars: map[string]string{
			"NODE_ENV": "production",
			"PORT":     "3000",
		},
		Probes: []string{"event_loop_lag", "heap_usage", "memory_rss"},
	},
	"python": {
		Name:        "Python",
		Extensions:  []string{".py"},
		Command:     "python",
		Args:        []string{},
		HealthCheck: "http://localhost:8000/health",
		EnvVars: map[string]string{
			"PYTHONPATH":        ".",
			"PYTHONUNBUFFERED": "1",
		},
		Probes: []string{"gil_wait_percent", "memory_usage", "thread_count"},
	},
	"java": {
		Name:        "Java",
		Extensions:  []string{".jar"},
		Command:     "java",
		Args:        []string{"-jar"},
		HealthCheck: "http://localhost:8080/actuator/health",
		EnvVars: map[string]string{
			"JAVA_OPTS": "-Xmx512m -Xms256m",
		},
		Probes: []string{"heap_usage", "gc_stats", "thread_count"},
	},
	"go": {
		Name:        "Go",
		Extensions:  []string{".go", ""},
		Command:     "./",
		Args:        []string{},
		HealthCheck: "http://localhost:8080/health",
		EnvVars: map[string]string{
			"GOMAXPROCS": "0",
		},
		Probes: []string{"goroutines", "heap_size", "gc_pause"},
	},
	"rust": {
		Name:        "Rust",
		Extensions:  []string{""},
		Command:     "./",
		Args:        []string{},
		HealthCheck: "http://localhost:3000/health",
		EnvVars: map[string]string{
			"RUST_LOG": "info",
		},
		Probes: []string{"thread_count", "memory_usage", "allocator_stats"},
	},
	"php": {
		Name:        "PHP",
		Extensions:  []string{".php"},
		Command:     "php",
		Args:        []string{},
		HealthCheck: "http://localhost:8000/health",
		EnvVars: map[string]string{
			"PHP_INI_SCAN_DIR": "/etc/php/conf.d",
		},
		Probes: []string{"memory_usage", "opcache_stats", "fpm_status"},
	},
}

func DetectLanguage(filePath string) (string, *LanguageTemplate) {
	ext := filepath.Ext(filePath)
	base := filepath.Base(filePath)
	
	for lang, template := range LanguageTemplates {
		for _, validExt := range template.Extensions {
			if ext == validExt {
				return lang, &template
			}
		}
	}
	
	switch {
	case strings.Contains(base, "package.json"):
		template := LanguageTemplates["node"]
		return "node", &template
	case strings.Contains(base, "requirements.txt"):
		template := LanguageTemplates["python"]
		return "python", &template
	case strings.Contains(base, "pom.xml"):
		template := LanguageTemplates["java"]
		return "java", &template
	case strings.Contains(base, "go.mod"):
		template := LanguageTemplates["go"]
		return "go", &template
	case strings.Contains(base, "Cargo.toml"):
		template := LanguageTemplates["rust"]
		return "rust", &template
	case strings.Contains(base, "composer.json"):
		template := LanguageTemplates["php"]
		return "php", &template
	}
	
	return "generic", nil
}

func GenerateConfig(lang, appPath string) (string, error) {
	template, exists := LanguageTemplates[lang]
	if !exists {
		return "", fmt.Errorf("unsupported language: %s", lang)
	}
	
	config := fmt.Sprintf(`# GProc configuration for %s application
processes:
  - name: %s-app
    command: %s
    args: ["%s"]
    working_dir: "%s"
    env:
`, template.Name, lang, template.Command, appPath, filepath.Dir(appPath))
	
	for key, value := range template.EnvVars {
		config += fmt.Sprintf("      %s: %s\n", key, value)
	}
	
	config += fmt.Sprintf(`    auto_restart: true
    max_restarts: 5
    health_check:
      url: "%s"
      interval: "30s"
    probes:
`, template.HealthCheck)
	
	for _, probe := range template.Probes {
		config += fmt.Sprintf("      - %s\n", probe)
	}
	
	return config, nil
}

func InitProject(lang, appPath string) error {
	config, err := GenerateConfig(lang, appPath)
	if err != nil {
		return err
	}
	
	configPath := "gproc.yaml"
	return os.WriteFile(configPath, []byte(config), 0644)
}