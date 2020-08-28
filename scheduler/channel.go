package scheduler

import (
	"fmt"
	"jobscheduler/heap"
	"runtime"
	"time"
)

type ChannelScheduler struct {
	c            chan heap.Node
	h            heap.Heap
	tmr          *time.Timer
	nextDeadline int64
}

func (s *ChannelScheduler) Start() {
	s.c = make(chan heap.Node, 0)
	go s.runScheduling()
}

func (s *ChannelScheduler) Stop() {
	if s.tmr != nil {
		if !s.tmr.Stop() {
			<-s.tmr.C
		}
	}
	close(s.c)
}

// Schedule part of inteface Scheduler
// Creates a SchedNode and puts it on the chan to the goroutine
// running runScheduling()
func (s *ChannelScheduler) Schedule(f func(), n int) {
	scheduleAt := time.Now().Add(time.Duration(n) * time.Millisecond)
	s.c <- &SchedNode{interval: n, fn: f, desiredTime: scheduleAt, desiredNS: scheduleAt.UnixNano()}
}

// runScheduling gets run by a single goroutine, from Scheduler.Start
func (s *ChannelScheduler) runScheduling() {
	for {
		if s.tmr != nil {
			select {
			case _ = <-s.tmr.C:
				// timer elapsed, run any functions that are due
				s.tmr = nil
				s.runFunction()
			case node := <-s.c:
				// new function to schedule
				if node == nil {
					// Scheduler.Stop() closed channel
					break
				}
				s.h = s.h.Insert(node)
			}
		} else {
			// Nothing scheduled, wait for new function to schedule to arrive
			node := <-s.c
			if node == nil {
				// Scheduler.Stop() closed channel
				break
			}
			s.h = s.h.Insert(node)
		}

		s.scheduleNext()
	}
}

// runFunction runs all the functions that need
// to execute right now. Slop will build up and it's
// possible for more than 1 function to want to run.
func (s *ChannelScheduler) runFunction() {
	nowNS := time.Now().UnixNano()
	for {
		var n heap.Node
		s.h, n = s.h.Delete()
		sn := n.(*SchedNode)
		go (sn.fn)()
		runtime.Gosched()

		if len(s.h) == 0 {
			break
		}

		if s.h[0].Value() > nowNS {
			break
		}
	}
}

// scheduleNext figures out what interval to sleep for the
// next function-to-run to execute at the proper time.
func (s *ChannelScheduler) scheduleNext() {
	if len(s.h) < 1 {
		return
	}
	// figure out interval
	sn := s.h[0].(*SchedNode)
	interv := sn.desiredTime.Sub(time.Now())
	// interv could be negative, in the past

	// might need to update timer instead of creating a new one
	// this would happen if scheduler.Schedule() gets called with
	// an interval of X - 1 millisec when s.tmr has X millisec left
	// before firing.
	if s.h[0].Value() < s.nextDeadline {
		s.tmr.Stop()
		s.tmr.Reset(interv)
		s.nextDeadline = s.h[0].Value()
		fmt.Printf("Rescheduling for wakeup at %s\n\n", sn.desiredTime.Format(time.RFC3339Nano))
		return
	}
	fmt.Printf("Scheduling for wakeup at %s\n\n", sn.desiredTime.Format(time.RFC3339Nano))
	s.nextDeadline = s.h[0].Value()
	s.tmr = time.NewTimer(interv)
}
