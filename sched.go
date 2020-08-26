package main

import (
	"fmt"
	"jobscheduler/heap"
	"sync"
	"time"
)

type Scheduler struct {
	hpl          sync.Mutex
	h            heap.Heap
	tmr          *time.Timer
	nextDeadline time.Time
}

func (s *Scheduler) Start() {
}

func (s *Scheduler) Stop() {
}

func (s *Scheduler) Schedule(f func(), n int) {
	scheduleAt := time.Now().Add(time.Duration(n) * time.Millisecond)
	fmt.Printf("Schedule wakeup at       %v\n", scheduleAt.Format(time.RFC3339Nano))
	sn := &SchedNode{interval: n, fn: f, desiredTime: scheduleAt}

	s.sched(sn)
}

func (s *Scheduler) doNext() {
	go func() {
		for {
			<-s.tmr.C
			var n heap.Node
			s.h, n = s.h.Delete()
			sn := n.(*SchedNode)
			(sn.fn)()

			if len(s.h) == 0 {
				break
			}
			// Set timer for next functino
		}
	}()
}

func (s *Scheduler) sched(n *SchedNode) {
	s.hpl.Lock()
	defer s.hpl.Unlock()
	s.h = s.h.Insert(n)

	s.scheduleNext()

	if len(s.h) == 1 {
		s.doNext()
	}
}

func (s *Scheduler) scheduleNext() {
	// figure out interval
	sn := s.h[0].(*SchedNode)
	fmt.Printf("Scheduling for wakeup at %s\n", sn.desiredTime.Format(time.RFC3339Nano))
	interv := sn.desiredTime.Sub(time.Now())
	// interv could be negative, in the past
	fmt.Printf("timer interval: %v\n", interv)

	// might need to update timer instead of creating a new one
	s.tmr = time.NewTimer(interv)
}

type SchedNode struct {
	interval    int       // milliseconds until this function runs
	desiredTime time.Time // desired time for function run
	fn          func()
}

func (sn *SchedNode) Value() int64 {
	// Should return desiredTime as a single number
	return sn.desiredTime.UnixNano()
}

func (sn *SchedNode) IsNil() bool {
	return sn == nil
}

func (sn *SchedNode) String() string {
	return fmt.Sprintf("%d/%s",
		sn.interval, sn.desiredTime.Format(time.RFC3339Nano))
}

type schd struct {
	executeAt time.Time
}

func (s *schd) runned() {
	fmt.Printf("Now:             %s\n", time.Now().Format(time.RFC3339Nano))
	fmt.Printf("Wanted to run at %s\n", s.executeAt.Format(time.RFC3339Nano))
}

func main() {
	s := &Scheduler{}

	s.Start()

	var x schd
	x.executeAt = time.Now().Add(time.Second * 2)
	s.Schedule(x.runned, 2000)

	time.Sleep(5 * time.Second)

	s.Stop()
}
