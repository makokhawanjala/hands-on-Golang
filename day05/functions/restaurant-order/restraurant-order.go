package main

import "fmt"

// You are ordering food at the restaurant
func placeOrder(dishName string, quantity int, isSpicy bool) string {
	order := fmt.Sprintf("Order: %d x %s", quantity, dishName)

	if isSpicy {
		order += " (EXTRA HOT. 🌶️ )"
	}

	return order
}

func main() {
	// customer ordering
	order1 := placeOrder("Chicken Wings", 12, true)
	order2 := placeOrder("Samosa", 15, false)
	order3 := placeOrder("Choma", 3, false)

	fmt.Println(order1)
	fmt.Println(order2)
	fmt.Println(order3)
}
