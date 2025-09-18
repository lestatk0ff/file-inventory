package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCreateCommand(t *testing.T) {
	tests := []struct {
		name        string
		setupFiles  []string
		config      Config
		expectError bool
		expectCount int
	}{
		{
			name:        "basic create",
			setupFiles:  []string{"file1.txt", "file2.txt"},
			config:      Config{},
			expectCount: 2,
		},
		{
			name:        "with sorting",
			setupFiles:  []string{"z.txt", "a.txt", "m.txt"},
			config:      Config{SortOutput: true},
			expectCount: 3,
		},
		{
			name:        "with relative paths",
			setupFiles:  []string{"file1.txt", "sub/file2.txt"},
			config:      Config{RelativePaths: true},
			expectCount: 2,
		},
		{
			name:        "nonexistent directory",
			setupFiles:  []string{},
			config:      Config{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dir string
			if !tt.expectError {
				dir = t.TempDir()
				// Setup files
				for _, f := range tt.setupFiles {
					fullPath := filepath.Join(dir, f)
					os.MkdirAll(filepath.Dir(fullPath), 0755)
					os.WriteFile(fullPath, []byte("test"), 0644)
				}
			} else {
				dir = "/nonexistent/directory"
			}

			output := filepath.Join(t.TempDir(), "test-inventory.txt")
			err := runCreateCommand(dir, output, tt.config)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if !tt.expectError {
				// Verify output file was created and has correct content
				data, err := os.ReadFile(output)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				content := string(data)
				lines := strings.Split(strings.TrimSpace(content), "\n")
				if content == "" {
					lines = []string{}
				}

				if len(lines) != tt.expectCount {
					t.Errorf("Expected %d lines in output, got %d", tt.expectCount, len(lines))
				}

				// Test sorting if enabled
				if tt.config.SortOutput && len(lines) > 1 {
					for i := 1; i < len(lines); i++ {
						if lines[i-1] > lines[i] {
							t.Error("Output is not sorted")
							break
						}
					}
				}
			}
		})
	}
}

func TestRunDiffCommand(t *testing.T) {
	tests := []struct {
		name        string
		content1    string
		content2    string
		expectError bool
	}{
		{
			name:     "valid files",
			content1: "file1.txt\nfile2.txt\n",
			content2: "file1.txt\nfile3.txt\n",
		},
		{
			name:        "nonexistent file1",
			content1:    "",
			content2:    "file1.txt\n",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file1 := "test-diff1.txt"
			file2 := "test-diff2.txt"
			defer os.Remove(file1)
			defer os.Remove(file2)

			if !tt.expectError {
				os.WriteFile(file1, []byte(tt.content1), 0644)
				os.WriteFile(file2, []byte(tt.content2), 0644)
			} else {
				// Only create file2 for nonexistent file1 test
				if tt.content2 != "" {
					os.WriteFile(file2, []byte(tt.content2), 0644)
				}
			}

			err := runDiffCommand(file1, file2)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCobraCreateCommand(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	file1 := filepath.Join(dir, "file1.mp3")
	file2 := filepath.Join(dir, "subdir", "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	// Test the actual cobra command
	var output string
	var sortOutput bool
	var relativePaths bool
	var includeHidden bool
	var excludePatterns []string
	var includePatterns []string

	createCmd := &cobra.Command{
		Use:   "create [DIR]",
		Short: "Create a file inventory for a directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateCommand(args[0], output, Config{
				SortOutput:      sortOutput,
				RelativePaths:   relativePaths,
				IncludeHidden:   includeHidden,
				ExcludePatterns: excludePatterns,
				IncludePatterns: includePatterns,
			})
		},
	}

	createCmd.Flags().StringVarP(&output, "output", "o", "file-inventory.txt", "Output file name")
	createCmd.Flags().BoolVar(&sortOutput, "sort", false, "Sort file paths in output")
	createCmd.Flags().BoolVar(&relativePaths, "relative", false, "Use relative paths from scan directory")
	createCmd.Flags().BoolVar(&includeHidden, "hidden", false, "Include hidden files and directories")
	createCmd.Flags().StringSliceVar(&excludePatterns, "exclude", []string{}, "Exclude patterns (glob)")
	createCmd.Flags().StringSliceVar(&includePatterns, "include", []string{}, "Include patterns (glob)")

	output = "test-inventory.txt"
	createCmd.SetArgs([]string{dir, "-o", output})

	if err := createCmd.Execute(); err != nil {
		t.Fatalf("Cobra create command failed: %v", err)
	}

	defer os.Remove(output)
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	outputStr := string(data)
	if !strings.Contains(outputStr, "file1.mp3") || !strings.Contains(outputStr, "file2.txt") {
		t.Errorf("Output missing expected files: %s", outputStr)
	}
}

func TestCobraDiffCommand(t *testing.T) {
	// Create two inventory files with a known diff
	fileA := "a.txt"
	fileB := "b.txt"
	os.WriteFile(fileA, []byte("foo\nbar\nuniqueA\n"), 0644)
	os.WriteFile(fileB, []byte("foo\nbar\nuniqueB\n"), 0644)
	defer os.Remove(fileA)
	defer os.Remove(fileB)

	diffCmd := &cobra.Command{
		Use:   "diff [FILE1] [FILE2]",
		Short: "Show diff between two inventory files",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiffCommand(args[0], args[1])
		},
	}

	diffCmd.SetArgs([]string{fileA, fileB})
	if err := diffCmd.Execute(); err != nil {
		t.Fatalf("Cobra diff command failed: %v", err)
	}
}
