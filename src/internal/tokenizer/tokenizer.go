package tokenizer

import (
	"fmt"
	"strings"
)

type Token interface {
}

type AddPlayersToken struct {
	Keys []string
}

type UpdatePlayersToken struct {
	WinKeys  []string
	LoseKeys []string
}

type SelectAllToken struct{}

const addPlayerExpression string = "ADD PLAYER"
const setWinExpression string = "SET WIN FOR "
const setLoseExpression string = "SET LOSE FOR "
const selectAllExpression string = "SELECT ALL"

func Tokenize(expression string) (interface{}, error) {
	sanitazedExpression, err := sanitazeExpression(expression)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), addPlayerExpression) {
		return tokenizeAddPlayersExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), setWinExpression) ||
		strings.HasPrefix(strings.ToUpper(sanitazedExpression), setLoseExpression) {
		return tokenizeUpdatePlayerExpression(sanitazedExpression)
	}

	if strings.ToUpper(sanitazedExpression) == selectAllExpression {
		return selectAllTokenizer(sanitazedExpression)
	}

	return nil,
		fmt.Errorf("delob error: Could not parse given expression: %s", sanitazedExpression)
}
