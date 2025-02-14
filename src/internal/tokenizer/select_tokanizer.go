package tokenizer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// SELECT Players;
// SELECT Players WHERE Key = 'zxc'
// SELECT Players WHERE Elo > 2500
// SELECT Players JOIN Matches
// SELECT Players ORDER BY Elo

// SELECT Players.history ??????
// SELECT Players JOIN Matches WHERE Key = 'zxc' AND Elo > 2500 ORDER BY Elo ASC

type SelectQuery struct {
	WhereClause WhereClause
	JoinMatches bool
	OrderBy     Field
	OrderDir    OrderDir
}

func (a SelectQuery) GetType() string {
	return SelectQueryType
}

func (a SelectQuery) ToJson() (string, error) {
	result, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

type WhereClause struct {
	Predicates      []Predicate
	ClauseOperators []ClauseOperators
}

type Predicate struct {
	LeftOperant  Field
	Operator     Operator
	RightOperand any
}

type ClauseOperators string

type Operator string

const (
	AND Field = "AND"
	OR  Field = "OR"
)

const (
	Equal           Operator = "="
	NotEqual        Operator = "<>"
	GreaterThan     Operator = ">"
	GreaterThanOrEq Operator = ">="
	LessThan        Operator = "<"
	LessThanOrEq    Operator = "<="
)

type Field string

const (
	Empty Field = ""
	Key   Field = "Key"
	Elo   Field = "Elo"
)

type OrderDir string

const (
	Asc  Field = "ASC"
	Desc Field = "DESC"
)

const SELECT_PLAYERS_EXPRESSION_COMPONENT = "SELECT"
const JOIN_EXPRESSION_COMPONENT = "JOIN"
const WHERE_EXPRESSION_COMPONENT = "WHERE"
const ORDER_BY_EXPRESSION_COMPONENT = "ORDER BY"

func newSelectOrder() SelectQuery {
	return SelectQuery{}
}

func (s *SelectQuery) validateSelectToken(token string) error {
	if token != "Players" {
		return fmt.Errorf("error in select clause - select is handling only Players collection")
	}
	return nil
}

func (s *SelectQuery) validateJoinToken(token string) {
	// TODO not yet implemented - should be connected with streaming player's elo history?
	if token == "Matches" {
		s.JoinMatches = true
	}
}

func (s *SelectQuery) validateOrderToken(token string) error {
	if token == "" {
		return nil
	}

	orderExpressionComponents := strings.Split(token, " ")

	if len(orderExpressionComponents) != 2 {
		return fmt.Errorf("wrong format of ORDER BY clause")
	}

	switch orderExpressionComponents[0] {
	case string(Elo):
		s.OrderBy = Elo
	case string(Key):
		s.OrderBy = Key
	}

	switch strings.ToUpper(orderExpressionComponents[1]) {
	case string(Asc):
		s.OrderDir = OrderDir(Asc)
	case string(Desc):
		s.OrderDir = OrderDir(Desc)
	}

	if s.OrderBy == "" || s.OrderDir == "" {
		return fmt.Errorf("wrong format of ORDER BY clause")
	}

	return nil
}

func tokenizeSelect(expression string) (ParsedExpression, error) {
	selectOrder := newSelectOrder()
	selectToken, joinToken, _, orderToken := getSelectExpressionTokens(expression)

	if errSelect := selectOrder.validateSelectToken(selectToken); errSelect != nil {
		return nil, errSelect
	}
	selectOrder.validateJoinToken(joinToken)
	selectOrder.validateOrderToken(orderToken)

	return selectOrder, nil
}

func getSelectExpressionTokens(rawExpression string) (string, string, string, string) {
	var partIdxs []int
	var partExpressions []string
	var selectToken, joinToken, whereToken, orderToken string
	normalizedExpression := strings.ToUpper(rawExpression)

	partIdxs, partExpressions = tryAppendToExpressionParts(partIdxs, partExpressions, normalizedExpression, SELECT_PLAYERS_EXPRESSION_COMPONENT)
	partIdxs, partExpressions = tryAppendToExpressionParts(partIdxs, partExpressions, normalizedExpression, JOIN_EXPRESSION_COMPONENT)
	partIdxs, partExpressions = tryAppendToExpressionParts(partIdxs, partExpressions, normalizedExpression, WHERE_EXPRESSION_COMPONENT)
	partIdxs, partExpressions = tryAppendToExpressionParts(partIdxs, partExpressions, normalizedExpression, ORDER_BY_EXPRESSION_COMPONENT)
	partIdxs = append(partIdxs, len(normalizedExpression))

	for i := range partExpressions {
		switch partExpressions[i] {
		case SELECT_PLAYERS_EXPRESSION_COMPONENT:
			selectToken = getExpressionToken(partIdxs, i, rawExpression, SELECT_PLAYERS_EXPRESSION_COMPONENT)
		case JOIN_EXPRESSION_COMPONENT:
			joinToken = getExpressionToken(partIdxs, i, rawExpression, JOIN_EXPRESSION_COMPONENT)
		case WHERE_EXPRESSION_COMPONENT:
			whereToken = getExpressionToken(partIdxs, i, rawExpression, WHERE_EXPRESSION_COMPONENT)
		case ORDER_BY_EXPRESSION_COMPONENT:
			orderToken = getExpressionToken(partIdxs, i, rawExpression, ORDER_BY_EXPRESSION_COMPONENT)
		}
	}

	return selectToken, joinToken, whereToken, orderToken
}

func getExpressionToken(partIdxs []int, i int, expression string, component string) string {
	return strings.TrimSpace(expression[partIdxs[i]+len(component) : partIdxs[i+1]])
}

func tryAppendToExpressionParts(partIdxs []int, partExpressions []string, expression string, component string) ([]int, []string) {
	if idx := strings.Index(expression, component); idx > -1 {
		partIdxs = append(partIdxs, idx)
		partExpressions = append(partExpressions, component)
	}
	return partIdxs, partExpressions
}
