package main

import (
	"fmt"
	"strconv"
)

// Name: makeCoffee
// Parameters: coffeeType (string), sugarSpoons (int)
// Body: the instructions

func makeCoffee(coffeeType string, sugarSpoons int) string {
	coffee := "Brewing " + coffeeType
	if sugarSpoons > 0 {
		coffee += " with " + strconv.Itoa(sugarSpoons) + " spoons of sugar"
	}

	return coffee
}

func main() {
	// using the above makeCoffee() function
	result := makeCoffee("espresso", 2)
	fmt.Println(result)

	result1 := makeCoffee("cappucino", 7)
	fmt.Println(result1)

	result2 := makeCoffee("mocha", 4)
	fmt.Println(result2)
}
