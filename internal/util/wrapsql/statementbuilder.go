package wrapsql

import (
	"fmt"
	"sort"
	"strings"
)

// InsertQuery is used to generate an insert query
type InsertQuery struct {
	IntoTable      string
	InjectedValues InjectedValues
}

// UpdateQuery is used to generate an update query
type UpdateQuery struct {
	UpdateTable    string
	InjectedValues InjectedValues
	WhereClause    WhereClause
}

// InjectedValues are a mapping of key/value pairs where the key is the name of the table column and the value is it's injected value.
type InjectedValues map[string]interface{}

// SelectStatement is used to generate a select statement
type SelectStatement struct {
	Selectors   []string
	FromTable   string
	JoinClauses []JoinClause
	WhereClause WhereClause
	Limit       int
}

// JoinClause is used to generate a JOIN clause
type JoinClause struct {
	JoinTable string
	On        OnClause
}

// OnClause is used to generate an ON clause
type OnClause struct {
	LeftSide  string
	RightSide string
}

// WhereOperation is used to generate a WHERE operation, such as "`ID` = ?"
type WhereOperation struct {
	LeftSide  string
	Operator  string
	RightSide string // only use if the RightSide needs to be wrapped in ``.
}

// WhereClause is used to generate a WHERE clause, which is a series of WhereOperations, such as "`ID` = ? AND `deletedAt' IS NULL"
type WhereClause struct {
	Operator        string // either AND or OR
	WhereOperations []WhereOperation
}

// GetSelectString returns a statement string intended for a SELECT call.
func GetSelectString(ss SelectStatement) string {
	selectString := getEscapedSequence(ss.Selectors)
	whereString := GetWhereString(ss.WhereClause)
	statement := fmt.Sprintf("SELECT %v FROM %v", selectString, ss.FromTable)
	joinString := GetJoinsString(ss.JoinClauses)
	if joinString != "" {
		statement = statement + " " + joinString
	}
	if whereString != "" {
		statement = statement + fmt.Sprintf(" WHERE %v", whereString)
	}
	if ss.Limit != 0 {
		statement = statement + fmt.Sprintf(" LIMIT %v", ss.Limit)
	}
	return statement
}

func getEscapedSequence(sequence []string) string {
	var escapedSequence []string
	for _, s := range sequence {
		if shouldBeEscaped(s) {
			escapedSequence = append(escapedSequence, "`"+strings.Replace(s, ".", "`.`", -1)+"`")
		} else {
			escapedSequence = append(escapedSequence, s)
		}
	}
	return strings.Join(escapedSequence, ",")
}

func shouldBeEscaped(s string) bool {
	return !strings.HasPrefix(s, "COUNT(")
}

func getEscapedString(s string) string {
	return getEscapedSequence([]string{s})
}

// GetJoinsString returns a string for the JOIN clause in the query
func GetJoinsString(joins []JoinClause) string {
	var joinStrings []string
	for _, join := range joins {
		joinStrings = append(joinStrings, getJoinString(join))
	}
	return strings.Join(joinStrings, " ")
}

func getJoinString(join JoinClause) string {
	return fmt.Sprintf("JOIN %v ON %v = %v", join.JoinTable, getEscapedString(join.On.LeftSide), getEscapedString(join.On.RightSide))
}

// GetWhereString returns a string for the WHERE clause in the query
func GetWhereString(where WhereClause) string {
	var operationStrings []string
	for _, operation := range where.WhereOperations {
		operationStrings = append(operationStrings, getWhereOperationString(operation))
	}
	return strings.Join(operationStrings, fmt.Sprintf(" %v ", where.Operator))
}

func getWhereOperationString(operation WhereOperation) string {
	operationString := fmt.Sprintf("%v %v", getEscapedString(operation.LeftSide), operation.Operator)
	if operation.RightSide != "" {
		operationString = operationString + fmt.Sprintf("%v", getEscapedString(operation.RightSide))
	}
	return operationString
}

// GetInsertString returns a statement string intended for an INSERT call.
func GetInsertString(iq InsertQuery) (string, []interface{}) {
	keys, valueStubs, values := getOrderedInsertValues(iq.InjectedValues)
	keysString := getEscapedSequence(keys)
	valueStubsString := strings.Join(valueStubs, ",")
	return fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", iq.IntoTable, keysString, valueStubsString), values
}

func getOrderedInsertValues(ivs InjectedValues) (keys []string, valueStubs []string, values []interface{}) {
	for key := range ivs {
		keys = append(keys, key)
		valueStubs = append(valueStubs, "?")
	}
	sort.Strings(keys)
	for _, key := range keys {
		values = append(values, ivs[key])
	}
	return
}

// GetUpdateString returns a statement string intended for an UPDATE call.
func GetUpdateString(iq UpdateQuery, whereClauseInjectedValues ...interface{}) (string, []interface{}) {
	keys, _, values := getOrderedInsertValues(iq.InjectedValues)
	values = append(values, whereClauseInjectedValues...)
	whereString := GetWhereString(iq.WhereClause)
	var setStrings []string
	for _, key := range keys {
		keyString := getEscapedString(key)
		setStrings = append(setStrings, fmt.Sprintf("%v = %v", keyString, "?"))
	}
	setString := strings.Join(setStrings, ",")
	statement := fmt.Sprintf("UPDATE %v SET %v", iq.UpdateTable, setString)
	if whereString != "" {
		statement = statement + fmt.Sprintf(" WHERE %v", whereString)
	}
	return statement, values
}
