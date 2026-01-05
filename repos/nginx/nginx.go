package nginx

import (
	"fmt"
	"make-it-offline/pkg/plugins"
	"os"
	"os/exec"
	"path/filepath"
)

type NginxPlugin struct {
	plugins.BasePlugin
}

func (p *NginxPlugin) GetName() string {
	return "nginx"
}

func (p *NginxPlugin) Generate(appVersion, osName, osVersion, arch string, formats []string) (string, error) {
	fmt.Printf("Generating nginx %s for %s %s (%s) in formats: %v\n", appVersion, osName, osVersion, arch, formats)

	outputName := fmt.Sprintf("nginx-%s-%s-%s-%s", appVersion, osName, osVersion, arch)
	outputPath, _ := filepath.Abs(filepath.Join("output", outputName))

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
			if err := p.generateRPM(outputPath, appVersion, osName, osVersion, arch); err != nil {
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

func (p *NginxPlugin) generateDockerCompose(outputPath, appVersion string) error {
	tmpl := `version: '3'
services:
  nginx:
    image: nginx:{{.Version}}
    ports:
      - "80:80"
`
	data := struct {
		Version string
	}{
		Version: appVersion,
	}
	return p.WriteDockerComposeWithTemplate(outputPath, tmpl, data)
}

func (p *NginxPlugin) getDockerCommand() string {
	if _, err := exec.LookPath("docker"); err == nil {
		return "docker"
	}
	if _, err := exec.LookPath("podman"); err == nil {
		return "podman"
	}
	return "docker" // fallback
}

func (p *NginxPlugin) generateRPM(outputPath, appVersion, osName, osVersion, arch string) error {
	dockerCmd := p.getDockerCommand()
	dockerFile := filepath.Join("docker", "os", fmt.Sprintf("%s-%s-%s", osName, osVersion, arch))
	if _, err := os.Stat(dockerFile); os.IsNotExist(err) {
		// 尝试匹配 redhat 到 readhat (处理拼写错误)
		if osName == "redhat" {
			dockerFile = filepath.Join("docker", "os", fmt.Sprintf("readhat-%s-%s", osVersion, arch))
		}
	}

	if _, err := os.Stat(dockerFile); os.IsNotExist(err) {
		return fmt.Errorf("dockerfile not found: %s", dockerFile)
	}

	imageName := fmt.Sprintf("make-it-offline/%s-%s-%s:latest", osName, osVersion, arch)
	buildCmd := exec.Command(dockerCmd, "build", "-t", imageName, "-f", dockerFile, filepath.Dir(dockerFile))
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("failed to build docker image: %v", err)
	}

	// 使用 yumdownloader 下载 nginx
	// 需要先安装 yum-utils 和 nginx repo
	// 尝试下载指定版本，如果失败则下载最新版本
	downloadCmd := fmt.Sprintf("mkdir -p /tmp/download && (yum install -y epel-release || true) && yum install -y yum-utils && (yumdownloader --destdir=/tmp/download --resolve nginx-%s || yumdownloader --destdir=/tmp/download --resolve nginx)", appVersion)

	runCmd := exec.Command(dockerCmd, "run", "--rm",
		"-v", fmt.Sprintf("%s:/output", outputPath),
		imageName,
		"sh", "-c", fmt.Sprintf("%s && cp /tmp/download/*.rpm /output/", downloadCmd))

	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr

	if err := runCmd.Run(); err != nil {
		return fmt.Errorf("failed to download RPM in docker: %v", err)
	}

	return nil
}

func init() {
	plugins.Register(&NginxPlugin{})
}
