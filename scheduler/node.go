package scheduler

import (
	"fmt"
	"time"
)

// SchedNode holds a function that we want to run at some
// later time. Fits interface heap.Node so that instances
// that fit Scheduler interface can hold it.
type SchedNode struct {
	interval    int       // milliseconds until this function runs
	desiredTime time.Time // desired time for function run
	desiredNS   int64
	fn          func()
}

func (sn *SchedNode) Value() int64 {
	return sn.desiredNS
}

func (sn *SchedNode) IsNil() bool {
	return sn == nil
}

func (sn *SchedNode) String() string {
	return fmt.Sprintf("%s", sn.desiredTime.Format(time.RFC3339Nano))
}
