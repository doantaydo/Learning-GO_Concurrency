package main

import (
	"fmt"
	"sync"
)

func printSomeThing(s string, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println(s)
}

func main() {
	// simple GoRoutine
	// go printSomeThing("This is the first thing to be printed!")
	// time.Sleep(1 * time.Second)

	// use WaitGroup
	var wg sync.WaitGroup

	words := []string{"alpha", "beta", "delta", "gamma", "pi", "zeta", "eta", "theta", "epsilon"}

	wg.Add(len(words))

	for i, x := range words {
		go printSomeThing(fmt.Sprintf("%d: %s", i, x), &wg)
	}

	wg.Wait()

	wg.Add(1)
	printSomeThing("This is the second thing to be printed!", &wg)
}
