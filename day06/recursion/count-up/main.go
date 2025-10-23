package main

import "fmt"

func countUp(n int) {
	if n == 10 {
		fmt.Println("Blast Off!")
		return
	}
	fmt.Println(n)
	countUp(n + 1)
}

func main() {
	countUp(0)
}
