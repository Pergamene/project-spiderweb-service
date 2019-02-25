package pageservice

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/stretchr/testify/require"
)

var pageService PageService
var ctx context.Context

func TestMain(m *testing.M) {
	pageService = PageService{}
	ctx = context.Background()
	result := m.Run()
	os.Exit(result)
}
func TestCreatePage(t *testing.T) {
	cases := []struct {
		name         string
		params       CreatePageParams
		expectedErr  error
		expectedPage page.Page
	}{}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			page, err := pageService.CreatePage(ctx, tc.params)
			errExpected := testErrorAgainstCase(t, err, tc.expectedErr)
			if errExpected {
				return
			}
			require.Equal(t, page, tc.expectedPage)
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
