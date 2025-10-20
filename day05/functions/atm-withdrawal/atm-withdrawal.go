package main

import (
	"fmt"
)

// ATM gives you money and a message receipt
func withdrawMoney(accountBalance float64, amount float64) (float64, string) {
	if amount > accountBalance {
		return accountBalance, "insuffiecient funds! ❌"
	}

	newBalance := accountBalance - amount
	message := fmt.Sprintf("Withdrawal was successful! New balance: $%.2f", newBalance)

	return newBalance, message
}

func main() {
	// Using the ATM
	balance := 1000.0
	balance, message := withdrawMoney(balance, 250.0)
	fmt.Println(message)
	fmt.Println(balance)

}
