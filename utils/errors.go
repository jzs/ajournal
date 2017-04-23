package utils

// NewAPIError returns a new API error that wraps the underlying error. Use this in Services such
// that the underlying error is not propagated to the end user.
func NewAPIError(err error, code int, desc string) error {
	return APIError{err: err, Status: code, Desc: desc}
}

// APIError type
type APIError struct {
	err    error
	Status int
	Desc   string
}

func (e APIError) Error() string {
	return e.Desc
}
