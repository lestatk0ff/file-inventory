package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// showDiff compares two files and prints lines unique to each in table format
func showDiff(file1, file2 string) error {
	set1, err := readFileLines(file1)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", file1, err)
	}

	set2, err := readFileLines(file2)
	if err != nil {
		return fmt.Errorf("error reading %s: %w", file2, err)
	}

	// Collect all unique files
	allFiles := make(map[string]struct{})
	for line := range set1 {
		if _, ok := set2[line]; !ok {
			allFiles[line] = struct{}{}
		}
	}
	for line := range set2 {
		if _, ok := set1[line]; !ok {
			allFiles[line] = struct{}{}
		}
	}

	// Sort files for consistent output
	var sortedFiles []string
	for filePath := range allFiles {
		sortedFiles = append(sortedFiles, filePath)
	}
	sort.Strings(sortedFiles)

	// Create table
	table := tablewriter.NewWriter(os.Stdout)
	table.Options(
		tablewriter.WithHeaderAutoFormat(tw.Off),
		tablewriter.WithRendition(tw.Rendition{
			Borders: tw.Border{
				Left:   tw.Off,
				Right:  tw.Off,
				Top:    tw.Off,
				Bottom: tw.Off,
			},
		}),
	)
	table.Header("file_path", file1, file2)

	// Add differences to table
	for _, filePath := range sortedFiles {
		_, inFile1 := set1[filePath]
		_, inFile2 := set2[filePath]

		file1Status := " "
		file2Status := " "

		if inFile1 && !inFile2 {
			file1Status = "+"
			file2Status = "-"
		} else if !inFile1 && inFile2 {
			file1Status = "-"
			file2Status = "+"
		}

		table.Append(filePath, file1Status, file2Status)
	}

	table.Render()
	return nil
}

// readFileLines reads a file and returns a set of its lines
func readFileLines(filename string) (map[string]struct{}, error) {
	set := make(map[string]struct{})

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" { // Skip empty lines
			set[line] = struct{}{}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return set, nil
}
