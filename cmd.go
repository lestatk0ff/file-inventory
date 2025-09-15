package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{Use: "file-inventory"}

	var output string
	var createCmd = &cobra.Command{
		Use:   "create [DIR]",
		Short: "Create a file inventory for a directory",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dirPath := args[0]
			files, err := findFiles(dirPath)
			if err != nil {
				fmt.Printf("Error walking the path: %v\n", err)
				os.Exit(1)
			}
			err = writeFileList(output, files)
			if err != nil {
				fmt.Printf("Error writing file list: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Inventory written to %s\n", output)
			fmt.Printf("Total files found: %d\n", len(files))
		},
	}
	createCmd.Flags().StringVarP(&output, "output", "o", "file-inventory", "Output file name (default is file-inventory.txt)")

	var diffCmd = &cobra.Command{
		Use:   "diff [FILE1] [FILE2]",
		Short: "Show diff between two inventory files",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			showDiff(args[0], args[1])
		},
	}

	rootCmd.AddCommand(createCmd, diffCmd)
	rootCmd.Execute()
}
