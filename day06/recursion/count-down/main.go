package main

import "fmt"

func countDown(n int) {
	if n == 0 {
		fmt.Println("Blast Off!")
		return
	}
	fmt.Println(n)
	countDown(n - 1)
}

func main() {
	countDown(5)
}
