package main

import (
	"fmt"
	"sync"
)

func main() {
	// m := make(map[int]int)
	// var wg sync.WaitGroup

	// for i := 0; i < 100; i++ {
	// 	wg.Add(1)
	// 	go func(val int) {
	// 		defer wg.Done()
	// 		m[val] = val // Race condition: concurrent writes to map
	// 	}(i)
	// }
	// wg.Wait()
	// // Depending on the Go version and specific conditions, this might panic
	// // due to concurrent map writes, or produce an incorrect map.
	// fmt.Println("Map size:", len(m))

	// Option 1: Using sync.Mutex for a regular map
	m := make(map[int]int)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			mu.Lock()
			m[val] = val
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	fmt.Println("Map size (with Mutex):", len(m))

	// Option 2: Using sync.Map (specifically designed for concurrent map access)
	var concurrentMap sync.Map
	var wg2 sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg2.Add(1)
		go func(val int) {
			defer wg2.Done()
			concurrentMap.Store(val, val) // Safe concurrent operation
		}(i)
	}
	wg2.Wait()
	size := 0
	concurrentMap.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	fmt.Println("Map size (with sync.Map):", size)
}
