package main

import (
	"bufio"
	"fmt"
	"os"
)

// showDiff compares two files and prints lines unique to each
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

	fmt.Println("Lines in", file1, "but not in", file2, ":")
	for line := range set1 {
		if _, ok := set2[line]; !ok {
			fmt.Println("-", line)
		}
	}
	fmt.Println("Lines in", file2, "but not in", file1, ":")
	for line := range set2 {
		if _, ok := set1[line]; !ok {
			fmt.Println("+", line)
		}
	}
}
