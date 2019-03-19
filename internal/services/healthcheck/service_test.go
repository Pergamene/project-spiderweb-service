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

type isHealthyCall struct {
	returnIsHealthy bool
	returnErr       error
}

func TestIsHealthy(t *testing.T) {
	cases := []struct {
		name            string
		isHealthyCalls  []isHealthyCall
		returnIsHealthy bool
		returnErr       error
	}{
		{
			name:            "db healthy",
			isHealthyCalls:  []isHealthyCall{{returnIsHealthy: true}},
			returnIsHealthy: true,
		},
		{
			name:            "db states its not healthy",
			isHealthyCalls:  []isHealthyCall{{returnIsHealthy: false}},
			returnIsHealthy: false,
		},
		{
			name: "db is not healthy",
			isHealthyCalls: []isHealthyCall{
				{
					returnIsHealthy: true,
					returnErr:       errors.New("failure"),
				},
			},
			returnErr: errors.New("failure"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckStore := new(mocks.HealthcheckStore)
			for index := range tc.isHealthyCalls {
				healthcheckStore.On("IsHealthy").Return(tc.isHealthyCalls[index].returnIsHealthy, tc.isHealthyCalls[index].returnErr)
			}
			healthcheckService = HealthcheckService{
				HealthcheckStore: healthcheckStore,
			}
			isHealthy, err := healthcheckService.IsHealthy(ctx)
			healthcheckStore.AssertNumberOfCalls(t, "IsHealthy", len(tc.isHealthyCalls))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, isHealthy, tc.returnIsHealthy)
		})
	}
}
