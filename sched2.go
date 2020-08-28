package main

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
	s := &scheduler.ChannelScheduler{}

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

	s.Stop()

	// Sleep for a while to let Scheduler.Stop() play out.
	time.Sleep(2 * time.Second)
	fmt.Printf("All scheduled jobs done\n")
}
