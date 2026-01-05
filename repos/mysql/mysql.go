package mysql

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type MysqlPlugin struct {
	plugins.BasePlugin
}

func (p *MysqlPlugin) GetName() string {
	return "mysql"
}

func (p *MysqlPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating mysql %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("mysql-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
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
		case "rpm", "yum":
			if err := p.generateRPM(outputPath, appVersion); err != nil {
				return "", err
			}
			commands = append(commands, "yum localinstall -y *.rpm")
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

func (p *MysqlPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3.8'
services:
  db:
    image: mysql:{{.Version}}
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: example_db
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
volumes:
  mysql_data:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func (p *MysqlPlugin) generateRPM(outputPath, appVersion string) error {
	rpmFile := filepath.Join(outputPath, fmt.Sprintf("mysql-%s.rpm", appVersion))
	return os.WriteFile(rpmFile, []byte("MOCK RPM CONTENT"), 0644)
}

func init() {
	plugins.Register(&MysqlPlugin{})
}
