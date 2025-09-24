package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// readNonEmptyLine prompts for input, trims it, and gives the user one retry if empty.
// We return the second attempt as-is (even if empty) so that the caller decides on a fallback.
// This separation keeps read logic simple and reuse-friendly
func readNonEmptyLine(prompt string) string {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(prompt + " ")
	text, err := in.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}
	text = strings.TrimSpace(text)
	if text == "" {
		fmt.Println("Oops, that was empty-try once more.")
		fmt.Print(prompt + " ")
		text, err = in.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return ""
		}
		text = strings.TrimSpace(text)
	}
	return text
}

func main() {
	// 1) prompt clearly
	name := readNonEmptyLine("What's your name?")
	// 2) Graceful fallback if still empty after retry
	if name == "" {
		name = "friend"
	}
	// 3) Respond with formatted output
	fmt.Printf("Nice to meet you, %s!\n", name)
	fmt.Println("Welcome to Go-small steps, clean code, daily gains")
}
// CI refresh
