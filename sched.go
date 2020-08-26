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
	now := time.Now()
	now.Add(time.Duration(n) * time.Millisecond)

	sn := &SchedNode{interval: n, fn: f, desiredTime: now}

	s.sched(sn)
}

func (s *Scheduler) doNext() {
	go func() {
		<-s.tmr.C
		var n heap.Node
		s.h, n = s.h.Delete()
		sn := n.(*SchedNode)
		(sn.fn)()

		if len(s.h) == 0 {
			return
		}
	}()
}

func (s *Scheduler) sched(n *SchedNode) {
	s.hpl.Lock()
	defer s.hpl.Unlock()
	s.h = s.h.Insert(n)

	if len(s.h) == 1 {
		// figure out interval
		sn := s.h[0].(*SchedNode)
		interv := sn.desiredTime.Sub(time.Now())
		// interv could be negative, in the past

		// start go routine
		s.tmr = time.NewTimer(interv)
		s.doNext()
		return
	}

	sn := s.h[0].(*SchedNode)
	intrvl := sn.desiredTime.Sub(s.nextDeadline)
	fmt.Printf("%v\n", intrvl)
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
	fmt.Printf("Now: %s\n", time.Now().Format(time.RFC3339Nano))
}

func main() {
	s := &Scheduler{}

	s.Start()

	var x schd
	x.executeAt = time.Now().Add(time.Second)

	time.Sleep(5 * time.Second)

	s.Schedule(x.runned, 2000)
	s.Stop()
}
