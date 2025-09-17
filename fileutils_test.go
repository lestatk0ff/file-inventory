package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindFiles(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	os.Mkdir(subdir, 0755)
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	files, err := findFiles(dir)
	if err != nil {
		t.Fatalf("findFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(files))
	}

	// Check that both files are found
	foundFile1 := false
	foundFile2 := false
	for _, f := range files {
		if strings.HasSuffix(f, "file1.txt") {
			foundFile1 = true
		}
		if strings.HasSuffix(f, "file2.txt") {
			foundFile2 = true
		}
	}

	if !foundFile1 {
		t.Error("file1.txt not found in results")
	}
	if !foundFile2 {
		t.Error("file2.txt not found in results")
	}
}

func TestInventoryDoesNotIncludeDirectories(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	os.Mkdir(subdir, 0755)
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	files, err := findFiles(dir)
	if err != nil {
		t.Fatalf("findFiles failed: %v", err)
	}
	for _, f := range files {
		info, err := os.Stat(f)
		if err != nil {
			t.Fatalf("Stat failed for %s: %v", f, err)
		}
		if info.IsDir() {
			t.Errorf("Directory %s found in file list!", f)
		}
	}
}

func TestWriteFileList(t *testing.T) {
	testFile := "test-write-filelist.txt"
	defer os.Remove(testFile)

	testFiles := []string{
		"/path/to/file1.txt",
		"/path/to/file2.mp3",
		"/another/path/file3.doc",
	}

	err := writeFileList(testFile, testFiles)
	if err != nil {
		t.Fatalf("writeFileList failed: %v", err)
	}

	// Read the file back and verify contents
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	content := string(data)
	for _, file := range testFiles {
		if !strings.Contains(content, file) {
			t.Errorf("File %s not found in output", file)
		}
	}

	// Check that we have the right number of lines
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != len(testFiles) {
		t.Errorf("Expected %d lines, got %d", len(testFiles), len(lines))
	}
}

func TestWriteFileListEmptySlice(t *testing.T) {
	testFile := "test-write-empty.txt"
	defer os.Remove(testFile)

	err := writeFileList(testFile, []string{})
	if err != nil {
		t.Fatalf("writeFileList with empty slice failed: %v", err)
	}

	// Check that file exists but is empty
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if len(data) != 0 {
		t.Errorf("Expected empty file, got %d bytes", len(data))
	}
}
