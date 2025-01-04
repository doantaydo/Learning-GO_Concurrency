package main

import (
	"fmt"
	"sync"
	"time"
)

// The Dining Philosophers problem is well known in computer science circles.
// Five philosophers, numbered from 0 through 4, live in a house where the table is laied for them;
// each philosophers has their own place at the table.
// Their only difficulty - besides those of philosophy - is that the dish served is a very difficult
// kind of spaghetti which has to be eaten with two forks. These are two forks next to each plate,
// this mean that no two neighbours may be eating simultaneously, since there are five philosophers and five forks.

// This is a simple implementation of Dijkstra's solution to the "Dining Philosophers" dilemma.

// Philosopher is a struct which stores info about a philosopher
type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// philosophers list of all philosophers
var philosophers = []Philosopher{
	{name: "Plato", leftFork: 4, rightFork: 0},
	{name: "Socrates", leftFork: 0, rightFork: 1},
	{name: "Aristotle", leftFork: 1, rightFork: 2},
	{name: "Pascal", leftFork: 2, rightFork: 3},
	{name: "Locke", leftFork: 3, rightFork: 4},
}

// define some variables
var hunger = 3                  // how many times does a person eat?
var eatTime = 1 * time.Second   // how long it takes to earTime
var thinkTime = 3 * time.Second // how long a philosopher thinks
var sleepTime = 1 * time.Second /// how long to wait when printing things out

// define variable to list the order in which philosophers finish dining an leave
var orderMutex sync.Mutex      // a mutex for the slide orderFinishedName
var orderFinishedName []string // the order in which philosophers finish dining an leave

func main() {
	// print out a welcome message
	fmt.Println("Dining Philosophers Problem")
	fmt.Println("---------------------------")
	fmt.Println("The Table is empty.")

	// start the meal
	dine()

	// print out finished message
	fmt.Println("The table is empty.")

	// print the order in which philosophers finish dining an leave
	fmt.Println("Order of completing meal")
	for i, name := range orderFinishedName {
		fmt.Printf("%d. %s\n", i+1, name)
	}
}

func dine() {
	// wg is the WaitGroup that keeps track of how many philosophers are still at the table.
	// When it reaches zero, everyone is finished and has left.
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	// We want everyone to be seated before they start eating.
	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks
	// forks are assigned using the fields leftFork and rightForks in the Philosopher type.
	// Each forks, then, can be found using the index (an integer), and each forks has a unique mutex.
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// fire off a goroutine for the current philosopher
		go diningProblem(philosophers[i], wg, forks, seated)
	}

	wg.Wait()
}

// diningProblem is the function fired off as a goroutine for each of philosophers.
// It takes one philosopher, WaitGroup to determine when everyone is done,
// a map containing the mutexes for every forks on the table, and a WaitGroup used to pause execution of every instance
// of this goroutine until everyone is seated at the table.
func diningProblem(philosopher Philosopher, wg *sync.WaitGroup, forks map[int]*sync.Mutex, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done()

	seated.Wait()

	// eat three times
	for i := hunger; i > 0; i-- {
		// get a lock on both forks
		if philosopher.leftFork > philosopher.rightFork {
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("\t%s takes the left fork.\n", philosopher.name)
			forks[philosopher.rightFork].Lock()
			fmt.Printf("\t%s takes the right fork.\n", philosopher.name)
		}

		fmt.Printf("\t%s has both forks and is eating.\n", philosopher.name)
		time.Sleep(eatTime)

		fmt.Printf("\t%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime)
		forks[philosopher.leftFork].Unlock()
		forks[philosopher.rightFork].Unlock()

		fmt.Printf("\t%s put down the forks.\n", philosopher.name)
	}

	fmt.Println(philosopher.name, "is satisfied.")
	fmt.Println(philosopher.name, "left the table.")
	orderMutex.Lock()
	orderFinishedName = append(orderFinishedName, philosopher.name)
	orderMutex.Unlock()
}
