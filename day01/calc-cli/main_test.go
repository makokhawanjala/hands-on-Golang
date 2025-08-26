package main

import "testing"

func TestMedian(t *testing.T) {
	cases := []struct{
		in []float64
		want float64
	}{ 
		{[]float64{}, 0},
		{[]float64{1}, 1},
		{[]float64{1,2}, 1.5},
		{[]float64{3,1,2}, 2},
		{[]float64{10,2,5,7}, 6},
	}
	for _, c := range cases {
		got := median(c.in)
		if got != c.want {
			t.Fatalf("median(%v)=%v want %v", c.in, got, c.want)
		}
	}
}
