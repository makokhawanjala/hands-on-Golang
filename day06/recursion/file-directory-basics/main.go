package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("=====Example 1: Reading a Single Directory==========")
	readSingleDirectory("testdata")

	fmt.Println("\n======Example 2: Checking if Path is File or Directory=========")
	checkFileOrDir("testdata/root-file.txt")
	checkFileOrDir("testdata/level1")

	fmt.Println("======= Example 3: Getting File Information======================")
	getFileInfo("testdata")
}

// readSingleDirectory reads only the immediate contents of a Directory
func readSingleDirectory(path string) {
	// open the directory
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	fmt.Printf("Contents of '%s':\n", path)
	for _, entry := range entries {
		// entry.Dir() tells us if it's a directory
		entryType := "FILE"
		if entry.IsDir() {
			entryType = "DIRECTORY"
		}
		fmt.Printf(" [%s] %s\n", entryType, entry.Name())
	}
}

// checkFileOrDir determines if a path is a file or directory
func checkFileOrDir(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if info.IsDir() {
		fmt.Printf("'%s' is a DIRECTORY \n", path)
	} else {
		fmt.Printf("'%s' is a FILE\n", path)
	}
}

// getFileInfo gets detailed information about a file
func getFileInfo(path string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("File: %s\n", info.Name())
	fmt.Printf("Size: %d bytes\n", info.Size())
	fmt.Printf("Permissions: %s\n", info.Mode())
	fmt.Printf("Modified: %s\n", info.ModTime())
	fmt.Printf("Is directory: %t\n", info.IsDir())
}
