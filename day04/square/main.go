package main

import "fmt"

func squareAndCube(n int) (int, int) {
	return n * n, n * n * n
}

func main() {
	s, _ := squareAndCube(5)
	fmt.Println("Square:", s)

	_, c := squareAndCube(3)
	fmt.Println("Cube:", c)

}
