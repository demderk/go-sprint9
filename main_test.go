package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomElementsZeroSize(t *testing.T) {
	res := generateRandomElements(0)
	require.Equal(t, res, []int{})
}

func TestMaximumValuesThreadedLessThanCHUNKS(t *testing.T) {
	arr := []int{}
	if CHUNKS < 2 {
		t.Fatalf("thread count less than 1")
	}
	for i := 0; i < CHUNKS-1; i++ {
		arr = append(arr, 10)
	}
	res := maxChunks(arr)
	require.Equal(t, 0, res)
}

func TestMaximumValuesThreadedOK(t *testing.T) {
	arr := []int{1, 2, 5, 45, 12, 5, 100, 65}
	res := maxChunks(arr)
	require.Equal(t, 100, res)

}

func TestMaximumValues(t *testing.T) {
	testCases := []struct {
		name     string
		input    []int
		expected int
	}{
		{"EmptySlice", []int{}, 0},
		{"SliceLenOne", []int{8}, 8},
		{"NormalValues", []int{1, 2, 5, 10, 5, 4, 3}, 10},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			require.Equal(t, testCase.expected, maximum(testCase.input))
		})
	}
}
