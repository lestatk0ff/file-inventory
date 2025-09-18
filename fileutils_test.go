package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindFilesVariousCases(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  []string
		setupDirs   []string
		config      Config
		expectCount int
		expectError bool
	}{
		{
			name:        "empty directory",
			setupDirs:   []string{},
			setupFiles:  []string{},
			config:      Config{},
			expectCount: 0,
		},
		{
			name:        "nested structure",
			setupDirs:   []string{"sub1", "sub1/sub2"},
			setupFiles:  []string{"file1.txt", "sub1/file2.txt", "sub1/sub2/file3.txt"},
			config:      Config{},
			expectCount: 3,
		},
		{
			name:        "with hidden files excluded",
			setupFiles:  []string{"visible.txt", ".hidden.txt"},
			config:      Config{IncludeHidden: false},
			expectCount: 1,
		},
		{
			name:        "with hidden files included",
			setupFiles:  []string{"visible.txt", ".hidden.txt"},
			config:      Config{IncludeHidden: true},
			expectCount: 2,
		},
		{
			name:        "with include patterns",
			setupFiles:  []string{"file1.txt", "file2.go", "file3.txt"},
			config:      Config{IncludePatterns: []string{"*.txt"}},
			expectCount: 2,
		},
		{
			name:        "with exclude patterns",
			setupFiles:  []string{"file1.txt", "file2.go", "file3.txt"},
			config:      Config{ExcludePatterns: []string{"*.go"}},
			expectCount: 2,
		},
		{
			name:        "sorted output",
			setupFiles:  []string{"z.txt", "a.txt", "m.txt"},
			config:      Config{SortOutput: true},
			expectCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()

			// Setup directories
			for _, d := range tt.setupDirs {
				os.MkdirAll(filepath.Join(dir, d), 0755)
			}

			// Setup files
			for _, f := range tt.setupFiles {
				fullPath := filepath.Join(dir, f)
				os.MkdirAll(filepath.Dir(fullPath), 0755)
				os.WriteFile(fullPath, []byte("test"), 0644)
			}

			files, err := findFilesWithConfig(dir, tt.config)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(files) != tt.expectCount {
				t.Errorf("Expected %d files, got %d", tt.expectCount, len(files))
			}

			// Test sorting if enabled
			if tt.config.SortOutput && len(files) > 1 {
				for i := 1; i < len(files); i++ {
					if files[i-1] > files[i] {
						t.Error("Files are not sorted")
						break
					}
				}
			}
		})
	}
}

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

func TestFindFilesInvalidDirectory(t *testing.T) {
	_, err := findFiles("/nonexistent/directory")
	if err == nil {
		t.Error("Expected error for nonexistent directory")
	}
}

func TestFindFilesNotADirectory(t *testing.T) {
	// Create a regular file
	tempFile := filepath.Join(t.TempDir(), "notadir.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)

	_, err := findFiles(tempFile)
	if err == nil {
		t.Error("Expected error when path is not a directory")
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
		// Convert relative path to absolute for stat check
		fullPath := f
		if !filepath.IsAbs(f) {
			fullPath = filepath.Join(dir, f)
		}
		info, err := os.Stat(fullPath)
		if err != nil {
			t.Fatalf("Stat failed for %s: %v", fullPath, err)
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

func TestRelativePaths(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	os.Mkdir(subdir, 0755)
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	config := Config{RelativePaths: true}
	files, err := findFilesWithConfig(dir, config)
	if err != nil {
		t.Fatalf("findFilesWithConfig failed: %v", err)
	}

	// All paths should be relative
	for _, f := range files {
		if filepath.IsAbs(f) {
			t.Errorf("Expected relative path, got absolute: %s", f)
		}
	}
}

func TestFullPaths(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	os.Mkdir(subdir, 0755)
	file1 := filepath.Join(dir, "file1.txt")
	file2 := filepath.Join(subdir, "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	config := Config{RelativePaths: false}
	files, err := findFilesWithConfig(dir, config)
	if err != nil {
		t.Fatalf("findFilesWithConfig failed: %v", err)
	}

	// All paths should be absolute
	for _, f := range files {
		if !filepath.IsAbs(f) {
			t.Errorf("Expected absolute path, got relative: %s", f)
		}
	}
}

func TestDefaultBehaviorIsRelative(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "file1.txt")
	os.WriteFile(file1, []byte("test1"), 0644)

	// Default config should use relative paths
	files, err := findFiles(dir)
	if err != nil {
		t.Fatalf("findFiles failed: %v", err)
	}

	// Should have relative paths by default
	for _, f := range files {
		if filepath.IsAbs(f) {
			t.Errorf("Expected relative path by default, got absolute: %s", f)
		}
	}
}
