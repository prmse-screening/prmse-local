package enums

type TaskState int

const (
	Preparing TaskState = iota
	Pending
	Processing
	Completed
	Failed
)
