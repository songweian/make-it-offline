package prometheus

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type PrometheusPlugin struct {
	plugins.BasePlugin
}

func (p *PrometheusPlugin) GetName() string {
	return "prometheus"
}

func (p *PrometheusPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating prometheus %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("prometheus-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
	outputPath := filepath.Join("output", outputName)

	if err := os.MkdirAll(outputPath, 0755); err != nil {
		return "", err
	}

	if err := os.WriteFile(filepath.Join(outputPath, "prometheus.yml"), []byte("global:\n  scrape_interval: 15s\nscrape_configs:\n  - job_name: 'prometheus'\n    static_configs:\n      - targets: ['localhost:9090']"), 0644); err != nil {
		return "", err
	}

	var commands []string
	for _, format := range formats {
		switch format {
		case "docker-compose":
			if err := p.generateDockerCompose(outputPath, appVersion); err != nil {
				return "", err
			}
			commands = append(commands, "docker-compose up -d")
		}
	}

	if err := p.WriteInstallScript(outputPath, commands); err != nil {
		return "", err
	}

	archivePath := outputPath + ".tar.gz"
	if err := p.CreateArchive(outputPath, archivePath); err != nil {
		return "", err
	}

	return archivePath, nil
}

func (p *PrometheusPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3'
services:
  prometheus:
    image: prom/prometheus:v{{.Version}}
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus
    command: --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus
volumes:
  prometheus_data:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func init() {
	plugins.Register(&PrometheusPlugin{})
}
