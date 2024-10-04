package utility

import (
	"math/rand"
)

// GenerateRandomNumber returns a random number with the specified number of digits.
// If digits is less than 1 or more than 18, it returns 0 because generating such numbers is impractical with int.
func GenerateRandomNumber(digits int) int {
	if digits < 1 || digits > 18 { // The limit is because the range of int can be exceeded otherwise
		return 0 // Returning 0 for cases where it's not possible to generate a number
	}

	min := powerOf10(digits - 1)
	max := powerOf10(digits) - 1

	return rand.Intn(max-min+1) + min
}

// powerOf10 returns 10 raised to the power of n (10^n).
func powerOf10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}
