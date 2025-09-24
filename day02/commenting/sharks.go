package main

import (
	"fmt"
)

func PrintSharks() {
	// Define sharks variable as a slice of strings
	sharks := []string{"hammerhead", "great white", "dogfish", "frilled"}

	// For loop that iterates over sharks list and prints each string item

	for _, shark := range sharks {
		fmt.Println(shark)
	}

}
