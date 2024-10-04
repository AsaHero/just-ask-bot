package utility

import (
	"math/rand"
	"time"
)

func RandomInt(from int, to int) int {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random integer within the specified range
	num1 := r.Intn(to-from+1) + from

	return num1
}

func TwoRandomInt(from int, to int) (int, int) {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate two random integers within the specified range
	num1 := r.Intn(to-from+1) + from
	num2 := r.Intn(to-from+1) + from

	return num1, num2
}

func ThereRandomInt(from int, to int) (int, int, int) {
	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate three random integers within the specified range
	num1 := r.Intn(to-from+1) + from
	num2 := r.Intn(to-from+1) + from
	num3 := r.Intn(to-from+1) + from

	return num1, num2, num3
}
