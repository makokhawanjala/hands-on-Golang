package main

import (
	"fmt"
	"strings"
)

// Function to accept name from the command line
// and print the entered name with a greeting
func main() {
	//print this in the CLI
	fmt.Println("Please enter your name.")

	//define a variable name which is a string to store the name
	var name string
	//Accept user input from the command line
	fmt.Scanln(&name)
	// Use TrimSpace function from Go's std library
	// to remove any extra space characters
	name = strings.TrimSpace(name)

	// Now output the name and the greeting
	fmt.Printf("Hi, %s! I'm Go!", name)
}
