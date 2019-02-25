package env

import (
	"os"

	"github.com/pkg/errors"
)

// Require retrieves the env var or errors if it doesn't exist.
func Require(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", errors.Errorf("could not find env var \"%v\"", key)
	}
	return value, nil
}

// Get retrieves the env var or returns the default if it doesn't exist.
func Get(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
