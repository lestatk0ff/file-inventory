package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "file-inventory",
		Short: "A tool for creating and comparing file inventories",
		Long:  "file-inventory helps you create file inventories and compare them to track changes in directories.",
	}

	// Global config variables
	var (
		output          string
		sortOutput      bool
		relativePaths   bool
		includeHidden   bool
		excludePatterns []string
		includePatterns []string
	)

	var createCmd = &cobra.Command{
		Use:   "create [DIR]",
		Short: "Create a file inventory for a directory",
		Long:  "Recursively scan a directory and create a text file listing all files found.",
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

	var diffCmd = &cobra.Command{
		Use:   "diff [FILE1] [FILE2]",
		Short: "Show diff between two inventory files",
		Long:  "Compare two inventory files and display differences in a formatted table.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiffCommand(args[0], args[1])
		},
	}

	rootCmd.AddCommand(createCmd, diffCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runCreateCommand(dirPath, output string, config Config) error {
	files, err := findFilesWithConfig(dirPath, config)
	if err != nil {
		return fmt.Errorf("failed to scan directory %q: %w", dirPath, err)
	}

	if err := writeFileList(output, files); err != nil {
		return fmt.Errorf("failed to write inventory to %q: %w", output, err)
	}

	fmt.Printf("Inventory written to %s\n", output)
	fmt.Printf("Total files found: %d\n", len(files))
	return nil
}

func runDiffCommand(file1, file2 string) error {
	if err := showDiff(file1, file2); err != nil {
		return fmt.Errorf("failed to compare files: %w", err)
	}
	return nil
}
