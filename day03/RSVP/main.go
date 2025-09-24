package main

import (
	"fmt"
	"runtime"
	"time"
)

type Rsvp struct {
	Name, Email, Phone string
	willAttend         bool
}

var responses = make([]*Rsvp, 0, 10)

func main() {
	start := time.Now()

	fmt.Println("Party RSVP system")

	for {
		var name, email, phone, attending string

		fmt.Print("Enter your name (or type 'exit' to quit): ")
		fmt.Scanln(&name)

		if name == "exit" {
			break
		}

		fmt.Print("Enter your email: ")
		fmt.Scanln(&email)

		fmt.Print("Enter your phone: ")
		fmt.Scanln(&phone)

		fmt.Print("Will you attend? (yes/no): ")
		fmt.Scanln(&attending)
		willAttend := attending == "yes"

		guest := &Rsvp{Name: name, Email: email, Phone: phone, willAttend: willAttend}
		responses = append(responses, guest)

	}

	// fmt.Println("\n Final Guest List:")

	// for _, r := range responses {
	//	fmt.Printf("%s {%s,%s} attending? %v\n", r.Name, r.Email, r.Phone, r.willAttend)
	// }

	yes, no := countAttending(responses)
	fmt.Printf("\nSummary: %d attending, %d not attending, total %d\n", yes, no, yes+no)

	var searchName string
	fmt.Print("\nSearch guest by name (or type 'exit'): ")
	fmt.Scanln(&searchName)

	if searchName != "exit" {
		guest := findGuest(responses, searchName)
		if guest != nil {
			fmt.Printf("Found:%s {%s, %s} attending? %v\n", guest.Name, guest.Email, guest.Phone, guest.willAttend)
		} else {
			fmt.Println("Guest not found")
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Memory Usage:\n")
	fmt.Printf("\tAlloc = %v MB\n", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MB\n", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)

}

func countAttending(list []*Rsvp) (int, int) {
	yes, no := 0, 0

	for _, r := range list {
		if r.willAttend {
			yes++
		} else {
			no++
		}
	}
	return yes, no
}

func findGuest(list []*Rsvp, name string) *Rsvp {
	for _, r := range list {
		if r.Name == name {
			return r
		}
	}
	return nil
}
