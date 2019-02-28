package pageservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store/mocks"
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

func getUpdatePageParams(guid, userID string) UpdatePageParams {
	return UpdatePageParams{
		Page:   getDefaultUpdatePagePage(guid),
		UserID: userID,
	}
}

func getDefaultUpdatePagePage(guid string) page.Page {
	return getPage(guid, "This is a new title", "This is a new summary")
}

func getPage(guid, title, summary string) page.Page {
	return page.Page{
		GUID:    guid,
		Title:   title,
		Summary: summary,
	}
}

func getStoreUnauthorizedErr(userID, tableID string, err error) error {
	return &storeerror.NotAuthorized{
		UserID:  userID,
		TableID: tableID,
		Err:     err,
	}
}

func TestUpdatePage(t *testing.T) {
	cases := []struct {
		name                                    string
		params                                  UpdatePageParams
		storeAssertCanModifyPageCalled          bool
		storeAssertCanModifyPageParamPageGUID   string
		storeAssertCanModifyPageParamPageUserID string
		storeAssertCanModifyPageReturnIsOwner   bool
		storeAssertCanModifyPageReturnErr       error
		storeUpdatePageCalled                   bool
		storeUpdatePageParamPage                page.Page
		storeUpdatePageReturnErr                error
		expectedErr                             error
	}{
		{
			name:   "test proper update",
			params: getUpdatePageParams("PG_1", "UR_1"),
			storeAssertCanModifyPageCalled:          true,
			storeAssertCanModifyPageParamPageGUID:   "PG_1",
			storeAssertCanModifyPageParamPageUserID: "UR_1",
			storeUpdatePageCalled:                   true,
			storeUpdatePageParamPage:                getDefaultUpdatePagePage("PG_1"),
		},
		{
			name:   "test unauthorized update",
			params: getUpdatePageParams("PG_1", "UR_1"),
			storeAssertCanModifyPageCalled:          true,
			storeAssertCanModifyPageParamPageGUID:   "PG_1",
			storeAssertCanModifyPageParamPageUserID: "UR_1",
			storeAssertCanModifyPageReturnErr:       getStoreUnauthorizedErr("UR_1", "PG_1", nil),
			storeUpdatePageCalled:                   false,
			expectedErr:                             errors.New("User UR_1 is not authorized to perform the action on the ID PG_1"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageStore := new(mocks.PageStore)
			pageStore.On("AssertCanModifyPage", tc.storeAssertCanModifyPageParamPageGUID, tc.storeAssertCanModifyPageParamPageUserID).Return(tc.storeAssertCanModifyPageReturnIsOwner, tc.storeAssertCanModifyPageReturnErr)
			pageStore.On("UpdatePage", tc.storeUpdatePageParamPage).Return(tc.storeUpdatePageReturnErr)
			pageService.PageStore = pageStore
			err := pageService.UpdatePage(ctx, tc.params)
			pageStore.AssertNumberOfCalls(t, "AssertCanModifyPage", getExpectedNumberOfCalls(tc.storeAssertCanModifyPageCalled))
			pageStore.AssertNumberOfCalls(t, "UpdatePage", getExpectedNumberOfCalls(tc.storeUpdatePageCalled))
			errExpected := testErrorAgainstCase(t, err, tc.expectedErr)
			if errExpected {
				return
			}
		})
	}
}

// returns 1 if the call was expected, or 0 otherwise
func getExpectedNumberOfCalls(isCalled bool) int {
	if isCalled {
		return 1
	}
	return 0
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
