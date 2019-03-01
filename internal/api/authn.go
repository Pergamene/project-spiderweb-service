package api

import (
	"net/http"
)

// AuthN struct for fulfilling authentication
type AuthN struct {
	Datacenter      string
	AdminAuthSecret string
}

// Different header key names
const (
	AdminAuthSecretHeaderKey = "X-ADMIN-AUTH-SECRET"
	UserIDHeaderKey          = "X-USER-ID"
)

// FailedAuthentication is an error that signifies that the request failed authentication.
type FailedAuthentication struct{}

func (e *FailedAuthentication) Error() string {
	return "not authenticated"
}

// Authenticate first checks if we are running locally, if so it will load AuthData from the headers
// it will also check for an admin auth secret authentication.
// FailedAuthentication will be returned if they are not authenticated.
func (a AuthN) Authenticate(r *http.Request) (AuthData, error) {
	if a.isAdmin(r) {
		if a.hasUserID(r) {
			return AuthData{
				Type:   AuthTypeProxyUser,
				UserID: r.Header.Get(UserIDHeaderKey),
			}, nil
		}
		return AuthData{
			Type: AuthTypeAdmin,
		}, nil
	}
	return AuthData{}, &FailedAuthentication{}
}

func (a AuthN) isAdmin(r *http.Request) bool {
	return r.Header.Get(AdminAuthSecretHeaderKey) == a.AdminAuthSecret || a.Datacenter == LocalDatacenterEnv
}

func (a AuthN) hasUserID(r *http.Request) bool {
	return r.Header.Get(UserIDHeaderKey) != ""
}
