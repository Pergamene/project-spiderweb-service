package mysqlstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

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
	statement, err := s.db.Prepare("INSERT INTO Page (`PageTemplate_ID`, `Version_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( ?, ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return record, err
	}
	t := time.Now()
	record.CreatedAt = &t
	record.UpdatedAt = &t
	defer statement.Close()
	result, err := statement.Exec(record.PageTemplate.ID, record.Version.ID, record.GUID, record.Title, record.Summary, record.PermissionType, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return record, err
	}
	id, err := result.LastInsertId()
	record.ID = id
	return record, nil
}

// CanEditPage checks if the given user can modify the given page. If not, a storeerror.NotAuthorized will be returned.
// Will also return whether or not the user is the original owner.
func (s PageStore) CanEditPage(pageGUID, userID string) (bool, error) {
	if pageGUID == "" {
		return false, errors.New("must provide a pageGUID to check privileges")
	}
	if userID == "" {
		return false, errors.New("must provide a userID to check privileges")
	}
	if s.db == nil {
		return false, &storeerror.DBNotSetUp{}
	}
	rows, err := s.db.Query(`
		SELECT PageOwner.isOwner 
		FROM PageOwner 
		JOIN Page ON PageOwner.Page_ID = Page.ID 
		JOIN User ON PageOwner.User_ID = User.ID
		WHERE Page.guid = ? 
		AND User.email = ? 
		LIMIT 1`, pageGUID, userID)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	var isOwner bool
	for rows.Next() {
		err = rows.Scan(&isOwner)
		if err != nil {
			return false, err
		}
		err = rows.Err()
		if err != nil {
			return false, err
		}
		return isOwner, nil
	}
	return false, &storeerror.NotAuthorized{
		UserID:  userID,
		TableID: pageGUID,
	}
}

// CanReadPage checks if the given user can read the given page. If not, a storeerror.NotAuthorized will be returned.
// Will also return whether or not the user is the original owner.
func (s PageStore) CanReadPage(pageGUID, userID string) (bool, error) {
	isOwner, err := s.CanEditPage(pageGUID, userID)
	if err != nil {
		return isOwner, err
	}
	if isOwner {
		return isOwner, nil
	}
	if pageGUID == "" {
		return isOwner, errors.New("must provide a pageGUID to check privileges")
	}
	if s.db == nil {
		return isOwner, &storeerror.DBNotSetUp{}
	}
	rows, err := s.db.Query("SELECT `permission` FROM `Page` WHERE `guid` = ?", pageGUID)
	if err != nil {
		return isOwner, err
	}
	defer rows.Close()
	var pagePermission string
	for rows.Next() {
		err = rows.Scan(&pagePermission)
		if err != nil {
			return isOwner, err
		}
		err = rows.Err()
		if err != nil {
			return isOwner, err
		}
		p, err := permission.GetPermissionType(pagePermission)
		if err != nil {
			return isOwner, err
		}
		if p.IsPublic() {
			return isOwner, nil
		}
	}
	return isOwner, &storeerror.NotAuthorized{
		UserID:  userID,
		TableID: pageGUID,
	}
}

// SetPage sets the given page.
func (s PageStore) SetPage(record page.Page) error {
	if record.GUID == "" {
		return errors.New("must provide record.GUID to update the page")
	}
	if record.Title == "" {
		return errors.New("must provide record.Title to update the page")
	}
	return errors.New("@TODO: still need to update all the other stuff that CreatePage has")
	// if s.db == nil {
	// 	return &storeerror.DBNotSetUp{}
	// }
	// statement, err := s.db.Prepare("UPDATE `Page` SET `title` = ?, `summary` = ? WHERE `guid` = ?")
	// if err != nil {
	// 	return err
	// }
	// defer statement.Close()
	// _, err = statement.Exec(record.Title, record.Summary, record.GUID)
	// return err
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
