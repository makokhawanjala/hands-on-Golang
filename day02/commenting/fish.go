package main

import (
	"fmt"
)

func PrintFish() {
	// Define a fish_type variable as a slice of strings
	fish_type := []string{"Tilapia", "Nile Perch", "Mud Fish"}

	// For loop that iterates over fish type list and prints each string item

	for _, fish := range fish_type {
		fmt.Println(fish)
	}

}
