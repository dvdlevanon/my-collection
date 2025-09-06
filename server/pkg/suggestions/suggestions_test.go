package suggestions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUniqueRandomNumbers_NoInfiniteLoop(t *testing.T) {
	t.Run("Request more numbers than available", func(t *testing.T) {
		// This should not hang and should return at most 3 numbers
		result := getUniqueRandomNumbers(3, 10)
		assert.LessOrEqual(t, len(result), 3, "Should not return more numbers than available")
		assert.GreaterOrEqual(t, len(result), 1, "Should return at least 1 number when max > 0")
	})

	t.Run("Request zero numbers", func(t *testing.T) {
		result := getUniqueRandomNumbers(5, 0)
		assert.Empty(t, result, "Should return empty slice when count is 0")
	})

	t.Run("Max is zero", func(t *testing.T) {
		result := getUniqueRandomNumbers(0, 5)
		assert.Empty(t, result, "Should return empty slice when max is 0")
	})

	t.Run("Normal case", func(t *testing.T) {
		result := getUniqueRandomNumbers(10, 5)
		assert.Len(t, result, 5, "Should return exactly 5 numbers")

		// Verify all numbers are unique
		seen := make(map[int]bool)
		for _, num := range result {
			assert.False(t, seen[num], "Number %d should be unique", num)
			assert.GreaterOrEqual(t, num, 1, "Number should be >= 1")
			assert.LessOrEqual(t, num, 10, "Number should be <= max")
			seen[num] = true
		}
	})

	t.Run("Equal count and max", func(t *testing.T) {
		result := getUniqueRandomNumbers(5, 5)
		assert.Len(t, result, 5, "Should return exactly 5 numbers")

		// Should contain all numbers from 1 to 5
		seen := make(map[int]bool)
		for _, num := range result {
			seen[num] = true
		}
		for i := 1; i <= 5; i++ {
			assert.True(t, seen[i], "Should contain number %d", i)
		}
	})
}
