package healthcheckservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var healthcheckService HealthcheckService
var ctx context.Context

func TestMain(m *testing.M) {
	healthcheckService = HealthcheckService{}
	ctx = context.Background()
	result := m.Run()
	os.Exit(result)
}
func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name                string
		expectedErr         error
		expectedHealthcheck bool
	}{
		{
			name:                "db healthy",
			expectedHealthcheck: true,
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			isHealthy, err := healthcheckService.IsHealthy(ctx)
			errExpected := testErrorAgainstCase(t, err, tc.expectedErr)
			if errExpected {
				return
			}
			require.Equal(t, isHealthy, tc.expectedHealthcheck)
		})
	}
}

// returns true if tcErr was expected
func testErrorAgainstCase(t *testing.T, err error, tcErr error) bool {
	if tcErr != nil {
		require.EqualError(t, err, tcErr.Error())
		return true
	}
	require.NoError(t, err)
	return false
}
