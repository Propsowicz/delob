package tokenizer

import (
	"errors"
	"fmt"
	"strings"
)

func sanitazeExpression(expression string) (string, error) {
	if expression[len(expression)-1:] != ";" {
		return "", fmt.Errorf("expression should end with semicolon")
	}
	return expression[:len(expression)-1], nil
}

func tryExtractExpressionFromBrackets(expression string) string {
	expression = strings.TrimSpace(expression)

	if expression[0] == '(' && expression[len(expression)-1] == ')' {
		return expression[1 : len(expression)-1]
	}
	return expression
}

func extractIdFromString(rawString string) (string, error) {
	str := strings.TrimSpace(rawString)
	if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		return str[1 : len(str)-1], nil
	} else {
		return "", errors.New("delob error: Key not enclosed in parentheses: " + rawString)
	}
}
