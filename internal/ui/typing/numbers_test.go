package typing

import (
	"math/rand"
	"testing"
)

func TestShouldInsertNumber(t *testing.T) {
	if ShouldInsertNumber(nil) {
		t.Fatalf("expected nil RNG to return false")
	}

	var foundTrue, foundFalse bool
	for seed := int64(0); seed < 1_000 && (!foundTrue || !foundFalse); seed++ {
		control := rand.New(rand.NewSource(seed))
		expected := control.Float64() < 0.1

		rng := rand.New(rand.NewSource(seed))
		if actual := ShouldInsertNumber(rng); actual != expected {
			t.Fatalf("expected result %v for seed %d, got %v", expected, seed, actual)
		}
		if expected {
			foundTrue = true
		} else {
			foundFalse = true
		}
	}

	if !foundTrue || !foundFalse {
		t.Fatalf("expected to observe both true and false outcomes within seed search; got true=%v false=%v", foundTrue, foundFalse)
	}
}

func TestRandomNumberString(t *testing.T) {
	if result := RandomNumberString(nil, 5); result != "" {
		t.Fatalf("expected nil rng to return empty string")
	}
	rng := rand.New(rand.NewSource(1))
	result := RandomNumberString(rng, 3)
	if len(result) == 0 || len(result) > 3 {
		t.Fatalf("expected random number length between 1 and 3, got %q", result)
	}
	if result[0] == '0' {
		t.Fatalf("expected first digit to be non-zero, got %q", result)
	}
}
