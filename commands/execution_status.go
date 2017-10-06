package commands

type ExecutionStatus int

const (
	Success ExecutionStatus = 0
	Failure ExecutionStatus = 1
)

func (status ExecutionStatus) ToInt() int {
	return int(status)
}
