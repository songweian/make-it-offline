package mattermost

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type MattermostPlugin struct {
	plugins.BasePlugin
}

func (p *MattermostPlugin) GetName() string {
	return "mattermost"
}

func (p *MattermostPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating mattermost %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("mattermost-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
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

func (p *MattermostPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3'
services:
  db:
    image: postgres:13-alpine
    restart: always
    environment:
      POSTGRES_USER: mmuser
      POSTGRES_PASSWORD: mostest_password
      POSTGRES_DB: mattermost
    volumes:
      - db_data:/var/lib/postgresql/data
  mattermost:
    image: mattermost/mattermost-team-edition:{{.Version}}
    restart: always
    depends_on:
      - db
    environment:
      MM_SQLSETTINGS_DRIVERNAME: postgres
      MM_SQLSETTINGS_DATASOURCE: postgres://mmuser:mostest_password@db:5432/mattermost?sslmode=disable&connect_timeout=10
    ports:
      - "8065:8065"
    volumes:
      - mm_config:/mattermost/config
      - mm_data:/mattermost/data
      - mm_logs:/mattermost/logs
volumes:
  db_data:
  mm_config:
  mm_data:
  mm_logs:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func init() {
	plugins.Register(&MattermostPlugin{})
}
