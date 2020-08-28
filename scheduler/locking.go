package scheduler

/*
 * A non-idiomatic Go solution, in Go.
 * Uses a binary heap as a priority queue to
 * organize the schedule of jobs.
 * Relies heavily on mutex locks to keep that
 * heap from getting messed up, because it has
 * a separate goroutine waiting for a time.Timer
 * to elapse. Other goroutines can insert nodes
 * into the binary heap, so something has to keep
 * the background goroutine and other goroutines from
 * racing to fix up the heap.
 *
 * Has a small advantage in that a background goroutine
 * doesn't always execute. If nothing is in the heap,
 * the background goroutine returns.
 */

import (
	"fmt"
	"jobscheduler/heap"
	"runtime"
	"sync"
	"time"
)

// LockingScheduler should fit interface Scheduler
type LockingScheduler struct {
	hpl          sync.Mutex
	h            heap.Heap
	tmr          *time.Timer
	nextDeadline int64
}

// Start makes *LockingScheduler fit interface Scheduler
func (s *LockingScheduler) Start() {
}

// Stop makes *LockingScheduler fit interface Scheduler
func (s *LockingScheduler) Stop() {
}

// Schedule makes *LockingScheduler fit interface Scheduler,
// but also gets called by random goroutines to get functions
// to run in the future.
func (s *LockingScheduler) Schedule(f func(), n int) {
	scheduleAt := time.Now().Add(time.Duration(n) * time.Millisecond)
	sn := &SchedNode{interval: n, fn: f, desiredTime: scheduleAt, desiredNS: scheduleAt.UnixNano()}

	s.sched(sn)
}

func (s *LockingScheduler) doNext() {
	go func() {
		for {
			<-s.tmr.C

			nowNS := time.Now().UnixNano()

			// pick up all the jobs to run, there may be
			// some sloppiness that means more than 1 job
			// wants to run right now
			for {

				s.hpl.Lock()
				var n heap.Node
				s.h, n = s.h.Delete()
				s.hpl.Unlock()
				sn := n.(*SchedNode)
				go (sn.fn)()
				runtime.Gosched()

				s.hpl.Lock()
				if len(s.h) == 0 {
					s.hpl.Unlock()
					break
				}

				if s.h[0].Value() > nowNS {
					s.hpl.Unlock()
					break
				}
				s.hpl.Unlock()
			}

			// Set timer for next functino
			s.scheduleNext()
		}
	}()
}

func (s *LockingScheduler) sched(n *SchedNode) {
	s.hpl.Lock()
	s.h = s.h.Insert(n)
	s.hpl.Unlock()

	s.scheduleNext()

	s.hpl.Lock()
	if len(s.h) == 1 {
		s.doNext()
	}
	s.hpl.Unlock()
}

func (s *LockingScheduler) scheduleNext() {
	s.hpl.Lock()
	defer s.hpl.Unlock()
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
