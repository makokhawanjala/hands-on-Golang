package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// walk recursively reads the directory at path and prints its entries.
// indent is a string used to visually indent child entries.
func walk(path string, indent string) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		// Return the error to the caller so it can decide what to do.
		return err
	}

	for _, entry := range entries {
		name := entry.Name()
		fmt.Println(indent + name)

		// If this entry is a directory, call walk recursively on it.
		if entry.IsDir() {
			subPath := filepath.Join(path, name)
			if err := walk(subPath, indent+"  "); err != nil {
				// Print the error and continue with the next entry.
				fmt.Fprintf(os.Stderr, "error reading %s: %v\n", subPath, err)
			}
		}
	}

	return nil
}

func main() {
	// Parse flags (so we can accept an optional path argument).
	flag.Parse()

	// Default root path is current directory.
	root := "."
	if flag.NArg() > 0 {
		// If the user passed an argument, use it as the root.
		root = flag.Arg(0)
	}

	// Convert the root path to an absolute path for clarity when printed.
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	// Print the root path we're about to walk.
	fmt.Println(absRoot)

	// Start the recursive walk. If walk returns an error, print it and exit.
	if err := walk(absRoot, ""); err != nil {
		fmt.Fprintf(os.Stderr, "walk error: %v\n", err)
		os.Exit(1)
	}
}
