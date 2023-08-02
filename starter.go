package starter

import (
	"sync"
	"time"
)

type Pistol struct {
	m       *sync.Mutex
	r       *sync.Cond
	s       *sync.Cond
	waiting int
}

// Ready creates a new starter pistol to sync the start time of multiple runners.
func Ready() *Pistol {
	m := &sync.Mutex{}
	r := sync.NewCond(m)
	s := sync.NewCond(m)

	return &Pistol{
		m:       m,
		r:       r,
		s:       s,
		waiting: 0,
	}
}

// Steady makes sure the expected number of runners are waiting for the Go.
func (p *Pistol) Steady(expect int) *Pistol {
	if expect <= 0 {
		// no-op
		return p
	}

	p.m.Lock()

	// fast path, the number is already satisfied
	if p.waiting >= expect {
		p.m.Unlock()
		return p
	}

	// wait for our condition to be fulfilled
	for p.waiting < expect {
		p.s.Wait()
	}

	p.m.Unlock()
	return p
}

// Go signals the runners to start and returns the current time for easy measurements.
func (p *Pistol) Go() time.Time {
	p.r.Broadcast()
	return time.Now()
}

// Wait blocks the runners until Go has been called.
func (p *Pistol) Wait() {
	p.m.Lock()

	p.waiting++
	p.s.Broadcast() // notify any Steadys that may be waiting
	p.r.Wait()
	p.waiting--

	p.m.Unlock()
}
