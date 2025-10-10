package typing

import (
	"math/rand"
	"strconv"
)

// ShouldInsertNumber returns true with ~10% probability
func ShouldInsertNumber(rng *rand.Rand) bool {
	if rng == nil {
		return false
	}
	return rng.Float64() < 0.1
}

// RandomNumberString generates a random string of digits with length between 1 and maxLen.
// The first digit is always 1-9 to avoid leading zeros, and subsequent digits are 0-9.
// This mirrors the behavior of Monkeytype's getNumbers function.
func RandomNumberString(rng *rand.Rand, maxLen int) string {
	if rng == nil || maxLen <= 0 {
		return ""
	}

	// Random length between 1 and maxLen
	length := rng.Intn(maxLen) + 1

	result := ""
	for i := 0; i < length; i++ {
		var digit int
		if i == 0 {
			// First digit: 1-9 (no leading zero)
			digit = rng.Intn(9) + 1
		} else {
			// Subsequent digits: 0-9
			digit = rng.Intn(10)
		}
		result += strconv.Itoa(digit)
	}

	return result
}
