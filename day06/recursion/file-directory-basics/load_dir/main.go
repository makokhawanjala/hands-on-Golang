package main

import (
	"fmt"
	"os"
)

func main() {
	readSingleDir("/mnt/f/dev-workspace/social_website/my_bookmarks_env/Scripts")
}

func readSingleDir(path string) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf(" %v\n", err)
		return
	}

	for _, entry := range entries {
		entryType := "FILE"

		if entry.IsDir() {
			entryType = "DIRECTORY"
		}

		fmt.Printf("%v : %v\n", entry, entryType)

	}
}
