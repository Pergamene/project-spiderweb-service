package page

import (
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/Pergamene/project-spiderweb-service/internal/models/version"
)

// Page is the main object that houses most information.
type Page struct {
	ID             int64 `json:"-"`
	Version        version.Version
	GUID           string `json:"id"`
	Title          string `json:"title"`
	Summary        string `json:"summary"`
	PermissionType permission.Type
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
