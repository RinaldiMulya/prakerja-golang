package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wg.Add(len(numbers))
	for _, number := range numbers {
		go func(n int) {
			defer wg.Done()
			(func() {
				fmt.Println(n)
			})() // IIFE
		}(number)
	}
	wg.Wait()
}