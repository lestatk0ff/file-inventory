package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

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

	showDiff(file1, file2)

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

	// The common file should not appear in differences (it's in both files)
	// But it might appear in the table with equal statuses, so we won't test for its absence
}

func TestShowDiffIdenticalFiles(t *testing.T) {
	// Create two identical files
	file1 := "test-identical1.txt"
	file2 := "test-identical2.txt"
	defer os.Remove(file1)
	defer os.Remove(file2)

	content := "file1.txt\nfile2.txt\nfile3.txt\n"
	os.WriteFile(file1, []byte(content), 0644)
	os.WriteFile(file2, []byte(content), 0644)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showDiff(file1, file2)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should show headers but no content rows (no differences)
	if !strings.Contains(output, file1) {
		t.Errorf("Output should contain file1 name in header")
	}
	if !strings.Contains(output, file2) {
		t.Errorf("Output should contain file2 name in header")
	}

	// Count the number of lines - should be minimal (just header and separator)
	lines := strings.Split(strings.TrimSpace(output), "\n")
	// Should have 2 lines: header line and separator line (no data rows for identical files)
	if len(lines) > 3 {
		t.Errorf("Expected minimal output for identical files, got %d lines", len(lines))
	}
}

func TestShowDiffEmptyFiles(t *testing.T) {
	// Create two empty files
	file1 := "test-empty1.txt"
	file2 := "test-empty2.txt"
	defer os.Remove(file1)
	defer os.Remove(file2)

	os.WriteFile(file1, []byte(""), 0644)
	os.WriteFile(file2, []byte(""), 0644)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showDiff(file1, file2)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should show headers but no differences
	if !strings.Contains(output, file1) {
		t.Errorf("Output should contain file1 name in header")
	}
	if !strings.Contains(output, file2) {
		t.Errorf("Output should contain file2 name in header")
	}
}

func TestShowDiffOneEmptyFile(t *testing.T) {
	file1 := "test-empty.txt"
	file2 := "test-nonempty.txt"
	defer os.Remove(file1)
	defer os.Remove(file2)

	os.WriteFile(file1, []byte(""), 0644)
	os.WriteFile(file2, []byte("unique_file.txt\n"), 0644)

	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	showDiff(file1, file2)

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Should show the unique file with proper +/- indicators
	if !strings.Contains(output, "unique_file.txt") {
		t.Errorf("Output should contain unique_file.txt")
	}
}
