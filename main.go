package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

const (
	SIZE   = 100_000_000
	CHUNKS = 8
)

// generateRandomElements generates random elements.
func generateRandomElements(size int) []int {
	arr := []int{}

	for i := 0; i < size; i++ {
		rnd := rand.Int31()
		arr = append(arr, int(rnd))
	}
	return arr
}

// maximum returns the maximum number of elements.
func maximum(data []int) int {
	if len(data) < 1 {
		return 0
	}
	max := data[0]
	for i := 0; i < len(data); i++ {
		if max < data[i] {
			max = data[i]
		}
	}
	return max
}

// maxChunks returns the maximum number of elements in a chunks.
func maxChunks(data []int) int {
	if len(data) < CHUNKS {
		return 0
	}

	chunksMax := make([]int, CHUNKS)
	chunkSize := len(data) / CHUNKS

	var wait sync.WaitGroup

	wait.Add(CHUNKS)
	for i := 0; i < CHUNKS; i++ {
		go func(i int) {
			defer wait.Done()
			start := chunkSize * i
			end := min(start+chunkSize, len(data))
			max := maximum(data[start:end])

			chunksMax[i] = max
		}(i)
	}
	wait.Wait()

	return maximum(chunksMax)
}

func main() {
	// С go 1.20 rand.seed – Deprecated (по данным VSCode). Тут не забыл, это осознанно.

	fmt.Printf("Генерируем %d целых чисел", SIZE)
	randomData := generateRandomElements(SIZE)

	fmt.Println("Ищем максимальное значение в один поток")
	timeStart := time.Now()
	max := maximum(randomData)

	elapsedNonThreaded := time.Since(timeStart)

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsedNonThreaded)

	fmt.Printf("Ищем максимальное значение в %d потоков\n", CHUNKS)

	timeStart = time.Now()
	max = maxChunks(randomData)

	elapsedThreaded := time.Since(timeStart)

	fmt.Printf("Максимальное значение элемента: %d\nВремя поиска: %d ms\n", max, elapsedThreaded)

	fmt.Printf("diff: %v", (elapsedNonThreaded - elapsedThreaded).Seconds())
}
