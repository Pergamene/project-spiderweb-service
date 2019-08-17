package mysqlstore

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/Pergamene/project-spiderweb-service/internal/util/guidgen"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"
	"github.com/pkg/errors"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
)

// PageStore is the mysql for pages
type PageStore struct {
	db *sql.DB
}

// NewPageStore returns a PageStore
func NewPageStore(mysqldb *sql.DB) PageStore {
	return PageStore{
		db: mysqldb,
	}
}

// CreatePage creates a new page.
func (s PageStore) CreatePage(record page.Page, ownerID int64) (page.Page, error) {
	if record.GUID == "" {
		return record, errors.New("must provide record.GUID to create the page")
	}
	if record.Title == "" {
		return record, errors.New("must provide record.Title to create the page")
	}
	if record.Version.ID == 0 {
		return record, errors.New("must provide record.Version.ID to create the page")
	}
	if record.PermissionType == "" {
		return record, errors.New("must provide record.PermissionType to create the page")
	}
	if record.PageTemplate.ID == 0 {
		return record, errors.New("must provide record.PageTemplate.ID to create the page")
	}
	if ownerID == 0 {
		return record, errors.New("must provide ownerID to create the page")
	}
	if s.db == nil {
		return record, &storeerror.DBNotSetUp{}
	}
	t := time.Now()
	record.CreatedAt = &t
	record.UpdatedAt = &t
	id, err := wrapsql.ExecSingleInsert(s.db, wrapsql.InsertQuery{
		IntoTable: "Page",
		InjectedValues: wrapsql.InjectedValues{
			"PageTemplate_ID": record.PageTemplate.ID,
			"Version_ID":      record.Version.ID,
			"guid":            record.GUID,
			"title":           record.Title,
			"summary":         record.Summary,
			"permission":      record.PermissionType,
			"createdAt":       record.CreatedAt,
			"updatedAt":       record.UpdatedAt,
		},
	})
	if err != nil {
		return record, err
	}
	record.ID = id
	_, err = wrapsql.ExecSingleInsert(s.db, wrapsql.InsertQuery{
		IntoTable: "PageOwner",
		InjectedValues: wrapsql.InjectedValues{
			"Page_ID": record.ID,
			"User_ID": ownerID,
			"isOwner": true,
		},
	})
	if err != nil {
		return record, err
	}
	return record, nil
}

// CanEditPage checks if the given user can modify the given page. If not, a storeerror.NotAuthorized will be returned.
// Will also return whether or not the user is the original owner.
func (s PageStore) CanEditPage(guid, userID string) (bool, error) {
	if guid == "" {
		return false, errors.New("must provide a guid to check privileges")
	}
	if userID == "" {
		return false, errors.New("must provide a userID to check privileges")
	}
	if s.db == nil {
		return false, &storeerror.DBNotSetUp{}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"PageOwner.isOwner"},
		FromTable: "PageOwner",
		JoinClauses: []wrapsql.JoinClause{
			{JoinTable: "Page", On: wrapsql.OnClause{LeftSide: "PageOwner.Page_ID", RightSide: "Page.ID"}},
			{JoinTable: "User", On: wrapsql.OnClause{LeftSide: "PageOwner.User_ID", RightSide: "User.ID"}},
		},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "Page.guid", Operator: "= ?"},
				{LeftSide: "User.guid", Operator: "= ?"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid, userID)
	var isOwner bool
	err = wrapsql.GetSingleRow(guid, rows, err, &isOwner)
	if _, ok := err.(*storeerror.NotFound); ok {
		return false, &storeerror.NotAuthorized{
			UserID:  userID,
			TableID: guid,
		}
	}
	return isOwner, err
}

