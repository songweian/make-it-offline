package postgresql

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type PostgresqlPlugin struct {
	plugins.BasePlugin
}

func (p *PostgresqlPlugin) GetName() string {
	return "postgresql"
}

func (p *PostgresqlPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating postgresql %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("postgresql-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
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
		case "deb", "apt":
			if err := p.generateDEB(outputPath, appVersion); err != nil {
				return "", err
			}
			commands = append(commands, "dpkg -i *.deb || apt-get install -f -y")
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

func (p *PostgresqlPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3.8'
services:
  db:
    image: postgres:{{.Version}}
    environment:
      POSTGRES_PASSWORD: example_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
volumes:
  postgres_data:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func (p *PostgresqlPlugin) generateDEB(outputPath, appVersion string) error {
	debFile := filepath.Join(outputPath, fmt.Sprintf("postgresql-%s.deb", appVersion))
	return os.WriteFile(debFile, []byte("MOCK DEB CONTENT"), 0644)
}

func init() {
	plugins.Register(&PostgresqlPlugin{})
}
