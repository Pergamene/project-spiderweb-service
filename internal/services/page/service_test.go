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

type canEditPageCall struct {
	paramPageGUID   string
	paramPageUserID string
	returnIsOwner   bool
	returnErr       error
}

type updatePageCalls struct {
	paramPage page.Page
	returnErr error
}

func TestUpdatePage(t *testing.T) {
	cases := []struct {
		name             string
		params           UpdatePageParams
		canEditPageCalls []canEditPageCall
		updatePageCalls  []updatePageCalls
		returnErr        error
	}{
		{
			name:   "test proper update",
			params: getUpdatePageParams("PG_1", "UR_1"),
			canEditPageCalls: []canEditPageCall{
				{
					paramPageGUID:   "PG_1",
					paramPageUserID: "UR_1",
				},
			},
			updatePageCalls: []updatePageCalls{{paramPage: getDefaultUpdatePagePage("PG_1")}},
		},
		{
			name:   "test unauthorized update",
			params: getUpdatePageParams("PG_1", "UR_1"),
			canEditPageCalls: []canEditPageCall{
				{
					paramPageGUID:   "PG_1",
					paramPageUserID: "UR_1",
					returnErr:       getStoreUnauthorizedErr("UR_1", "PG_1", nil),
				},
			},
			returnErr: errors.New("User UR_1 is not authorized to perform the action on the ID PG_1"),
		},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			pageStore := new(mocks.PageStore)
			for index := range tc.canEditPageCalls {
				pageStore.On("CanEditPage", tc.canEditPageCalls[index].paramPageGUID, tc.canEditPageCalls[index].paramPageUserID).Return(tc.canEditPageCalls[index].returnIsOwner, tc.canEditPageCalls[index].returnErr)
			}
			for index := range tc.updatePageCalls {
				pageStore.On("UpdatePage", tc.updatePageCalls[index].paramPage).Return(tc.updatePageCalls[index].returnErr)
			}
			pageService = PageService{
				PageStore: pageStore,
			}
			err := pageService.UpdatePage(ctx, tc.params)
			pageStore.AssertNumberOfCalls(t, "CanEditPage", len(tc.canEditPageCalls))
			pageStore.AssertNumberOfCalls(t, "UpdatePage", len(tc.updatePageCalls))
			errExpected := testutils.TestErrorAgainstCase(t, err, tc.returnErr)
			if errExpected {
				return
			}
		})
	}
}
