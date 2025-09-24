package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("Welcome to Go Calculator")
	fmt.Println("Type 'exit' anytime to quit.")

	for {
		// First number
		var input string
		fmt.Print("Enter first number (or 'exit'):")
		fmt.Scanln(&input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye")
			break
		}

		num1, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid number, try again.")
			continue
		}

		// second number
		fmt.Print("Enter second number (Or 'exit'):")
		fmt.Scanln(&input)

		if strings.ToLower(input) == "exit" {
			fmt.Println("Goodbye")
			break
		}

		num2, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Println("Invalid number, try again.")
			continue
		}

		// Operator
		fmt.Print("Enter operation (+,-,*,/):")
		var op string
		fmt.Scanln(&op)

		result, err := calculate(num1, num2, op)
		if err != nil {
			fmt.Println("Error", err)
		} else {
			fmt.Printf("Result: %.2f\n", result)
		}
	}
}

func calculate(num1, num2 float64, op string) (float64, error) {
	switch op {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil

	case "*":
		return num1 * num2, nil

	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("cannot divide by zero")
		}
		return num1 / num2, nil
	default:
		return 0, fmt.Errorf("invalid operator: %s", op)
	}

}
