package mysqlstore

import (
	"fmt"
	"strings"
)

type whereOperation struct {
	leftSide  string
	operator  string
	rightSide string // only use if the rightSide needs to be wrapped in ``.
}

type whereClause struct {
	operator        string // either AND or OR
	whereOperations []whereOperation
}

func newSelectStatement(selectors []string, from string, where whereClause, limit int) string {
	selectString := "`" + strings.Join(selectors, "`,`") + "`"
	whereString := getWhereString(where)
	statement := fmt.Sprintf("SELECT %v FROM %v", selectString, from)
	if whereString != "" {
		statement = statement + fmt.Sprintf(" WHERE %v", whereString)
	}
	if limit != 0 {
		statement = statement + fmt.Sprintf(" LIMIT %v", limit)
	}
	return statement
}

func getWhereString(where whereClause) string {
	var operationStrings []string
	for _, operation := range where.whereOperations {
		operationStrings = append(operationStrings, getWhereOperationString(operation))
	}
	return strings.Join(operationStrings, where.operator)
}

func getWhereOperationString(operation whereOperation) string {
	operationString := fmt.Sprintf("`%v` %v", operation.leftSide, operation.operator)
	if operation.rightSide != "" {
		operationString = operationString + fmt.Sprintf(`%v`, operation.rightSide)
	}
	return operationString
}
