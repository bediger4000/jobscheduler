package main

import (
	"flag"
	"fmt"
	"jobscheduler/scheduler"
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
	useScheduler := flag.String("s", "locking", "select scheduler locking or channel")
	flag.Parse()

	var s scheduler.Scheduler

	switch *useScheduler {
	case "channel":
		s = &scheduler.ChannelScheduler{}
	case "locking":
		s = &scheduler.LockingScheduler{}
	}

	s.Start()

	var do action = schedule
	wg := &sync.WaitGroup{}
	serialNo := 0

	for i := 0; i < flag.NArg(); i++ {

		n, err := strconv.Atoi(flag.Arg(i))
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

	s.Stop()

	// Sleep for a while to let Scheduler.Stop() play out.
	time.Sleep(2 * time.Second)
	fmt.Printf("All scheduled jobs done\n")
}
