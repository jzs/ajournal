package utils

func NewAPIError(err error, code int, desc string) error {
	return APIError{err: err, Status: code, Desc: desc}
}

type APIError struct {
	err    error
	Status int
	Desc   string
}

func (e APIError) Error() string {
	return e.Desc
}
