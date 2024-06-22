package main

import "fmt"

func main() {
	number := []string{}
	for i := 0; i < 523; i++ {
		number = append(number, fmt.Sprint(i))
	}

	batches := make([][]string, 6)
	batchSize := 100
	for i, v := range number {
		batchIndex := i / batchSize
		batch := batches[batchIndex]
		batches[batchIndex] = append(batch, v)
	}
}
