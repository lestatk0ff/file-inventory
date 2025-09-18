package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Config holds configuration options for file operations
type Config struct {
	SortOutput    bool
	RelativePaths bool
	IncludeHidden bool
	ExcludePatterns []string
	IncludePatterns []string
}

func writeFileList(filename string, files []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for _, file := range files {
		if _, err := fmt.Fprintln(writer, file); err != nil {
			return fmt.Errorf("failed to write file entry: %w", err)
		}
	}
	return nil
}

func findFiles(dirPath string) ([]string, error) {
	return findFilesWithConfig(dirPath, Config{})
}

func findFilesWithConfig(dirPath string, config Config) ([]string, error) {
	// Validate input directory
	if info, err := os.Stat(dirPath); err != nil {
		return nil, fmt.Errorf("cannot access directory: %w", err)
	} else if !info.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", dirPath)
	}

	var files []string
	var count int

	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Log warning but continue processing
			fmt.Fprintf(os.Stderr, "Warning: skipping %q: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		// Skip hidden files if not included
		if !config.IncludeHidden && isHidden(path) {
			return nil
		}

		// Apply include/exclude patterns
		if !shouldIncludeFile(path, config) {
			return nil
		}

		// Convert to relative path if requested
		finalPath := path
		if config.RelativePaths {
			if relPath, err := filepath.Rel(dirPath, path); err == nil {
				finalPath = relPath
			}
		}

		files = append(files, finalPath)
		count++

		// Show progress for large directories
		if count%1000 == 0 {
			fmt.Fprintf(os.Stderr, "Found %d files...\r", count)
		}

		return nil
	})

	if count > 0 && count%1000 == 0 {
		fmt.Fprintf(os.Stderr, "\n")
	}

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %w", err)
	}

	// Sort output if requested
	if config.SortOutput {
		sort.Strings(files)
	}

	return files, nil
}

func isHidden(path string) bool {
	base := filepath.Base(path)
	return strings.HasPrefix(base, ".")
}

func shouldIncludeFile(path string, config Config) bool {
	// If include patterns are specified, file must match at least one
	if len(config.IncludePatterns) > 0 {
		matched := false
		for _, pattern := range config.IncludePatterns {
			if match, _ := filepath.Match(pattern, filepath.Base(path)); match {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// If exclude patterns are specified, file must not match any
	for _, pattern := range config.ExcludePatterns {
		if match, _ := filepath.Match(pattern, filepath.Base(path)); match {
			return false
		}
	}

	return true
}
