package tokenizer

import (
	"fmt"
	"strings"
)

// ADD PLAYER 'Tomek', 'Romek';
// SET WIN FOR 'Tomek' AND LOSE FOR 'Romek'; - i dont like it
// SELECT

type TokenizedExpression struct {
	ProcessMethod ProcessMethod
	Arguments     []string
}

type ProcessMethod int8

const (
	AddPlayer ProcessMethod = iota
	UpdatePlayers
	SelectAll
)

const addPlayerMethod string = "ADD PLAYER "
const setWinMethod string = "SET WIN FOR "
const setLoseMethod string = "SET LOSE FOR "
const selectAll string = "SELECT ALL"

func Tokenize(expression string) ([]TokenizedExpression, error) {
	expression, err := sanitazeExpression(expression)
	if err != nil {
		return []TokenizedExpression{}, err
	}

	if strings.HasPrefix(strings.ToUpper(expression), addPlayerMethod) {
		return addPlayerTokenizer(expression)
	}

	if strings.HasPrefix(strings.ToUpper(expression), setWinMethod) ||
		strings.HasPrefix(strings.ToUpper(expression), setLoseMethod) {
		return updatePlayerTokenizer(expression)
	}

	if strings.ToUpper(expression) == selectAll {
		return selectAllTokenizer(expression)
	}

	return []TokenizedExpression{},
		fmt.Errorf("delob error: Could not parse given expression: %s", expression)
}
