package tokenizer

// SELECT Players;
// SELECT Players WHERE Key = 'zxc'
// SELECT Players WHERE Elo > 2500
// SELECT Players JOIN Matches
// SELECT Players ORDER BY Elo

// SELECT Players JOIN Matches Where Key = 'zxc' AND Elo > 2500 ORDER BY Elo ASC

type SelectOrder struct {
	WhereClause WhereClause
	JoinMatches bool
	OrderBy     Field
	OrderDir    OrderDir
}

type WhereClause struct {
	Filters         []Filter
	ClauseOperators []ClauseOperators
}

type Filter struct {
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

const SELECT_PLAYERS_EXPRESSION string = "SELECT Players"

func tokenize_select(expression string) (interface{}, error) {

	return SelectOrder{}, nil
}
