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

func extractIdFromString(rawString string) (string, error) {
	str := strings.TrimSpace(rawString)
	if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
		return str[1 : len(str)-1], nil
	} else {
		return "", errors.New("delob error: Key not enclosed in parentheses: " + rawString)
	}
}

func convertStringToProcessMethod(str string) (ProcessMethod, error) {
	switch str {
	case "0":
		return AddPlayer, nil
	case "1":
		return SetWin, nil
	case "2":
		return SetLose, nil
	default:
		return AddPlayer, fmt.Errorf("cannot convert string to ProcessMethod enum")
	}
}
