package store

// HealthcheckStore defines the required functionality for any associated store.
type HealthcheckStore interface {
	IsHealthy() (bool, error)
}