// CanReadPage checks if the given user can read the given page. If not, a storeerror.NotAuthorized will be returned.
// Will also return whether or not the user is the original owner.
func (s PageStore) CanReadPage(guid, userID string) (bool, error) {
	isOwner, err := s.CanEditPage(guid, userID)
	if err != nil {
		if _, ok := err.(*storeerror.NotAuthorized); !ok {
			return isOwner, err
		}
	}
	if isOwner {
		return isOwner, nil
	}
	if s.db == nil {
		return isOwner, &storeerror.DBNotSetUp{}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"permission"},
		FromTable: "Page",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	var pagePermission string
	err = wrapsql.GetSingleRow(guid, rows, err, &pagePermission)
	if _, ok := err.(*storeerror.NotFound); ok {
		return false, &storeerror.NotAuthorized{
			UserID:  userID,
			TableID: guid,
		}
	}
	p, err := permission.GetPermissionType(pagePermission)
	if err != nil {
		return false, err
	}
	return p.IsPublic(), nil
}

// SetPage sets the given page.
func (s PageStore) SetPage(record page.Page) error {
	if record.GUID == "" {
		return errors.New("must provide record.GUID to update the page")
	}
	query := wrapsql.UpdateQuery{
		UpdateTable:    "Page",
		InjectedValues: wrapsql.InjectedValues{},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
			},
		},
	}
	if record.Title != "" {
		query.InjectedValues["title"] = record.Title
	}
	if record.Summary != "" {
		query.InjectedValues["summary"] = record.Summary
	}
	if record.Version.ID != 0 {
		query.InjectedValues["Version_ID"] = record.Version.ID
	}
	if record.PermissionType != "" {
		query.InjectedValues["permission"] = record.PermissionType
	}
	if record.PageTemplate.ID != 0 {
		query.InjectedValues["PageTemplate_ID"] = record.PageTemplate.ID
	}
	if record.DeletedAt != nil {
		query.InjectedValues["deletedAt"] = record.DeletedAt
	}
	return wrapsql.ExecSingleUpdate(s.db, query, record.GUID)
}

// GetPage returns back the given page.
func (s PageStore) GetPage(guid string) (page.Page, error) {
	if guid == "" {
		return page.Page{}, errors.New("must provide guid to get the page")
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"Page.ID", "Version.guid", "PageTemplate.guid", "Page.title", "Page.summary", "Page.permission", "Page.createdAt", "Page.updatedAt"},
		FromTable: "Page",
		JoinClauses: []wrapsql.JoinClause{
			{JoinTable: "Version", On: wrapsql.OnClause{LeftSide: "Page.Version_ID", RightSide: "Version.ID"}},
			{JoinTable: "PageTemplate", On: wrapsql.OnClause{LeftSide: "Page.PageTemplate_ID", RightSide: "PageTemplate.ID"}},
		},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "Page.guid", Operator: "= ?"},
				{LeftSide: "Page.deletedAt", Operator: "IS NULL"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	p := page.Page{
		GUID: guid,
	}
	var permissionString string
	err = wrapsql.GetSingleRow(guid, rows, err, &p.ID, &p.Version.GUID, &p.PageTemplate.GUID, &p.Title, &p.Summary, &permissionString, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return page.Page{}, err
	}
	pt, err := permission.GetPermissionType(permissionString)
	if err != nil {
		return page.Page{}, err
	}
	p.PermissionType = pt
	return p, err
}

