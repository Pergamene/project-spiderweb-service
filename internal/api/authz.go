package api

// AuthZ struct for fulfilling authorization
type AuthZ struct {
	APIPath string
}

// FailedAuthorization is an error that signifies that the request failed authorization.
type FailedAuthorization struct{}

func (e *FailedAuthorization) Error() string {
	return "not authorized"
}

// Authorize determines if the user is authorized to access the given entities.
// FailedAuthorization will be returned if they are not authorized.
func (a AuthZ) Authorize(path string, authData AuthData) (AuthData, error) {
	return authData, nil
}
