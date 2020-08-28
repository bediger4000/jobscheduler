package scheduler

import (
	"fmt"
	"sync"
	"time"
)

type Schd struct {
	SerialNo  int
	ExecuteAt time.Time
	Wg        *sync.WaitGroup
}

func (s *Schd) Run() {
	fmt.Printf("%d Now:             %s\n", s.SerialNo, time.Now().Format(time.RFC3339Nano))
	fmt.Printf("%d Wanted to run at %s\n\n", s.SerialNo, s.ExecuteAt.Format(time.RFC3339Nano))
	s.Wg.Done()
}
