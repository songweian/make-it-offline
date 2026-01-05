package grafana

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type GrafanaPlugin struct {
	plugins.BasePlugin
}

func (p *GrafanaPlugin) GetName() string {
	return "grafana"
}

func (p *GrafanaPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating grafana %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("grafana-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
	outputPath := filepath.Join("output", outputName)

	if err := os.MkdirAll(outputPath, 0755); err != nil {
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

func (p *GrafanaPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3'
services:
  grafana:
    image: grafana/grafana:{{.Version}}
    ports:
      - "3000:3000"
    volumes:
      - grafana_data:/var/lib/grafana
    environment:
      GF_SECURITY_ADMIN_PASSWORD: admin
volumes:
  grafana_data:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func init() {
	plugins.Register(&GrafanaPlugin{})
}
