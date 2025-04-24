package enums

type TaskState int

const (
	Preparing TaskState = iota
	Pending
	Processing
	Completed
	Failed
)

func (s TaskState) String() string {
	switch s {
	case Preparing:
		return "Preparing"
	case Pending:
		return "Pending"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}
