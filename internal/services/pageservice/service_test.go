package pageservice

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/Pergamene/project-spiderweb-service/internal/util/testutils"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/stores/store/mocks"
)

var pageService PageService
var ctx context.Context

func TestMain(m *testing.M) {
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
		name                         string
		params                       UpdatePageParams
		CanModifyPageCalled          bool
		CanModifyPageParamPageGUID   string
		CanModifyPageParamPageUserID string
		CanModifyPageReturnIsOwner   bool
		CanModifyPageReturnErr       error
		storeUpdatePageCalled        bool
		storeUpdatePageParamPage     page.Page
		storeUpdatePageReturnErr     error
		returnErr                    error
	}{
		{
			name:                         "test proper update",
			params:                       getUpdatePageParams("PG_1", "UR_1"),
			CanModifyPageCalled:          true,
			CanModifyPageParamPageGUID:   "PG_1",
			CanModifyPageParamPageUserID: "UR_1",
			storeUpdatePageCalled:        true,
			storeUpdatePageParamPage:     getDefaultUpdatePagePage("PG_1"),
		},
		{
			name:                         "test unauthorized update",
			params:                       getUpdatePageParams("PG_1", "UR_1"),
			CanModifyPageCalled:          true,
			CanModifyPageParamPageGUID:   "PG_1",
			CanModifyPageParamPageUserID: "UR_1",
			CanModifyPageReturnErr:       getStoreUnauthorizedErr("UR_1", "PG_1", nil),
			storeUpdatePageCalled:        false,
			returnErr:                    errors.New("User UR_1 is not authorized to perform the action on the ID PG_1"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageStore := new(mocks.PageStore)
			pageStore.On("CanModifyPage", tc.CanModifyPageParamPageGUID, tc.CanModifyPageParamPageUserID).Return(tc.CanModifyPageReturnIsOwner, tc.CanModifyPageReturnErr)
			pageStore.On("UpdatePage", tc.storeUpdatePageParamPage).Return(tc.storeUpdatePageReturnErr)
			pageService = PageService{
				PageStore: pageStore,
			}
			err := pageService.UpdatePage(ctx, tc.params)
			pageStore.AssertNumberOfCalls(t, "CanModifyPage", testutils.GetExpectedNumberOfCalls(tc.CanModifyPageCalled))
			pageStore.AssertNumberOfCalls(t, "UpdatePage", testutils.GetExpectedNumberOfCalls(tc.storeUpdatePageCalled))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
		})
	}
}
