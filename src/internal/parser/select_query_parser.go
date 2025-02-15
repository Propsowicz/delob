package parser

import "encoding/json"

type SelectQuery struct {
	WhereClause WhereClause
	JoinMatches bool
	OrderBy     Field
	OrderDir    OrderDir
}

func newSelectQuery(traceId string, tokens []Token) (ParsedExpression, error) {
	// happy implementation for now
	if len(tokens) == 1 && tokens[0].Token == SelectPlayers {
		return SelectQuery{
			JoinMatches: false,
		}, nil
	}
	if len(tokens) == 2 {
		if tokens[1].Token == OrderByAsc {
			return SelectQuery{
				JoinMatches: false,
				OrderBy:     Field(tokens[1].Value[0]),
				OrderDir:    Asc,
			}, nil
		} else {
			return SelectQuery{
				JoinMatches: false,
				OrderBy:     Field(tokens[1].Value[0]),
				OrderDir:    Desc,
			}, nil
		}
	}

	return nil, errorCannotGenerateParsedExpression(traceId)
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
