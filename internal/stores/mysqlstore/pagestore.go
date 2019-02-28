package mysqlstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/stores/storeerror"

	"github.com/Pergamene/project-spiderweb-service/internal/models/page"
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
		return record, errors.New("must provide record.Version to create the page")
	}
	if record.PermissionType == "" {
		return record, errors.New("must provide record.PermissionType to create the page")
	}

	statement, err := s.db.Prepare("INSERT INTO Page (`Version_ID`, `guid`, `title`, `summary`, `permission`, `createdAt`, `updatedAt`) VALUES( ?, ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return record, err
	}
	record.CreatedAt = time.Now()
	record.UpdatedAt = time.Now()
	defer statement.Close()
	result, err := statement.Exec(record.Version.ID, record.GUID, record.Title, record.Summary, record.PermissionType, record.CreatedAt, record.UpdatedAt)
	if err != nil {
		return record, err
	}
	id, err := result.LastInsertId()
	record.ID = id
	return record, nil
}

// CanModifyPage checks if the given user can modify the given page. If not, a storeerror.NotAuthorized will be returned.
// Will also return whether or not the user is the original owner.
func (s PageStore) CanModifyPage(pageGUID, userID string) (bool, error) {
	if pageGUID == "" {
		return false, errors.New("must provide a pageGUID to check privileges")
	}
	if userID == "" {
		return false, errors.New("must provide a userID to check privileges")
	}
	rows, err := s.db.Query("SELECT `PageOwner`.`isOwner` FROM `PageOwner` JOIN `Page` ON `PageOwner`.`Page_ID` = `Page`.`ID` WHERE `Page`.`guid` = ? AND `PageOwner`.`User_ID` = ? LIMIT 1", pageGUID, userID)
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
		return isOwner, nil
	}
	err = rows.Err()
	if err != nil {
		return false, err
	}
	return false, &storeerror.NotAuthorized{
		UserID:  userID,
		TableID: pageGUID,
	}
}

// UpdatePage updatse the given page.
func (s PageStore) UpdatePage(record page.Page) error {
	if record.GUID == "" {
		return errors.New("must provide record.GUID to update the page")
	}
	if record.Title == "" {
		return errors.New("must provide record.Title to update the page")
	}
	statement, err := s.db.Prepare("UPDATE `Page` SET `title` = ?, `summary` = ? WHERE `guid` = ?")
	if err != nil {
		return err
	}
	defer statement.Close()
	_, err = statement.Exec(record.Title, record.Summary, record.GUID)
	return err
}
