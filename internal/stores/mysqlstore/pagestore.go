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
	"github.com/Pergamene/project-spiderweb-service/internal/models/property"
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

// UpdatePage sets the given page.
func (s PageStore) UpdatePage(record page.Page) error {
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
	thisPageID := int64(0)
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

func (s PageStore) getPageID(guid string) (int64, error) {
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
	var pageID int64
	err = wrapsql.GetSingleRow(guid, rows, err, &pageID)
	return pageID, err
}

// RemovePage marks the given page and removed by setting the deletedAt property.
func (s PageStore) RemovePage(guid string) error {
	t := time.Now()
	return s.UpdatePage(page.Page{
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

// GetPageProperties returns the page's properties.
func (s PageStore) GetPageProperties(pageGUID string) (returnProperties []property.Property, returnErr error) {
	if pageGUID == "" {
		returnErr = errors.New("must provide pageGUID to get the page properties")
		return
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"Property.ID", "Property.type", "Property.key", "PagePropertyString.value", "PagePropertyNumber.value", "PagePropertyOrder.order"},
		FromTable: "Page",
		JoinClauses: []wrapsql.JoinClause{
			{JoinTable: "PagePropertyString", On: wrapsql.OnClause{LeftSide: "Page.ID", RightSide: "PagePropertyString.Page_ID"}},
			{JoinTable: "PagePropertyNumber", On: wrapsql.OnClause{LeftSide: "Page.ID", RightSide: "PagePropertyNumber.Page_ID"}},
			{JoinTable: "Property", On: wrapsql.OnClause{LeftSide: "PagePropertyString.Property_ID", RightSide: "Property.ID"}},
			{JoinTable: "Property", On: wrapsql.OnClause{LeftSide: "PagePropertyNumber.Property_ID", RightSide: "Property.ID"}},
			{JoinTable: "PagePropertyOrder", On: wrapsql.OnClause{LeftSide: "PagePropertyOrder.Page_ID", RightSide: "Page.ID"}},
			{JoinTable: "PagePropertyOrder", On: wrapsql.OnClause{LeftSide: "PagePropertyOrder.Property_ID", RightSide: "Property.ID"}},
		},
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "Page.guid", Operator: "= ?"},
				{LeftSide: "Page.deletedAt", Operator: "IS NULL"},
				{LeftSide: "Property.deletedAt", Operator: "IS NULL"},
				{LeftSide: "PagePropertyString.deletedAt", Operator: "IS NULL"},
				{LeftSide: "PagePropertyNumber.deletedAt", Operator: "IS NULL"},
			},
		},
		OrderClause: wrapsql.OrderClause{
			Column: "PagePropertyOrder.order",
			SortBy: "ASC",
		},
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), pageGUID)
	if err != nil {
		returnErr = err
		return
	}
	if err := rows.Err(); err != nil {
		returnErr = err
		return
	}
	var orderValues []int64
	defer rows.Close()
	for rows.Next() {
		var orderValue int64
		dbp := property.DBProperty{}
		err := rows.Scan(&dbp.ID, &dbp.Type, &dbp.Key, &dbp.StringValue, &dbp.NumberValue, &orderValue)
		if err != nil {
			returnErr = err
			return
		}
		p, err := dbp.GetProperty()
		if err != nil {
			returnErr = err
			return
		}
		returnProperties = append(returnProperties, p)
		orderValues = append(orderValues, orderValue)
	}
	if len(returnProperties) == 0 {
		returnProperties = make([]property.Property, 0)
	}
	return
}

// ReplacePageProperties replaces the current page's properties with the new properties.
func (s PageStore) ReplacePageProperties(pageGUID string, pageProperties []property.Property) error {
	// @TODO: all this needs to be wrapped into a transaction with rollback.
	if pageGUID == "" {
		return errors.New("must provide pageGUID to replace the page properties")
	}
	pageID, err := s.getPageID(pageGUID)
	if err != nil {
		return errors.Wrapf(err, "unable to get Page.ID for guid: %v", pageID)
	}
	err = s.setPagePropertyIDs(pageProperties)
	if err != nil {
		return errors.Wrap(err, "unable to get Property.ID for the pageProperties")
	}
	err = s.deletePageProperties(pageID)
	if err != nil {
		return errors.Wrap(err, "unable to delete page properties")
	}
	err = s.addPagePropertyOrders(pageID, pageProperties)
	if err != nil {
		return errors.Wrap(err, "unable to add page properties orders")
	}
	err = s.addTypedPageProperties(pageID, pageProperties, property.TypeNumber)
	if err != nil {
		return errors.Wrap(err, "unable to add number type page properties")
	}
	err = s.addTypedPageProperties(pageID, pageProperties, property.TypeString)
	if err != nil {
		return errors.Wrap(err, "unable to add string type page properties")
	}
	return nil
}

func (s PageStore) addPagePropertyOrders(pageID int64, pageProperties []property.Property) error {
	query := wrapsql.BatchInsertQuery{
		IntoTable: "PagePropertyOrder",
	}
	for i, pageProperty := range pageProperties {
		query.BatchInjectedValues["Page_ID"] = append(query.BatchInjectedValues["Page_ID"], pageID)
		query.BatchInjectedValues["Property_ID"] = append(query.BatchInjectedValues["Property_ID"], pageProperty.ID)
		query.BatchInjectedValues["order"] = append(query.BatchInjectedValues["order"], i)
	}
	err := wrapsql.ExecBatchInsert(s.db, query)
	if err != nil {
		return errors.Wrap(err, "unable to insert page property order")
	}
	return nil
}

