package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestShowDiffVariousCases(t *testing.T) {
	tests := []struct {
		name     string
		content1 string
		content2 string
		expectDiffs []string
		expectNoDiffs bool
	}{
		{
			name:     "basic differences",
			content1: "common.txt\nonly_in_file1.txt\n",
			content2: "common.txt\nonly_in_file2.txt\n",
			expectDiffs: []string{"only_in_file1.txt", "only_in_file2.txt"},
		},
		{
			name:     "identical files",
			content1: "file1.txt\nfile2.txt\nfile3.txt\n",
			content2: "file1.txt\nfile2.txt\nfile3.txt\n",
			expectNoDiffs: true,
		},
		{
			name:     "empty files",
			content1: "",
			content2: "",
			expectNoDiffs: true,
		},
		{
			name:     "one empty file",
			content1: "",
			content2: "unique_file.txt\n",
			expectDiffs: []string{"unique_file.txt"},
		},
		{
			name:     "sorted output",
			content1: "z.txt\na.txt\n",
			content2: "m.txt\nb.txt\n",
			expectDiffs: []string{"a.txt", "b.txt", "m.txt", "z.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file1 := "test-diff1.txt"
			file2 := "test-diff2.txt"
			defer os.Remove(file1)
			defer os.Remove(file2)

			os.WriteFile(file1, []byte(tt.content1), 0644)
			os.WriteFile(file2, []byte(tt.content2), 0644)

			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := showDiff(file1, file2)
			if err != nil {
				t.Fatalf("showDiff failed: %v", err)
			}

			w.Close()
			os.Stdout = oldStdout

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			// Check headers are present
			if !strings.Contains(output, file1) {
				t.Errorf("Output should contain file1 name: %s", file1)
			}
			if !strings.Contains(output, file2) {
				t.Errorf("Output should contain file2 name: %s", file2)
			}

			if tt.expectNoDiffs {
				// For identical files, should have minimal output
				lines := strings.Split(strings.TrimSpace(output), "\n")
				if len(lines) > 3 {
					t.Errorf("Expected minimal output for identical files, got %d lines", len(lines))
				}
			} else {
				// Check expected differences are present
				for _, diff := range tt.expectDiffs {
					if !strings.Contains(output, diff) {
						t.Errorf("Output should contain %s", diff)
					}
				}
			}
		})
	}
}

func TestShowDiffBasic(t *testing.T) {
	// Create two temporary files with different content
	file1 := "test-diff1.txt"
	file2 := "test-diff2.txt"
	defer os.Remove(file1)
	defer os.Remove(file2)

	// File1 content
	content1 := "common.txt\nonly_in_file1.txt\n"
	os.WriteFile(file1, []byte(content1), 0644)

	// File2 content
	content2 := "common.txt\nonly_in_file2.txt\n"
	os.WriteFile(file2, []byte(content2), 0644)

	// Capture stdout to test table output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := showDiff(file1, file2)
	if err != nil {
		t.Fatalf("showDiff failed: %v", err)
	}

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check that the output contains the file names in headers
	if !strings.Contains(output, file1) {
		t.Errorf("Output should contain file1 name: %s", file1)
	}
	if !strings.Contains(output, file2) {
		t.Errorf("Output should contain file2 name: %s", file2)
	}

	// Check that differences are shown
	if !strings.Contains(output, "only_in_file1.txt") {
		t.Errorf("Output should contain only_in_file1.txt")
	}
	if !strings.Contains(output, "only_in_file2.txt") {
		t.Errorf("Output should contain only_in_file2.txt")
	}
}

func TestShowDiffNonexistentFile(t *testing.T) {
	file1 := "nonexistent1.txt"
	file2 := "nonexistent2.txt"

	err := showDiff(file1, file2)
	if err == nil {
		t.Error("Expected error for nonexistent files")
	}
}

func TestReadFileLines(t *testing.T) {
	testFile := "test-read-lines.txt"
	defer os.Remove(testFile)

	content := "line1\nline2\n\nline4\n"
	os.WriteFile(testFile, []byte(content), 0644)

	lines, err := readFileLines(testFile)
	if err != nil {
		t.Fatalf("readFileLines failed: %v", err)
	}

	expectedLines := []string{"line1", "line2", "line4"}
	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for _, line := range expectedLines {
		if _, exists := lines[line]; !exists {
			t.Errorf("Expected line %q not found", line)
		}
	}

	// Empty lines should be skipped
	if _, exists := lines[""]; exists {
		t.Error("Empty line should not be included")
	}
}

func TestReadFileLinesNonexistent(t *testing.T) {
	_, err := readFileLines("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}
