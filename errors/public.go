package errors

// Public wraps the original error with a new error that
// has a `Public() string` method that will return a message
// that is acceptable to display to the public.
// Added a `StatusCode() int` method to asociate a custom
// status code for the respose.
// This error can also be unwrapped using the traditional
// `errors` package approach.
func Public(err error, msg string, statusCode int) error {
	return publicError{err, msg, statusCode}
}

type publicError struct {
	err        error
	msg        string
	statusCode int
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}

func (pe publicError) StatusCode() int {
	return pe.statusCode
}
