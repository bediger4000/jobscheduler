package main

import (
	"fmt"
	"jobscheduler/heap"
	"sync"
	"time"
)

func main() {
	s := &Scheduler{}

	s.Start()

	var x schd
	x.executeAt = time.Now().Add(time.Second * 2)
	s.Schedule(x.runned, 2000)

	var y schd
	y.executeAt = time.Now().Add(time.Second * 3)
	s.Schedule(y.runned, 3000)

	var z schd
	z.executeAt = time.Now().Add(time.Second * 7)
	s.Schedule(z.runned, 7000)

	fmt.Printf("sleeping 10 seconds\n")
	time.Sleep(10 * time.Second)

	s.Stop()
}

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

			nowNS := time.Now().UnixNano()

			// pick up all the jobs to run, there may be
			// some sloppiness that means more than 1 job
			// wants to run right now
			for {

				var n heap.Node
				s.h, n = s.h.Delete()
				sn := n.(*SchedNode)
				(sn.fn)()

				if len(s.h) == 0 {
					break
				}

				if s.h[0].Value() > nowNS {
					break
				}
			}

			// Set timer for next functino
			s.scheduleNext()
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
	fmt.Printf("heap has length %d\n", len(s.h))
}

func (s *Scheduler) scheduleNext() {
	if len(s.h) < 1 {
		return
	}
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
