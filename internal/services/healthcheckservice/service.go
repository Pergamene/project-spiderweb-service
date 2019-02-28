package healthcheckservice

import (
	"context"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/store"
)

// HealthcheckService is the service for handling healthcheck-related APIs
type HealthcheckService struct {
	HealthcheckStore store.HealthcheckStore
}

// IsHealthy creates a new healthcheck.
func (s HealthcheckService) IsHealthy(ctx context.Context) (bool, error) {
	return s.HealthcheckStore.IsHealthy()
}