func (s PageStore) addTypedPageProperties(pageID int64, pageProperties []property.Property, propertyType property.Type) error {
	scopedPageProperties := getTypedProperties(pageProperties, propertyType)
	if len(scopedPageProperties) == 0 {
		return nil
	}
	t := time.Now()
	tableName := ""
	switch propertyType {
	case property.TypeNumber:
		tableName = "PagePropertyNumber"
	case property.TypeString:
		tableName = "PagePropertyString"
	default:
		return errors.Errorf("unsupported page property type for instert: %v", propertyType)
	}
	query := wrapsql.BatchInsertQuery{
		IntoTable: tableName,
	}
	for i, pageProperty := range scopedPageProperties {
		query.BatchInjectedValues["Page_ID"] = append(query.BatchInjectedValues["Page_ID"], pageID)
		query.BatchInjectedValues["Property_ID"] = append(query.BatchInjectedValues["Property_ID"], pageProperty.ID)
		// @TODO: I don't actually think we need a Version_ID anywhere in the schema.  Since upping the version of a page will create a new Page.ID (even if it's the same Page.guid)
		// then properties will be automatically linked to that new version because of the new Page.ID
		query.BatchInjectedValues["Version_ID"] = append(query.BatchInjectedValues["Version_ID"], 0)
		query.BatchInjectedValues["value"] = append(query.BatchInjectedValues["value"], pageProperty.Value)
		// @TODO: figure out permission
		query.BatchInjectedValues["permission"] = append(query.BatchInjectedValues["permission"], "PR")
		query.BatchInjectedValues["createdAt"] = append(query.BatchInjectedValues["createdAt"], t)
		query.BatchInjectedValues["updatedAt"] = append(query.BatchInjectedValues["updatedAt"], t)
	}
	err := wrapsql.ExecBatchInsert(s.db, query)
	if err != nil {
		return errors.Wrap(err, "unable to insert page property order")
	}
	return nil
}

func getTypedProperties(properties []property.Property, propertyType property.Type) (returnProperties []property.Property) {
	returnProperties = make([]property.Property, 0)
	for _, p := range properties {
		if p.Type == propertyType {
			returnProperties = append(returnProperties, p)
		}
	}
	return
}

func (s PageStore) deletePageProperties(pageID int64) error {
	genericWhereClause := wrapsql.WhereClause{
		Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
			{LeftSide: "Page_ID", Operator: "= ?"},
		},
	}
	query := wrapsql.DeleteQuery{
		FromTable:   "PagePropertyOrder",
		WhereClause: genericWhereClause,
	}
	err := wrapsql.ExecDelete(s.db, query, pageID)
	if err != nil {
		return errors.Wrap(err, "unable to delete from PagePropertyOrder")
	}
	query = wrapsql.DeleteQuery{
		FromTable:   "PagePropertyNumber",
		WhereClause: genericWhereClause,
	}
	err = wrapsql.ExecDelete(s.db, query, pageID)
	if err != nil {
		return errors.Wrap(err, "unable to delete from PagePropertyNumber")
	}
	query = wrapsql.DeleteQuery{
		FromTable:   "PagePropertyString",
		WhereClause: genericWhereClause,
	}
	err = wrapsql.ExecDelete(s.db, query, pageID)
	if err != nil {
		return errors.Wrap(err, "unable to delete from PagePropertyString")
	}
	return nil
}

func (s PageStore) setPagePropertyIDs(pageProperties []property.Property) error {
	var keys []string
	for _, p := range pageProperties {
		keys = append(keys, p.Key)
	}
	pps, err := s.getPropertyIDs(keys)
	if err != nil {
		return err
	}
	for _, pp := range pps {
		for i := range pageProperties {
			if pageProperties[i].Key == pp.Key {
				pageProperties[i].ID = pp.ID
			}
		}
	}
	for i, p := range pageProperties {
		if p.ID == 0 {
			return errors.Errorf("unable to find the ID for the property at %v with key %v", i, p.Key)
		}
	}
	return nil
}

func (s PageStore) getPropertyIDs(propertyKeys []string) (returnProperties []property.Property, returnErr error) {
	for i, propertyKey := range propertyKeys {
		if propertyKey == "" {
			return nil, errors.Errorf("property key at %v must be non-zero value", i)
		}
	}
	statement := wrapsql.SelectStatement{
		Selectors: []string{"ID", "key"},
		FromTable: "Property",
		WhereClause: wrapsql.WhereClause{
			Operator: "AND", WhereOperations: []wrapsql.WhereOperation{
				{LeftSide: "key", Operator: "IN (" + wrapsql.GetNValueStubList(len(propertyKeys)) + ")"},
			},
		},
	}
	rows, err := s.db.Query(wrapsql.GetSelectString(statement), propertyKeys)
	if err != nil {
		returnErr = err
		return
	}
	if err := rows.Err(); err != nil {
		returnErr = err
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := property.Property{}
		err := rows.Scan(&p.ID, &p.Key)
		if err != nil {
			returnErr = err
			return
		}
		returnProperties = append(returnProperties, p)
	}
	if len(returnProperties) == 0 {
		returnProperties = make([]property.Property, 0)
	}
	return
}
