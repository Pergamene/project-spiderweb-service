package testutils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestErrorAgainstCase returns true if tcErr was expected
func TestErrorAgainstCase(t *testing.T, err error, tcErr error) bool {
	if tcErr != nil {
		require.EqualError(t, err, tcErr.Error())
		return true
	}
	require.NoError(t, err)
	return false
}
