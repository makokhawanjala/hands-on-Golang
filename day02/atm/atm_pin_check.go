package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// seed random generator
	// Using deprecated rand.Seed for compatibility
	rand.Seed(time.Now().UnixNano())

	// Generate random 4-digit PIN
	correctPIN := fmt.Sprintf("%04d", rand.Intn(10000))
	fmt.Println("[DEBUG] Today's PIN is:", correctPIN)

	var enteredPIN string
	attempts := 0

	// User enters a PIN
	// If correct ---> show Access granted
	// If wrong ---> Allow Retry
	// After 3 wrong attempts ---> "Card Blocked"
	for attempts < 3 {
		fmt.Print("Enter your 4-digit PIN: ")
		if _, err := fmt.Scanln(&enteredPIN); err != nil {
			fmt.Println("Error reading PIN:", err)
			continue
		}

		if enteredPIN == correctPIN {
			fmt.Println("Access granted. Welcome to your account.")

			// Give option to change pin
			fmt.Print("Do you want to change your PIN? (yes/no): ")
			var choice string
			if _, err := fmt.Scanln(&choice); err != nil {
				fmt.Println("Error reading choice:", err)
				return
			}

			if choice == "yes" {
				var newPIN string
				fmt.Print("Enter your new 4-digit PIN: ")
				if _, err := fmt.Scanln(&newPIN); err != nil {
					fmt.Println("Error reading new PIN:", err)
					return
				}
				correctPIN = newPIN
				fmt.Println("PIN successfully changed! Your new PIN is set.")
				fmt.Printf("[DEBUG] New PIN is now: %s\n", correctPIN)
			} else {
				fmt.Println("Do you want to transact?")
			}
			return
		} else {
			attempts++
			fmt.Printf("Wrong PIN at %s. Attempts left: %d\n",
				time.Now().Format("2006-01-02 15:04:05"),
				3-attempts)
		}
	}
	fmt.Println("Card blocked. Too many wrong attempts.")
}
