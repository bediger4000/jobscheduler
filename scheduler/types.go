// package scheduler implements a generic function scheduling
// system.
// Contains two working scheduling systems, one using
// mutex locks to handle concurrent data access, the other
// Go channels.
package scheduler

// Scheduler describes a generic scheduling system
type Scheduler interface {
	Start()
	Stop()
	Schedule(f func(), n int)
}
