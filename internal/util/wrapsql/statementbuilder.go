package wrapsql

import (
	"fmt"
	"strings"
)

type InsertQuery struct {
	IntoTable      string
	InjectedValues InjectedValues
}

type UpdateQuery struct {
	UpdateTable    string
	InjectedValues InjectedValues
	WhereClause    WhereClause
}

type InjectedValues map[string]interface{}

type SelectStatement struct {
	Selectors   []string
	FromTable   string
	JoinClauses []JoinClause
	WhereClause WhereClause
	Limit       int
}

type JoinClause struct {
	JoinTable string
	On        OnClause
}

type OnClause struct {
	LeftSide  string
	RightSide string
}

type WhereOperation struct {
	LeftSide  string
	Operator  string
	RightSide string // only use if the RightSide needs to be wrapped in ``.
}

type WhereClause struct {
	Operator        string // either AND or OR
	WhereOperations []WhereOperation
}

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
	s := "`" + strings.Join(sequence, "`,`") + "`"
	return strings.Replace(s, ".", "`.`", -1)
}

func getEscapedString(s string) string {
	return getEscapedSequence([]string{s})
}

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

func GetWhereString(where WhereClause) string {
	var operationStrings []string
	for _, operation := range where.WhereOperations {
		operationStrings = append(operationStrings, getWhereOperationString(operation))
	}
	return strings.Join(operationStrings, where.Operator)
}

func getWhereOperationString(operation WhereOperation) string {
	operationString := fmt.Sprintf("%v %v", getEscapedString(operation.LeftSide), operation.Operator)
	if operation.RightSide != "" {
		operationString = operationString + fmt.Sprintf("%v", getEscapedString(operation.RightSide))
	}
	return operationString
}

func GetInsertString(iq InsertQuery) (string, []interface{}) {
	keys, valueStubs, values := getOrderedInsertValues(iq.InjectedValues)
	keysString := getEscapedSequence(keys)
	valueStubsString := strings.Join(valueStubs, ",")
	return fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", iq.IntoTable, keysString, valueStubsString), values
}

func getOrderedInsertValues(ivs InjectedValues) (keys []string, valueStubs []string, values []interface{}) {
	for key, value := range ivs {
		keys = append(keys, key)
		values = append(values, value)
		valueStubs = append(valueStubs, "?")
	}
	return
}

func GetUpdateString(iq UpdateQuery, whereClauseInjectedValues ...interface{}) (string, []interface{}) {
	keys, _, values := getOrderedInsertValues(iq.InjectedValues)
	values = append(values, whereClauseInjectedValues)
	var setStrings []string
	for _, key := range keys {
		keyString := getEscapedString(key)
		setStrings = append(setStrings, fmt.Sprintf("%v = %v", keyString, "?"))
	}
	setString := strings.Join(setStrings, ",")
	return fmt.Sprintf("UPDATE %v SET %v", iq.UpdateTable, setString), values
}
