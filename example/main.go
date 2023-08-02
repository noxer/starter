package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/noxer/starter"
)

func main() {
	pistol := starter.Ready()
	wg := &sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()

			prepareRunner()
			pistol.Wait()
			run()
		}()
	}

	fmt.Println("Waiting for the runners to be ready...")
	start := pistol.Steady(100).Go()
	wg.Wait()
	fmt.Printf("It took %s to complete the tasks\n", time.Since(start))
}

func prepareRunner() { time.Sleep(time.Second) }
func run()           { /* do something important */ }
