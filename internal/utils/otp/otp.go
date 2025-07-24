package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateOTP(length int) string {
	rand.Seed(time.Now().UnixNano())
	min := int64(1 * intPow(10, length-1))
	max := int64(intPow(10, length) - 1)
	return fmt.Sprintf("%d", rand.Int63n(max-min)+min)
}

func intPow(a, b int) int {
	result := 1
	for b > 0 {
		result *= a
		b--
	}
	return result
}
