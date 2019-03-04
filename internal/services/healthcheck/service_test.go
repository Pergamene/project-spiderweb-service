package healthcheckservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/store/mocks"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

var healthcheckService HealthcheckService
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	result := m.Run()
	os.Exit(result)
}
func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name                     string
		storeIsHealthyCalled     bool
		storeIsHealthyReturnBool bool
		storeIsHealthyReturnErr  error
		returnIsHealthy          bool
		returnErr                error
	}{
		{
			name:                     "db healthy",
			storeIsHealthyCalled:     true,
			storeIsHealthyReturnBool: true,
			returnIsHealthy:          true,
		},
		{
			name:                     "db states its not healthy",
			storeIsHealthyCalled:     true,
			storeIsHealthyReturnBool: false,
			returnIsHealthy:          false,
		},
		{
			name:                    "db is not healthy",
			storeIsHealthyCalled:    true,
			storeIsHealthyReturnErr: errors.New("failure"),
			returnErr:               errors.New("failure"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckStore := new(mocks.HealthcheckStore)
			healthcheckStore.On("IsHealthy").Return(tc.storeIsHealthyReturnBool, tc.storeIsHealthyReturnErr)
			healthcheckService = HealthcheckService{
				HealthcheckStore: healthcheckStore,
			}
			isHealthy, err := healthcheckService.IsHealthy(ctx)
			healthcheckStore.AssertNumberOfCalls(t, "IsHealthy", testutils.GetExpectedNumberOfCalls(tc.storeIsHealthyCalled))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, isHealthy, tc.returnIsHealthy)
		})
	}
}
