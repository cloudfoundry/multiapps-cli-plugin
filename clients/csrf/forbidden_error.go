package csrf

const RetryIsNeeded = "retry is needed"

type ForbiddenError struct {
	value string
	ID    int
}

func (e *ForbiddenError) Error() string {
	return RetryIsNeeded
}
