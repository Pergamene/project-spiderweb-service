package store

import "github.com/Pergamene/project-spiderweb-service/internal/models/appuser"

// UserStore defines the required functionality for any associated store.
type UserStore interface {
	GetUser(userGUID string) (appuser.User, error)
}
