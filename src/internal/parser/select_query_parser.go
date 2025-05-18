package parser

import (
	"encoding/json"
)

type SelectQueryComponent string

const (
	KeyComponent     SelectQueryComponent = "key"
	EloComponent     SelectQueryComponent = "elo"
	EventsComponent  SelectQueryComponent = "events"
	MatchesComponent SelectQueryComponent = "matches"
)

type SelectQuery struct {
	SelectQueryComponents []SelectQueryComponent
	WhereClause           WhereClause
	JoinMatches           bool
	OrderBy               Field
	OrderDir              OrderDir
}

var componentMap = map[string]SelectQueryComponent{
	string(KeyComponent):     KeyComponent,
	string(EloComponent):     EloComponent,
	string(EventsComponent):  EventsComponent,
	string(MatchesComponent): MatchesComponent,
}

func newSelectQuery(traceId string, tokens []Token) (ParsedExpression, error) {
	// happy implementation for now

	selectQuery := SelectQuery{
		SelectQueryComponents: []SelectQueryComponent{},
	}

	if tokens[0].Token != SelectPlayers {
		// if tokens[0].Token != SelectPlayers || len(tokens[0].Value) == 0 {
		return nil, errorCannotGenerateParsedExpression(traceId)
	}

	for _, tokenValue := range tokens[0].Value {
		queryComponent, ok := componentMap[tokenValue]
		if !ok {
			return nil, errorCannotGenerateParsedExpression(traceId)
		}
		selectQuery.SelectQueryComponents = append(selectQuery.SelectQueryComponents, queryComponent)
	}

	if len(tokens) == 2 && (tokens[1].Token == OrderByAsc || tokens[1].Token == OrderByDesc) {
		selectQuery.OrderBy = Field(tokens[1].Value[0])

		if tokens[1].Token == OrderByAsc {
			selectQuery.OrderDir = Asc
		} else {
			selectQuery.OrderDir = Desc
		}
	}
	return selectQuery, nil
}

func (a SelectQuery) GetType() ExpressionType {
	return SelectQueryType
}

func (a SelectQuery) GetStringType() string {
	return string(SelectQueryType)
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
	Asc  OrderDir = "ASC"
	Desc OrderDir = "DESC"
)
