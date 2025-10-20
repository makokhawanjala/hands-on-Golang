package main

import "fmt"

// addItemsToCart adds new items to the slice and returns the updated cart
func addItemsToCart(cart []string, items ...string) []string {
	for _, item := range items {
		cart = append(cart, item)
	}
	return cart
}

// updateInventory updates an item quantity in a map
func updateInventory(inventory map[string]int, item string, quantity int) {
	inventory[item] = quantity
}

func main() {
	// Slice example
	myCart := []string{"Milk", "Bread"}
	myCart = addItemsToCart(myCart, "Eggs", "Butter", "Sugar", "Salt", "Mango Juice")
	fmt.Println("🛒 My Cart:", myCart)

	// Map example
	stock := map[string]int{"Apples": 50, "Oranges": 30}
	updateInventory(stock, "Apples", 45)
	fmt.Println("📦 Inventory:", stock)
}
