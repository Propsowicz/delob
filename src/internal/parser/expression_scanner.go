package parser

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

type ExpressionScanner struct {
	ExpressionType ExpressionType

	tokens []Token
}

type Token struct {
	Token TokenType
	Value []string
}

type ExpressionType string

const (
	AddPlayersCommand ExpressionType = "AddPlayersCommand"
	AddMatchCommand   ExpressionType = "AddMatchCommand"
	SelectQuery       ExpressionType = "SelectQuery"
)

type TokenType string

const (
	AddPlayer  TokenType = "AddPlayer"
	AddPlayers TokenType = "AddPlayers"
)

func NewExpressionScanner() ExpressionScanner {
	return ExpressionScanner{}
}

func (sc *ExpressionScanner) GetLogicalTokens() ([]Token, error) {
	if sc.ExpressionType == "" {
		return nil, fmt.Errorf("expression type is not determined")
	}
	return sc.tokens, nil
}

const add string = "add"
const player string = "player"
const players string = "players"

func (sc *ExpressionScanner) ScanRawExpression(expression string) error {
	sanitazedExpression, err := sanitazeExpression(expression)
	if err != nil {
		return nil
	}
	result := []Token{}
	searchedExpression := []string{add}

	expressionComponents := strings.Split(sanitazedExpression, " ")

	for i := 0; i < len(expressionComponents); i++ {
		if slices.Contains(searchedExpression, strings.ToLower(expressionComponents[i])) {
			foundAddPlayerExpression, t, iLambda, err := scanForAddPlayerExpression(i, expressionComponents, sanitazedExpression)
			if foundAddPlayerExpression {
				sc.ExpressionType = ExpressionType(AddPlayer)
				result = append(result, t)
			}
			if err != nil {
				return err
			}
			i += iLambda
		}
	}

	sc.tokens = result
	return nil
}

func scanForAddPlayerExpression(i int, expressionComponents []string, rawExpression string) (bool, Token, int, error) {
	if strings.ToLower(expressionComponents[i]) == add {
		if strings.ToLower(expressionComponents[i+1]) == player {
			return true,
				Token{
					AddPlayer, []string{
						strings.Trim(expressionComponents[i+2], "'")},
				},
				2,
				nil
		}
		if strings.ToLower(expressionComponents[i+1]) == players {
			keys, err := extractKeysFromRawExpression(rawExpression)
			if err != nil {
				return true, Token{}, 0, err
			}

			return true,
				Token{
					AddPlayers, keys,
				},
				2,
				nil
		}
		return true, Token{}, 0, fmt.Errorf("wrong format of expression")
	}
	return false, Token{}, 0, nil
}

func sanitazeExpression(expression string) (string, error) {
	if expression[len(expression)-1:] != ";" {
		return "", fmt.Errorf("expression should end with semicolon")
	}
	return expression[:len(expression)-1], nil
}

func extractKeysFromRawExpression(expression string) ([]string, error) {
	pattern := `\(\s*'([^']+)'(?:\s*,\s*'([^']+)')*\s*\)`
	r := regexp.MustCompile(pattern)
	match := r.FindString(expression)

	if match == "" {
		return nil, fmt.Errorf("cannot extract keys from given expression - %s", expression)
	}

	content := strings.Trim(match, "()")
	items := strings.Split(content, ",")

	var result []string
	for _, item := range items {
		result = append(result, strings.Trim(strings.TrimSpace(item), "'"))
	}

	return result, nil
}
