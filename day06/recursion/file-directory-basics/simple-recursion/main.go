package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	fmt.Println("=========Simple Recursive Directory Walker=========\n")

	err := walkDirectory("/mnt/f/dev-workspace", 0)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// walkDirectory recursively walks through directories
// depth parameter helps us visualize the tree structure
func walkDirectory(path string, depth int) error {
	// Step 1: Read the current directory
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	// Step 2: Process each entry
	for _, entry := range entries {
		// Create indentation based on depth
		indent := ""
		for i := 0; i < depth; i++ {
			indent += " "
		}

		// Build the full path
		fullPath := filepath.Join(path, entry.Name())

		//Print current item
		if entry.IsDir() {
			fmt.Printf("%s[DIR] %s\n", indent, entry.Name())

			// Step 3: THE RECURSSION - if it's a directory, go into it!
			err := walkDirectory(fullPath, depth+1)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("%s[FILE] %s\n", indent, entry.Name())

		}
	}
	return nil
}
