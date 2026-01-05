package plugins

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBasePlugin_WriteInstallScript(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "plugin_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	p := &BasePlugin{}
	commands := []string{"echo hello", "ls -l"}
	err = p.WriteInstallScript(tmpDir, commands)
	if err != nil {
		t.Fatalf("WriteInstallScript failed: %v", err)
	}

	scriptPath := filepath.Join(tmpDir, "install.sh")
	content, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatal(err)
	}

	expected := "#!/bin/bash\nset -e\necho hello\nls -l\n"
	if string(content) != expected {
		t.Errorf("Expected content %q, got %q", expected, string(content))
	}

	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0111 == 0 {
		t.Error("Expected script to be executable")
	}
}

func TestBasePlugin_WriteDockerComposeWithTemplate(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "plugin_test_docker")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	p := &BasePlugin{}
	tmpl := "image: nginx:{{.Version}}"
	data := struct{ Version string }{Version: "1.21"}

	err = p.WriteDockerComposeWithTemplate(tmpDir, tmpl, data)
	if err != nil {
		t.Fatalf("WriteDockerComposeWithTemplate failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(tmpDir, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}

	expected := "image: nginx:1.21"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
	}
}
