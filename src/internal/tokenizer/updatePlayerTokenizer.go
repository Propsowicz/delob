package tokenizer

import (
	"fmt"
	"strings"
)

func tokenizeUpdatePlayerExpression(expression string) (interface{}, error) {
	splitedUpdateExpression, err := splitUpdateExpression(expression)
	if err != nil {
		return nil, err
	}

	winKeys, loseKeys, errExtract := extractArgumentsForUpdatePlayerMethod(splitedUpdateExpression)
	if errExtract != nil {
		return nil, errExtract
	}

	return UpdatePlayersToken{
		WinKeys:  winKeys,
		LoseKeys: loseKeys,
	}, nil
}

func splitUpdateExpression(expression string) ([]string, error) {
	splittedExpression := strings.Split(expression, "AND")
	if len(splittedExpression) != 2 {
		return splittedExpression, fmt.Errorf("update expression should contain exatcly two paramaters - win and lose")
	}

	return splittedExpression, nil
}

func extractArgumentsForUpdatePlayerMethod(expressions []string) ([]string, []string, error) {
	var winKeys, loseKeys []string
	var errWinKeys, errLoseKeys error

	for i := range expressions {
		expressionArgs := strings.Split(strings.TrimSpace(expressions[i]), "FOR")

		switch strings.ToUpper(expressionArgs[0]) {
		case "SET WIN ", "WIN ":
			winKeys, errWinKeys = extractKeysFromPartialExpression(expressionArgs[1])
			if errWinKeys != nil {
				return winKeys, loseKeys, errWinKeys
			}

		case "SET LOSE ", "LOSE ":
			loseKeys, errLoseKeys = extractKeysFromPartialExpression(expressionArgs[1])
			if errLoseKeys != nil {
				return winKeys, loseKeys, errLoseKeys
			}
		}
	}
	return winKeys, loseKeys, nil
}

func extractKeysFromPartialExpression(partialExpression string) ([]string, error) {
	var result []string
	keys := strings.Split(tryExtractExpressionFromBrackets(partialExpression), ",")

	for i := range keys {
		key, err := extractIdFromString(strings.TrimSpace(keys[i]))
		if err != nil {
			return nil, err
		}

		result = append(result, key)

	}
	return result, nil
}
