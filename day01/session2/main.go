package main

import (
	"fmt"
	"strings"
)

// Accepts a name from the CLI, process it and output it in a clear format
func main() {
	println("Enter your name:")
	var name string
	fmt.Scanln(&name)

	name = strings.TrimSpace(name)
	if name == "" {
		name = "friend"
	}

	fmt.Printf("Hello, %s", name)

}
