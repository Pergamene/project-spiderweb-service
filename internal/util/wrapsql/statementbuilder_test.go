package wrapsql

import (
	"os"
	"testing"
	"time"

	"github.com/Pergamene/project-spiderweb-service/internal/models/permission"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	result := m.Run()
	os.Exit(result)
}

func TestGetSelectString(t *testing.T) {
	cases := []struct {
		name                 string
		paramSelectStatement SelectStatement
		returnStatement      string
	}{
		{
			name: "test 'get page' statement",
			paramSelectStatement: SelectStatement{
				Selectors: []string{"Page.ID", "Version.guid", "PageTemplate.guid", "Page.title", "Page.summary", "Page.permission", "Page.createdAt", "Page.updatedAt"},
				FromTable: "Page",
				JoinClauses: []JoinClause{
					{JoinTable: "Version", On: OnClause{LeftSide: "Page.Version_ID", RightSide: "Version.ID"}},
					{JoinTable: "PageTemplate", On: OnClause{LeftSide: "Page.PageTemplate_ID", RightSide: "PageTemplate.ID"}},
				},
				WhereClause: WhereClause{
					Operator: "AND", WhereOperations: []WhereOperation{
						{LeftSide: "guid", Operator: "= ?"},
						{LeftSide: "deletedAt", Operator: "IS NULL"},
					},
				},
				Limit: 1,
			},
			returnStatement: "SELECT `Page`.`ID`,`Version`.`guid`,`PageTemplate`.`guid`,`Page`.`title`,`Page`.`summary`,`Page`.`permission`,`Page`.`createdAt`,`Page`.`updatedAt` FROM Page JOIN Version ON `Page`.`Version_ID` = `Version`.`ID` JOIN PageTemplate ON `Page`.`PageTemplate_ID` = `PageTemplate`.`ID` WHERE `guid` = ? AND `deletedAt` IS NULL LIMIT 1",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetSelectString(tc.paramSelectStatement)
			require.Equal(t, tc.returnStatement, result)
		})
	}
}

func TestGetInsertString(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name             string
		paramInsertQuery InsertQuery
		returnQuery      string
		returnValues     []interface{}
	}{
		{
			name: "test 'insert page' statement",
			paramInsertQuery: InsertQuery{
				IntoTable: "Page",
				InjectedValues: InjectedValues{
					"PageTemplate_ID": 1,
					"Version_ID":      2,
					"guid":            "PG_1",
					"title":           "Test Title",
					"summary":         "Test Summary",
					"permission":      permission.TypePublic,
					"createdAt":       &now,
					"updatedAt":       nil,
				},
			},
			returnQuery:  "INSERT INTO Page (`PageTemplate_ID`,`Version_ID`,`createdAt`,`guid`,`permission`,`summary`,`title`,`updatedAt`) VALUES (?,?,?,?,?,?,?,?)",
			returnValues: []interface{}{1, 2, &now, "PG_1", permission.TypePublic, "Test Summary", "Test Title", nil},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query, valueStubs := GetInsertString(tc.paramInsertQuery)
			require.Equal(t, tc.returnQuery, query)
			require.Equal(t, tc.returnValues, valueStubs)
		})
	}
}

func TestGetUpdateString(t *testing.T) {
	cases := []struct {
		name                           string
		paramUpdateQuery               UpdateQuery
		paramWhereClauseInjectedValues []interface{}
		returnQuery                    string
		returnValues                   []interface{}
	}{
		{
			name: "test 'update page' statement",
			paramUpdateQuery: UpdateQuery{
				InjectedValues: InjectedValues{
					"title":      "Test Title",
					"summary":    "Test Summary",
					"Version_ID": 1,
					"permission": permission.TypePublic,
				},
				UpdateTable: "Page",
				WhereClause: WhereClause{
					Operator: "AND", WhereOperations: []WhereOperation{
						{LeftSide: "guid", Operator: "= ?"},
					},
				},
			},
			paramWhereClauseInjectedValues: []interface{}{
				"PG_1",
			},
			returnQuery:  "UPDATE Page SET `Version_ID` = ?,`permission` = ?,`summary` = ?,`title` = ? WHERE `guid` = ?",
			returnValues: []interface{}{1, permission.TypePublic, "Test Summary", "Test Title", "PG_1"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			query, valueStubs := GetUpdateString(tc.paramUpdateQuery, tc.paramWhereClauseInjectedValues...)
			require.Equal(t, tc.returnQuery, query)
			require.Equal(t, tc.returnValues, valueStubs)
		})
	}
}
