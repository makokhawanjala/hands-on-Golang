package main

import "fmt"

func main() {
	//active users set
	active := make(map[int]struct{})

	// users logging in
	active[100] = struct{}{}
	active[101] = struct{}{}
	active[102] = struct{}{}
	active[103] = struct{}{}
	active[104] = struct{}{}
	active[105] = struct{}{}

	fmt.Println("Current active users map:")
	fmt.Println(active)

	fmt.Println("Number of active users:", len(active))

	//check if user is active
	user := 102

	if _, ok := active[user]; ok {
		fmt.Printf("User %d is online ✅\n", user)
	}

	// user logs out
	delete(active, user)

	fmt.Println("After user 202 logs out:")
	fmt.Println(active)

	fmt.Println("Number of active users:", len(active))
}

