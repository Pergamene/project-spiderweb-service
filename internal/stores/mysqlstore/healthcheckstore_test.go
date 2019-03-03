package mysqlstore

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"
	"github.com/stretchr/testify/require"
)

func TestHealthcheckIsHealthy(t *testing.T) {
	cases := []struct {
		name                   string
		shouldReplaceDBWithNil bool
		preTestQueries         []string
		returnIsHealthy        bool
		returnErr              error
	}{
		{
			name:            "db healthy",
			preTestQueries:  []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"ok\")"},
			returnIsHealthy: true,
		},
		{
			name:            "db not healthy",
			preTestQueries:  []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"er\")"},
			returnIsHealthy: false,
		},
		{
			name:            "db not healthy because entry doesn't exist",
			preTestQueries:  []string{},
			returnIsHealthy: false,
		},
		{
			name:                   "db not setup",
			shouldReplaceDBWithNil: true,
			preTestQueries:         []string{"INSERT INTO `healthcheck` (`status`) VALUES (\"ok\")"},
			returnIsHealthy:        false,
			returnErr:              errors.New("DB is not configured"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			healthcheckStore := HealthcheckStore{
				db: mysqldb,
			}
			if tc.shouldReplaceDBWithNil {
				healthcheckStore.db = nil
			}
			err := clearTableForTest(healthcheckStore.db, "healthcheck")
			require.NoError(t, err)
			err = execPreTestQueries(healthcheckStore.db, tc.preTestQueries)
			require.NoError(t, err)
			isHealthy, err := healthcheckStore.IsHealthy()
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
			require.Equal(t, tc.returnIsHealthy, isHealthy)
		})
	}
}
