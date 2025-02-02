package tokenizer

import (
	"delob/internal/shared"
	"fmt"
	"strings"
)

type Token interface {
}

type AddPlayersToken struct {
	Keys []string
}

type AddMatchToken struct {
	TeamOneKeys []string
	TeamTwoKeys []string
	MatchResult shared.MatchResult
}

type SelectAllToken struct{}

const ADD_PLAYER_EXPRESSION string = "ADD PLAYER"
const SET_WIN_EXPRESSION string = "SET WIN FOR "
const SET_LOSE_EXPRESSION string = "SET LOSE FOR "
const DRAW_EXPRESSION string = "SET DRAW BETWEEN "
const SELECT_ALL_EXPRESSION string = "SELECT ALL"

func Tokenize(expression string) (interface{}, error) {
	sanitazedExpression, err := sanitazeExpression(expression)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), ADD_PLAYER_EXPRESSION) {
		return tokenizeAddPlayersExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), SET_WIN_EXPRESSION) ||
		strings.HasPrefix(strings.ToUpper(sanitazedExpression), SET_LOSE_EXPRESSION) {
		return tokenizeDecisiveMatchResultExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), DRAW_EXPRESSION) {
		return tokenizeDrawMatchResultExpression(sanitazedExpression)
	}

	if strings.ToUpper(sanitazedExpression) == SELECT_ALL_EXPRESSION {
		return selectAllTokenizer(sanitazedExpression)
	}

	return nil,
		fmt.Errorf("delob error: Could not parse given expression: %s", sanitazedExpression)
}
