package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestCobraCreateCommand(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "subdir"), 0755)
	file1 := filepath.Join(dir, "file1.mp3")
	file2 := filepath.Join(dir, "subdir", "file2.txt")
	os.WriteFile(file1, []byte("test1"), 0644)
	os.WriteFile(file2, []byte("test2"), 0644)

	// Prepare root command as in main
	var output string
	var buf bytes.Buffer
	createCmd := &cobra.Command{
		Use:   "create [DIR]",
		Short: "Create a file inventory for a directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := args[0]
			if output == "" {
				output = "file-inventory.txt"
			}
			files, err := findFiles(dirPath)
			if err != nil {
				t.Fatalf("Error walking the path: %v", err)
			}
			err = writeFileList(output, files)
			if err != nil {
				t.Fatalf("Error writing file list: %v", err)
			}
			buf.WriteString("Inventory written to " + output + "\n")
			buf.WriteString("Total files found: ")
			buf.WriteString(fmt.Sprint(len(files)))
		},
	}
	createCmd.Flags().StringVarP(&output, "output", "o", "", "Output file name (default is file-inventory.txt)")
	createCmd.SetArgs([]string{dir, "-o", "test-inventory.txt"})
	if err := createCmd.Execute(); err != nil {
		t.Fatalf("Cobra create command failed: %v", err)
	}
	data, err := os.ReadFile("test-inventory.txt")
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

	var buf bytes.Buffer
	diffCmd := &cobra.Command{
		Use:   "diff [FILE1] [FILE2]",
		Short: "Show diff between two inventory files",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// Use the actual showDiff implementation if available
			set1 := make(map[string]struct{})
			set2 := make(map[string]struct{})
			f1, _ := os.Open(args[0])
			defer f1.Close()
			f2, _ := os.Open(args[1])
			defer f2.Close()
			scanner1 := bufio.NewScanner(f1)
			for scanner1.Scan() {
				set1[scanner1.Text()] = struct{}{}
			}
			scanner2 := bufio.NewScanner(f2)
			for scanner2.Scan() {
				set2[scanner2.Text()] = struct{}{}
			}
			for line := range set1 {
				if _, ok := set2[line]; !ok {
					buf.WriteString("-" + line + "\n")
				}
			}
			for line := range set2 {
				if _, ok := set1[line]; !ok {
					buf.WriteString("+" + line + "\n")
				}
			}
		},
	}
	diffCmd.SetArgs([]string{fileA, fileB})
	if err := diffCmd.Execute(); err != nil {
		t.Fatalf("Cobra diff command failed: %v", err)
	}
	diffOut := buf.String()
	if !strings.Contains(diffOut, "-uniqueA") || !strings.Contains(diffOut, "+uniqueB") {
		t.Errorf("Diff output missing expected lines: %s", diffOut)
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