// GetPages returns a list of pages based on the nextBatchId
func (s PageStore) GetPages(userID, thisBatchID string, limit int) (pages []page.Page, total int, nextBatchID string, returnErr error) {
	if userID == "" {
		returnErr = errors.New("must provide userID to get pages")
		return
	}
	var err error
	thisPageID := 0
	if thisBatchID != "" {
		thisPageID, err = s.getPageID(thisBatchID)
		if err != nil {
			returnErr = errors.Wrapf(err, "unable to use thisBatchID: %v", thisBatchID)
			return
		}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"Page.guid", "Page.ID", "Version.guid", "PageTemplate.guid", "Page.title", "Page.summary", "Page.permission", "Page.createdAt", "Page.updatedAt"},
		FromTable: "Page",
		JoinClauses: []wrapsql.JoinClause{
			{JoinTable: "PageOwner", On: wrapsql.OnClause{LeftSide: "PageOwner.Page_ID", RightSide: "Page.ID"}},
			{JoinTable: "User", On: wrapsql.OnClause{LeftSide: "PageOwner.User_ID", RightSide: "User.ID"}},
			{JoinTable: "Version", On: wrapsql.OnClause{LeftSide: "Page.Version_ID", RightSide: "Version.ID"}},
			{JoinTable: "PageTemplate", On: wrapsql.OnClause{LeftSide: "Page.PageTemplate_ID", RightSide: "PageTemplate.ID"}},
		},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "Page.ID", Operator: fmt.Sprintf(">= %v", thisPageID)},
				{LeftSide: "User.guid", Operator: "= ?"},
				{LeftSide: "Page.deletedAt", Operator: "IS NULL"},
			},
		},
		Limit: limit + 1, // plus one so we can get an extra record to determine the nextBatchID
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), userID)
	if err != nil {
		returnErr = err
		return
	}
	if err := rows.Err(); err != nil {
		returnErr = err
		return
	}
	var permissionString string
	defer rows.Close()
	for rows.Next() {
		p := page.Page{}
		err := rows.Scan(&p.GUID, &p.ID, &p.Version.GUID, &p.PageTemplate.GUID, &p.Title, &p.Summary, &permissionString, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			returnErr = err
			return
		}
		pt, err := permission.GetPermissionType(permissionString)
		if err != nil {
			returnErr = err
			return
		}
		p.PermissionType = pt
		pages = append(pages, p)
	}
	if len(pages) == 0 {
		pages = make([]page.Page, 0)
	}
	if len(pages) > limit {
		lastPage := pages[len(pages)-1]
		nextBatchID = lastPage.GUID
		pages = pages[:len(pages)-1]
	}
	total, err = s.getTotalPages(userID)
	if err != nil {
		returnErr = err
	}
	return
}

func (s PageStore) getTotalPages(userID string) (int, error) {
	if userID == "" {
		return -1, errors.New("must provide userID to get pages")
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"COUNT(1)"},
		FromTable: "Page",
		JoinClauses: []wrapsql.JoinClause{
			{JoinTable: "PageOwner", On: wrapsql.OnClause{LeftSide: "PageOwner.Page_ID", RightSide: "Page.ID"}},
			{JoinTable: "User", On: wrapsql.OnClause{LeftSide: "PageOwner.User_ID", RightSide: "User.ID"}},
		},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "User.guid", Operator: "= ?"},
				{LeftSide: "Page.deletedAt", Operator: "IS NULL"},
			},
		},
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), userID)
	var total int
	err = wrapsql.GetSingleRow(userID, rows, err, &total)
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (s PageStore) getPageID(guid string) (int, error) {
	if guid == "" {
		return -1, errors.New("must provide guid to get the page id")
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"ID"},
		FromTable: "Page",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "guid", Operator: "= ?"},
			},
		},
		Limit: 1,
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), guid)
	var pageID int
	err = wrapsql.GetSingleRow(guid, rows, err, &pageID)
	return pageID, err
}

// RemovePage marks the given page and removed by setting the deletedAt property.
func (s PageStore) RemovePage(guid string) error {
	t := time.Now()
	return s.SetPage(page.Page{
		GUID:      guid,
		DeletedAt: &t,
	})
}

// GetUniquePageGUID returns a guid for the page that is guaranteed to be unique or errors.
// If the proposedPageGuid is not a zero-value and not unique, it will error.
func (s PageStore) GetUniquePageGUID(proposedPageGUID string) (string, error) {
	return s.getUniquePageGUID(proposedPageGUID, 0)
}

func (s PageStore) getUniquePageGUID(proposedPageGUID string, retry int) (string, error) {
	err := guidgen.CheckProposedGUID(proposedPageGUID, "PG", 15)
	if err != nil {
		return "", err
	}
	return getUniqueGUID(s.db, "PG", 15, "Page", proposedPageGUID, 0)
}
