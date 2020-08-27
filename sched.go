package main

/*
 * A non-idiomatic Go solution, in Go.
 * Uses a binary heap as a priority queue to
 * organize the schedule of jobs.
 * Relies heavily on mutex locks to keep that
 * heap from getting messed up, because it has
 * a separate goroutine waiting for a time.Timer
 * to elapse.
 */

import (
	"fmt"
	"jobscheduler/heap"
	"os"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type action int

const (
	sleep    action = iota
	schedule action = iota
)

func main() {
	s := &Scheduler{}

	s.Start()

	var do action = schedule

	wg := &sync.WaitGroup{}

	serialNo := 0

	for _, str := range os.Args[1:] {

		n, err := strconv.Atoi(str)
		if err != nil {
			continue
		}

		switch do {
		case schedule:

			wg.Add(1)
			x := &schd{
				executeAt: time.Now().Add(time.Millisecond * time.Duration(n)),
				wg:        wg,
				serialNo:  serialNo,
			}
			serialNo++

			s.Schedule(x.runned, n)

			do = sleep
		case sleep:
			sleepInterval := time.Millisecond * time.Duration(n)
			fmt.Printf("\nsleeping for %v\n", sleepInterval)
			time.Sleep(sleepInterval)
			fmt.Println()
			do = schedule
		}
	}

	wg.Wait()
	fmt.Printf("All scheduled jobs done\n")

	s.Stop()
}

type Scheduler struct {
	hpl          sync.Mutex
	h            heap.Heap
	tmr          *time.Timer
	nextDeadline int64
}

func (s *Scheduler) Start() {
}

func (s *Scheduler) Stop() {
}

func (s *Scheduler) Schedule(f func(), n int) {
	scheduleAt := time.Now().Add(time.Duration(n) * time.Millisecond)
	sn := &SchedNode{interval: n, fn: f, desiredTime: scheduleAt, desiredNS: scheduleAt.UnixNano()}

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

func (s *Scheduler) sched(n *SchedNode) {
	s.hpl.Lock()
	s.h = s.h.Insert(n)
	s.hpl.Unlock()

	s.scheduleNext()

	if len(s.h) == 1 {
		s.doNext()
	}
}

func (s *Scheduler) scheduleNext() {
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

type SchedNode struct {
	interval    int       // milliseconds until this function runs
	desiredTime time.Time // desired time for function run
	desiredNS   int64
	fn          func()
}

func (sn *SchedNode) Value() int64 {
	// Should return desiredTime as a single number
	return sn.desiredNS
}

func (sn *SchedNode) IsNil() bool {
	return sn == nil
}

func (sn *SchedNode) String() string {
	return fmt.Sprintf("%d/%s",
		sn.interval, sn.desiredTime.Format(time.RFC3339Nano))
}

type schd struct {
	serialNo  int
	executeAt time.Time
	wg        *sync.WaitGroup
}

func (s *schd) runned() {
	fmt.Printf("%d Now:             %s\n", s.serialNo, time.Now().Format(time.RFC3339Nano))
	fmt.Printf("%d Wanted to run at %s\n\n", s.serialNo, s.executeAt.Format(time.RFC3339Nano))
	s.wg.Done()
}
