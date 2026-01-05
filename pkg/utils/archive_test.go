package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateArchive(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "archive_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create some files and directories
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.MkdirAll(filepath.Join(sourceDir, "subdir"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sourceDir, "subdir", "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	targetFile := filepath.Join(tmpDir, "test.tar.gz")

	// Run CreateArchive
	if err := CreateArchive(sourceDir, targetFile); err != nil {
		t.Fatalf("CreateArchive failed: %v", err)
	}

	// Verify the archive
	f, err := os.Open(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	gr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	expectedFiles := map[string]bool{
		"file1.txt":        false,
		"subdir/":          false,
		"subdir/file2.txt": false,
	}

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := expectedFiles[header.Name]; ok {
			expectedFiles[header.Name] = true
		} else {
			t.Errorf("Unexpected file in archive: %s", header.Name)
		}
	}

	for name, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file %s not found in archive", name)
		}
	}
}
