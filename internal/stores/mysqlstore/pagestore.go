package mysqlstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"
	"github.com/Pergamene/project-spiderweb-service/internal/util/wrapsql"

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
func (s PageStore) CreatePage(record page.Page, ownerID string) (page.Page, error) {
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
	if s.db == nil {
		return record, &storeerror.DBNotSetUp{}
	}
	t := time.Now()
	record.CreatedAt = &t
	record.UpdatedAt = &t
	id, err := wrapsql.ExecSingleInsert(s.db, wrapsql.InsertQuery{
		IntoTable: "Page",
		InjectedValues: wrapsql.InsertInjectedValues{
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
				{LeftSide: "User.email", Operator: "= ?"},
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
		return isOwner, err
	}
	if isOwner {
		return isOwner, nil
	}
	if guid == "" {
		return isOwner, errors.New("must provide a guid to check privileges")
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
		UpdateTable: "Page",
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
	return wrapsql.ExecSingleUpdate(s.db, query, record.GUID)
}

// GetPage returns back the given page.
func (s PageStore) GetPage(guid string) (page.Page, error) {
	return page.Page{}, errors.New("@TODO: still need to make a different between GetPage and GetEntirePage")
	// if s.db == nil {
	// 	return page.Page{}, &storeerror.DBNotSetUp{}
	// }
	// rows, err := s.db.Query("SELECT `Version_ID`, `PageTemplate_ID`, `title`, `summary`, `permission`, `createdAt`, `updatedAt` FROM `Page` WHERE `guid` = ? AND `deletedAt` IS NULL LIMIT 1", guid)
	// if err != nil {
	// 	return page.Page{}, err
	// }
	// defer rows.Close()
	// var versionID int64
	// var pageTemplateID int64
	// var title string
	// var summary string
	// var permissionString string
	// var createdAt *time.Time
	// var updatedAt *time.Time
	// for rows.Next() {
	// 	err = rows.Scan(&versionID, &pageTemplateID, &title, &summary, &permissionString, &createdAt, &updatedAt)
	// 	if err != nil {
	// 		return page.Page{}, err
	// 	}
	// 	err = rows.Err()
	// 	if err != nil {
	// 		return page.Page{}, err
	// 	}
	// 	p, err := permission.GetPermissionType(permissionString)
	// 	if err != nil {
	// 		return page.Page{}, err
	// 	}
	// 	return page.Page{
	// 		Version: version.Version{
	// 			ID: versionID,
	// 		},
	// 		PageTemplate: pagetemplate.PageTemplate{
	// 			ID: pageTemplateID,
	// 		},
	// 		GUID:           guid,
	// 		Title:          title,
	// 		Summary:        summary,
	// 		PermissionType: p,
	// 		CreatedAt:      createdAt,
	// 		UpdatedAt:      updatedAt,
	// 	}, nil
	// }
	// return page.Page{}, &storeerror.NotFound{
	// 	ID: guid,
	// }
}

// GetEntirePage returns back the given page populated with details, properties, etc.
func (s PageStore) GetEntirePage(guid string) (page.Page, error) {
	return page.Page{}, errors.New("@TODO: GetEntirePage")
}

// GetPages returns a list of pages based on the nextBatchId
func (s PageStore) GetPages(userID string, nextBatchID string) ([]page.Page, int, string, error) {
	return nil, 0, "", errors.New("@TODO: GetPages")
}

// RemovePage marks the given page and removed by setting the deletedAt property.
func (s PageStore) RemovePage(guid string) error {
	return errors.New("@TODO: RemovePage")
}

// GetUniquePageGUID returns a guid for the page that is guaranteed to be unique or errors.
// If the proposedPageGuid is not a zero-value and not unique, it will error.
func (s PageStore) GetUniquePageGUID(proposedPageGUID string) (string, error) {
	return proposedPageGUID, errors.New("@TODO: GetUniquePageGUID")
}
