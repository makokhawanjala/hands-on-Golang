package main

import (
	"errors"
	"fmt"
)

// processPayment simulates a payment transaction. takes amount and cardNumber and returns an error
func processPayment(amount float64, cardNumber string) error {
	// validation checks
	if amount <= 0 {
		return errors.New("payment amount should be be greater than zero")
	}

	if len(cardNumber) < 16 {
		return errors.New("your card number is too short")
	}
	fmt.Printf("Processing Payment of $%.2f\n", amount)
	return nil
}

func main() {
	fmt.Println("===================E-commerce Payment System=========================")
	// Example 1: Valid transaction
	fmt.Println("Valid Payment")
	err := processPayment(99.99, "1234567890123456")
	if err != nil {
		// This block runs where there is a problem
		fmt.Println("Payment failed: ", err)

	} else {
		// This block runs when everything is alright.
		fmt.Println("Payment successful!")
	}

	// Example 2: Invalid amount
	fmt.Println("Invalid Amount")
	err = processPayment(-10.00, "1234567890123456")
	if err != nil {
		fmt.Println("Payment failed:", err)
	}

	// Example 3: invalid card

	fmt.Println("Invalid Card")
	err = processPayment(50.00, "123")
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
	}
}
