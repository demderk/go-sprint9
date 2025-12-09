package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateRandomElementsZeroSize(t *testing.T) {
	_, err := generateRandomElements(0)
	require.ErrorIs(t, err, ErrZeroSize)
}

func TestMaximumEmptySlice(t *testing.T) {
	_, err := maximum([]int{})
	require.ErrorIs(t, err, ErrEmptySlice)
}

func TestMaximumSliceLenOne(t *testing.T) {
	res, err := maximum([]int{8})
	require.NoError(t, err)
	require.Equal(t, 8, res)
}

func TestMaximumValues(t *testing.T) {
	res, err := maximum([]int{1, 2, 5, 10, 5, 4, 3})
	require.NoError(t, err)
	require.Equal(t, 10, res)
}

func TestMaximumValuesThreadedLessThanCHUNKS(t *testing.T) {
	arr := []int{}
	if CHUNKS < 2 {
		t.Fatalf("thread count less than 1")
	}
	for i := 0; i < CHUNKS-1; i++ {
		arr = append(arr, 10)
	}
	_, err := maxChunks(arr)
	require.ErrorIs(t, err, ErrDataLenLessThanThreads)
}

func TestMaximumValuesThreadedOK(t *testing.T) {
	arr := []int{1, 2, 5, 45, 12, 5, 100, 65}
	res, err := maxChunks(arr)
	require.NoError(t, err)
	require.Equal(t, 100, res)
}
