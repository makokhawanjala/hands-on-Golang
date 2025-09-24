package main

import "fmt"

func Phones() {
	// Declare phoneType variable as a string of slices
	phoneType := []string{"iPhone", "Samsung", "Redmi", "Oppo", "Huaweii", "Itel", "Lenovo"}

	// For loop to go through phone types and print each string item
	for _, phone := range phoneType {
		fmt.Println(phone)
	}

}
