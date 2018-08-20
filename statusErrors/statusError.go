package statusErrors

//Error extends go error with method for HTTP status
type Error interface {
	error
	Status() int
}

//StatusError holds an error with an html response code
type StatusError struct {
	Err        error `json:"error,omitempty"`
	StatusCode int   `json:"statusCode,omitempty"`
}

//StatusError implements error interface
func (e StatusError) Error() string {
	return e.Err.Error()
}

//Status returns HTTP status code
func (e StatusError) Status() int {
	return e.StatusCode
}

//New returns a new StatusError
func New(err error, statusCode int) StatusError {
	return StatusError{
		Err:        err,
		StatusCode: statusCode,
	}
}