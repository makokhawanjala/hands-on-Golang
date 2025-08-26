package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func median(nums []float64) float64 {
	n := len(nums)
	if n == 0 { return 0 }
	sorted := append([]float64(nil), nums...)
	sort.Float64s(sorted)
	if n%2==1 { 
		return sorted[n/2] 
	}
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

func sum(nums []float64) float64 {
	s := 0.0
	for _, v := range nums { s += v }
	return s
}

func avg(nums []float64) float64 {
	if len(nums)==0 { return 0 }
	return sum(nums)/float64(len(nums))
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	var nums []float64
	for scanner.Scan() {
		tok := scanner.Text()
		f, err := strconv.ParseFloat(tok, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "skip %q: %v\n", tok, err)
			continue
		}
		nums = append(nums, f)
	}
	fmt.Printf("count=%d sum=%.4f avg=%.4f median=%.4f\n", len(nums), sum(nums), avg(nums), median(nums))
}
