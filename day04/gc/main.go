package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	// Print initial GC stats
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	fmt.Printf("Initial GC count: %v\n", stats.NumGC)

	// Allocate a lot of objects to trigger GC
	for i := 0; i < 10_000_000; i++ {
		_ = make([]byte, 10) // small allocations
	}

	// Enable GC debug logging (shows pause times)
	debug.SetGCPercent(100)                    // GC runs more often
	debug.SetGCPercent(debug.SetGCPercent(-1)) // restore default

	// Force a GC manually
	fmt.Println("Forcing GC...")
	start := time.Now()
	runtime.GC()
	elapsed := time.Since(start)

	// Read updated GC stats
	runtime.ReadMemStats(&stats)
	fmt.Printf("GC count after: %v\n", stats.NumGC)
	fmt.Printf("Manual GC took: %v\n", elapsed)
}
