package random

import (
	"testing"
)

func TestGenerateCanonicalFloat(t *testing.T) {
	// Run 1000 iterations to ensure stability and range correctness
	for i := 0; i < 1000; i++ {
		val, err := GenerateCanonicalFloat()
		if err != nil {
			t.Fatalf("Iteration %d: unexpected error: %v", i, err)
		}

		if val < 0.0 || val >= 1.0 {
			t.Errorf("Iteration %d: value %f out of range [0.0, 1.0)", i, val)
		}
	}
}

func TestCreateRandInt(t *testing.T) {
	min := int64(10)
	max := int64(20)

	for i := 0; i < 100; i++ {
		val, err := CreateRandInt(min, max)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if val < min || val > max {
			t.Errorf("Value %d out of range [%d, %d]", val, min, max)
		}
	}
}

func TestGetRandArrVal(t *testing.T) {
	slice := []string{"A", "B", "C"}

	for i := 0; i < 50; i++ {
		val, err := GetRandArrVal(slice)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		found := false
		for _, item := range slice {
			if item == val {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Returned value %v not found in source slice", val)
		}
	}
}
