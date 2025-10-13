package main

import (
	"fmt"
	"time"
)

func main() {
	// Current time
	now := time.Now()
	fmt.Println("Now:     ", now.Format("02-Jan-2006"))

	// Add 1 day
	tomorrow := time.Now().AddDate(0, 0, 1)
	fmt.Println("Tommorrow:		", tomorrow.Format("02-Jan-2006"))

	// Subtract 1 day
	yesterday := now.AddDate(0, 0, -1)
	fmt.Println("Yesterday:  ", yesterday.Format("02-Jan-2006"))

	// Add 1 month
	nextMonth := now.AddDate(0, 1, 0)
	fmt.Println("Next Month:	", nextMonth.Format("02-Jan-2006"))

	// Subtract 1 Year
	lastYear := now.AddDate(-1, 0, 0)
	fmt.Println("Last Year:		", lastYear.Format("02-Jan-2006"))

	// Add 1 year, 2 months, 10 days
	custom := now.AddDate(1, 2, 10)
	fmt.Println("Custom (+1Y +2M +10D", custom.Format("02-Jan-2006"))
}
