package main

import (
	"fmt"
	"strings"
)

// Accepts a name from the CLI, process it and output it in a clear format
func main() {
	fmt.Println("Enter your name:") // Use fmt.Println instead of println for consistency
	var name string

	if _, err := fmt.Scanln(&name); err != nil {
		fmt.Println("Error reading input:", err)
		return
	}

	name = strings.TrimSpace(name)
	if name == "" {
		name = "friend"
	}

	fmt.Printf("Hello, %s\n", name) // Added newline for better output formatting
}

