package starter

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNotBlocking(t *testing.T) {
	Ready().Steady(0).Go()
}

func TestOneRunner(t *testing.T) {
	p := Ready()
	start := time.Now()
	done := make(chan struct{})

	go func() {
		time.Sleep(2 * time.Second)
		p.Wait()
		close(done)
	}()

	now := p.Steady(1).Go()

	require.GreaterOrEqual(t, now.Sub(start), time.Second)
	<-done // this will time out if the channel has not been closed
}

func TestNRunners(t *testing.T) {
	counts := []int{1, 10, 100, 1_000, 10_000, 100_000, 1_000_000}
	for _, count := range counts {
		count := count
		t.Run(strconv.Itoa(count), func(t *testing.T) {
			t.Parallel() // these tests take a long time, run them in parallel

			p := Ready()
			wg := &sync.WaitGroup{}
			wg.Add(count)

			t.Logf("Starting %d runners...", count)
			for i := 0; i < count; i++ {
				go func() {
					defer wg.Done()

					// sleep between 1 and 3 seconds
					time.Sleep(time.Duration(rand.Int63n(int64(2*time.Second))) + time.Second)

					p.Wait()
				}()
			}
			t.Logf("Started %d runners...", count)

			p.Steady(count)
			t.Logf("%d runners ready", count)

			start := p.Go()
			wg.Wait()
			t.Logf("Time to unlock %d runners: %s", count, time.Since(start))
		})
	}
}
