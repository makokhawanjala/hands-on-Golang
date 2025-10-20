package main

import "fmt"

// Named returns: quotient and remainder are pre-declared
func divide(dividend int, divisor int) (quotient int, remainder int) {
	quotient = dividend / divisor
	remainder = dividend % divisor
	return // naked return - automatically returns quotient and remainder
}

func add(num1 int, num2 int) (answer int) {
	answer = num1 + num2

	return
}

func subtract(num1 int, num2 int) (answer int) {
	answer = num1 - num2

	return
}

func main() {
	q, r := divide(17, 5)
	fmt.Printf("17 ➗ 5 = %d remainder %d\n", q, r)
	answer := add(2, 4)
	answer1 := subtract(20002234334, 43354578)
	fmt.Printf("2️⃣ ➕  4️⃣ 🟰 %d\n", answer)
	fmt.Printf("2️⃣ ➖  4️⃣ 🟰 %d\n", answer1)
}
