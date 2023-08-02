[![GoDoc](https://godoc.org/github.com/noxer/starter?status.svg)](https://godoc.org/github.com/noxer/starter)
[![Go Report Card](https://goreportcard.com/badge/github.com/noxer/starter)](https://goreportcard.com/report/github.com/noxer/starter)

# starter
For some performance tests, I need to prepare a number of runners and then wait for all of them to be ready. Only then should the experiment be started. In order to coordinate those runners, I've written this metaphorical starter pistol.

## Installation
You can import this module into your program using `go get`:

```sh
go get -u github.com/noxer/starter
```

## Example
A short example to get you started:

```go
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
```

## Usage
Initialize a new starter pistol by calling `starter.Ready()`:

```go
pistol := starter.Ready()
```

Now you pass the pistol to your runners (this is where the metaphor falls apart):

```go
go func() {
    // prepare the runner
    // ...

    // then wait for the bang
    pistol.Wait()

    // and do the work
    // ...
}()
```

You can now wait for the runners to be ready (aka. waiting):

```go
// wait for 10 runners to be ready
pistol.Steady(10)
```

It's perfectly fine to call `Steady` multiple times, you can use it to provide an indication about how many runners are ready:

```go
pistol.Steady(1)
fmt.Println("First runner is ready!")

pistol.Steady(100)
fmt.Println("All runners are ready!")
```

When all runners are ready, give the signal. `Go` will return the current timestamp, so you can measure the time it took:

```go
start := pistol.Go()
```

Be aware that it may take some time for the workers to actually start, waking up a million runners may take 100ms.

Although discuraged, you can reuse the pistol for the next run. This also means that runners that are late to call `Wait` will block even if `Go` had previously been called.

## License
MIT
