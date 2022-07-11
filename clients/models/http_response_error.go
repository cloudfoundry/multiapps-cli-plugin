package models

type HttpResponseError struct {
	Underlying error
}

func (e HttpResponseError) Error() string {
	return e.Underlying.Error()
}
