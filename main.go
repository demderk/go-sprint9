package main

import (
	"errors"
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

var (
	ErrZeroSize               = errors.New("size must be greater than zero")
	ErrEmptySlice             = errors.New("input slice is empty")
	ErrDataLenLessThanThreads = errors.New("data len less than threads (CHUNKS)")
)

// generateRandomElements generates random elements.
func generateRandomElements(size int) ([]int, error) {
	arr := []int{}
	if size < 1 {
		return []int{}, ErrZeroSize
	}

	for i := 0; i < size; i++ {
		rnd := rand.Int31()
		arr = append(arr, int(rnd))
	}
	return arr, nil
}

// maximum returns the maximum number of elements.
func maximum(data []int) (int, error) {
	if len(data) < 1 {
		return 0, ErrEmptySlice
	}
	return slices.Max(data), nil // :))
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) (int, error) {
	if len(data) < CHUNKS {
		return 0, ErrDataLenLessThanThreads
	}

	chunksMax := []int{}
	chunkSize := len(data) / CHUNKS

	var wait sync.WaitGroup
	var mu sync.Mutex
	wait.Add(CHUNKS)

	errCh := make(chan error, 1)

	for i := 0; i < CHUNKS; i++ {
		go func(i int) {
			defer wait.Done()
			start := chunkSize * i
			end := min(start+chunkSize, len(data))
			max, err := maximum(data[start:end])
			if err != nil {
				select {
				case errCh <- err:
					return
				default:
					return
				}
			}
			mu.Lock()
			chunksMax = append(chunksMax, max)
			mu.Unlock()
		}(i)
	}
	wait.Wait()
	close(errCh)

	if err, ok := <-errCh; ok {
		return 0, fmt.Errorf("maxChunks goroutine was failed: %w", err)
	}

	return maximum(chunksMax)
}

func main() {
	// С go 1.20 rand.seed – Deprecated (по данным VSCode). Тут не забыл, это осознанно.

	fmt.Printf("Генерируем %d целых чисел", SIZE)
	randomData, err := generateRandomElements(SIZE)
	if err != nil {
		fmt.Printf("random number generation was failed: %v", err)
	}

	fmt.Println("Ищем максимальное значение в один поток")
	timeStart := time.Now()
	max, err := maximum(randomData)
	if err != nil {
		fmt.Printf("non-threaded generation was failed: %v", err)
	}

	elapsedNonThreaded := time.Since(timeStart)

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsedNonThreaded)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)

	timeStart = time.Now()
	max, err = maxChunks(randomData)
	if err != nil {
		fmt.Printf("threaded generation was failed: %v", err)
		return
	}

	elapsedThreaded := time.Since(timeStart)

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsedThreaded)

	fmt.Printf("diff: %v", (elapsedNonThreaded - elapsedThreaded).Seconds())
}
