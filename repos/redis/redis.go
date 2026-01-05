package redis

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"path/filepath"
)

type RedisPlugin struct {
	plugins.BasePlugin
}

func (p *RedisPlugin) GetName() string {
	return "redis"
}

func (p *RedisPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating redis %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("redis-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
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

func (p *RedisPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3.8'
services:
  redis:
    image: redis:{{.Version}}
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
volumes:
  redis_data:
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func (p *RedisPlugin) generateRPM(outputPath, appVersion string) error {
	rpmFile := filepath.Join(outputPath, fmt.Sprintf("redis-%s.rpm", appVersion))
	return os.WriteFile(rpmFile, []byte("MOCK RPM CONTENT"), 0644)
}

func init() {
	plugins.Register(&RedisPlugin{})
}
