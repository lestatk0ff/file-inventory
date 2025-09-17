package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

// showDiff compares two files and prints lines unique to each in table format
func showDiff(file1, file2 string) {
	set1 := make(map[string]struct{})
	set2 := make(map[string]struct{})

	f1, err := os.Open(file1)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", file1, err)
		os.Exit(1)
	}
	defer f1.Close()
	f2, err := os.Open(file2)
	if err != nil {
		fmt.Printf("Error opening %s: %v\n", file2, err)
		os.Exit(1)
	}
	defer f2.Close()

	scanner1 := bufio.NewScanner(f1)
	for scanner1.Scan() {
		set1[scanner1.Text()] = struct{}{}
	}
	scanner2 := bufio.NewScanner(f2)
	for scanner2.Scan() {
		set2[scanner2.Text()] = struct{}{}
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
	for filePath := range allFiles {
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
}
