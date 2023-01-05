package main

import (
	"fmt"
)

// greatest common divisor (GCD) via Euclidean algorithm
func GCD(a, b int) int {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

// find Least Common Multiple (LCM) via GCD
func LCM(a, b int) int {
	result := a * b / GCD(a, b)
	return result
}

func main() {
	fmt.Println(GCD(10, 15))
	fmt.Println(LCM(2, 3))
}
