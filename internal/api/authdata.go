package api

import (
	"context"

	"github.com/pkg/errors"
)

// AuthType is a valid type of authentication.
type AuthType string

// valid AuthType values.
const (
	AuthTypeAdmin     AuthType = "admin"
	AuthTypeProxyUser AuthType = "proxyUser"
)

// AuthData are the data for authn/authz.
type AuthData struct {
	Type   AuthType
	UserID string
}

type authKeyType string

const authKey = authKeyType("auth")

// SetDataOnContext sets AuthData on the context.
func SetDataOnContext(ctx context.Context, authData AuthData) context.Context {
	return context.WithValue(ctx, authKey, authData)
}

// GetDataFromContext returns AuthData from the context.
func GetDataFromContext(ctx context.Context) (AuthData, error) {
	d, ok := ctx.Value(authKey).(AuthData)
	if !ok {
		return d, errors.New("authData not found in context")
	}
	return d, nil
}

// IsAdmin returns true if the AuthData is of admin.
func (ad AuthData) IsAdmin() bool {
	return ad.Type == AuthTypeAdmin
}
