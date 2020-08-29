package scheduler

import (
	"fmt"
	"jobscheduler/heap"
	"time"
)

// SchedNode holds a function that we want to run at some
// later time. Fits interface heap.Node so that instances
// that fit Scheduler interface can hold it.
// SchedNode fits interface heap.Node
type SchedNode struct {
	interval    int       // milliseconds until this function runs
	desiredTime time.Time // desired time for function run
	desiredNS   int64
	fn          func()
}

// Get compiler to enforce interface compliance
var _ heap.Node = (*SchedNode)(nil)

func (sn *SchedNode) Value() int64 {
	return sn.desiredNS
}

func (sn *SchedNode) IsNil() bool {
	return sn == nil
}

func (sn *SchedNode) String() string {
	return fmt.Sprintf("%s", sn.desiredTime.Format(time.RFC3339Nano))
}
