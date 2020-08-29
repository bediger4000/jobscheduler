package scheduler

import (
	"fmt"
	"sync"
	"time"
)

// Schd holds information so that it can show when
// a func should have run, and when it did run.
// Also has a WaitGroup so that the main goroutine can
// exit at an appropriate time.
type Schd struct {
	SerialNo  int
	ExecuteAt time.Time
	Wg        *sync.WaitGroup
}

// Run says when it executed, when it should have executed,
// and notifies a waiting thread via a WaitGroup
func (s *Schd) Run() {
	fmt.Printf("%d Now:             %s\n", s.SerialNo, time.Now().Format(time.RFC3339Nano))
	fmt.Printf("%d Wanted to run at %s\n\n", s.SerialNo, s.ExecuteAt.Format(time.RFC3339Nano))
	s.Wg.Done()
}
