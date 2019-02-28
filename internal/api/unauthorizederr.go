package api

// UnauthorizedErr is an error that signifies that the request failed authentication.
type UnauthorizedErr struct{}

func (e *UnauthorizedErr) Error() string {
	return "unauthorized error"
}
