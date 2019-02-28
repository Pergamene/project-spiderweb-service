package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// GetExpectedNumberOfCalls returns 1 if the call was expected, or 0 otherwise
func GetExpectedNumberOfCalls(isCalled bool) int {
	if isCalled {
		return 1
	}
	return 0
}

// TestErrorAgainstCase returns true if tcErr was expected
func TestErrorAgainstCase(t *testing.T, err error, tcErr error) bool {
	if tcErr != nil {
		require.EqualError(t, err, tcErr.Error())
		return true
	}
	require.NoError(t, err)
	return false
}
