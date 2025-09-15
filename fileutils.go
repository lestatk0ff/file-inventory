package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func writeFileList(filename string, files []string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, file := range files {
		if _, err := fmt.Fprintln(f, file); err != nil {
			return err
		}
	}
	return nil
}

func findFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
