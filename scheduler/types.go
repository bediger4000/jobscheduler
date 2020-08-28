package scheduler

// Scheduler describes a generic scheduling system
type Scheduler interface {
	Start()
	Stop()
	Schedule(f func(), n int)
}
