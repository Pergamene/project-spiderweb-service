package api

// InternalErr is an error that signifies that the request failed unexpectedly.
type InternalErr struct{}

func (e *InternalErr) Error() string {
	return "internal server error"
}
