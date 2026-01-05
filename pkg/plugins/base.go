package plugins

import (
	"gopkg.in/yaml.v3"
	"make-it-offline/pkg/utils"
	"os"
	"path/filepath"
	"text/template"
)

type Plugin interface {
	GetName() string
	Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error)
}

type BasePlugin struct{}

func (b *BasePlugin) WriteDockerCompose(outputPath string, composeContent interface{}) error {
	data, err := yaml.Marshal(composeContent)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(outputPath, "docker-compose.yml"), data, 0644)
}

func (b *BasePlugin) WriteDockerComposeWithTemplate(outputPath string, templateStr string, data interface{}) error {
	tmpl, err := template.New("docker-compose").Parse(templateStr)
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(outputPath, "docker-compose.yml"))
	if err != nil {
		return err
	}
	defer f.Close()
	return tmpl.Execute(f, data)
}

func (b *BasePlugin) WriteInstallScript(outputPath string, commands []string) error {
	scriptPath := filepath.Join(outputPath, "install.sh")
	content := "#!/bin/bash\nset -e\n"
	for _, cmd := range commands {
		content += cmd + "\n"
	}
	err := os.WriteFile(scriptPath, []byte(content), 0755)
	if err != nil {
		return err
	}
	return os.Chmod(scriptPath, 0755)
}

func (b *BasePlugin) CreateArchive(source, target string) error {
	return utils.CreateArchive(source, target)
}

var registry = make(map[string]Plugin)

func Register(p Plugin) {
	registry[p.GetName()] = p
}

func GetPlugin(name string) (Plugin, bool) {
	p, ok := registry[name]
	return p, ok
}
