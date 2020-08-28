package scheduler

/*
 * A job scheduler that uses channels to manage concurrency.
 * A single goroutine manages the binary heap that orders functions
 * to run, so no need to lock/unlock all over the place.
 * Uses a binary heap as a priority queue to
 * organize the schedule of jobs.
 */

import (
	"fmt"
	"jobscheduler/heap"
	"runtime"
	"time"
)

type ChannelScheduler struct {
	c            chan heap.Node
	done         chan bool
	h            heap.Heap
	tmr          *time.Timer
	nextDeadline int64
}

// Start part of inteface Scheduler
// Gets the background goroutine running
func (s *ChannelScheduler) Start() {
	s.c = make(chan heap.Node, 0)
	s.done = make(chan bool, 0)
	go s.runScheduling()
}

// Stop part of inteface Scheduler
// Tries to get the background goroutine to return
func (s *ChannelScheduler) Stop() {
	if s.tmr != nil {
		if !s.tmr.Stop() {
			<-s.tmr.C
		}
	}
	s.done <- true
	close(s.done)
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
// The control of waiting for timers to lapse, and/or receive new functions
// to schedule happens in this method.
func (s *ChannelScheduler) runScheduling() {
DONE:
	for {
		if len(s.h) > 0 {
			select {
			case _ = <-s.tmr.C:
				// timer elapsed, run any functions that are due
				s.runFunction()
				s.scheduleNext()
			case node := <-s.c:
				// new function to schedule
				if node != nil {
					s.h = s.h.Insert(node)
					s.scheduleNext()
				}
			case _ = <-s.done:
				// Scheduler.Stop() got called
				break DONE
			}
		} else {
			// Nothing scheduled, wait for new function to schedule to arrive
			select {
			case node := <-s.c:
				if node != nil {
					s.h = s.h.Insert(node)
					s.scheduleNext()
				}
			case _ = <-s.done:
				// Scheduler.Stop() got called
				break DONE
			}

		}
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
// The background goroutine ends up executing this.
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
