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
	"jobscheduler/scheduler"
	"os"
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
	s := &scheduler.LockingScheduler{}

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
			x := &scheduler.Schd{
				ExecuteAt: time.Now().Add(time.Millisecond * time.Duration(n)),
				Wg:        wg,
				SerialNo:  serialNo,
			}
			serialNo++

			s.Schedule(x.Run, n)

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
