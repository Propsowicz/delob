package tokenizer

import (
	"fmt"
	"strings"
)

func updatePlayerTokenizer(expression string) ([]TokenizedExpression, error) {
	result := []TokenizedExpression{}
	splitedUpdateExpression, err := splitUpdateExpression(expression)
	if err != nil {
		return result, err
	}

	var args []string

	for i := range splitedUpdateExpression {
		arg, errExtract := extractArgumentsForUpdatePlayerMethod(splitedUpdateExpression[i])
		if errExtract != nil {
			return result, errExtract
		}
		args = append(args, arg[0])
		args = append(args, arg[1])
	}

	result = append(result, TokenizedExpression{
		ProcessMethod: UpdatePlayers,
		Arguments:     args,
	})
	return result, nil
}

func splitUpdateExpression(expression string) ([]string, error) {
	splittedExpression := strings.Split(expression, "AND")
	if len(splittedExpression) != 2 {
		return splittedExpression, fmt.Errorf("update expression should contain exatcly two paramaters - win and lose")
	}

	return splittedExpression, nil
}

func extractArgumentsForUpdatePlayerMethod(expression string) ([]string, error) {
	var result []string

	expressionArgs := strings.Split(strings.TrimSpace(expression), " ")
	expressionArgument := expressionArgs[len(expressionArgs)-1]

	switch strings.ToUpper(expressionArgs[len(expressionArgs)-3]) {
	case "WIN":
		result = append(result, "WIN")
	case "LOSE":
		result = append(result, "LOSE")
	}

	id, err := extractIdFromString(expressionArgument)
	if err != nil {
		return result, err
	}
	result = append(result, id)

	if len(result) > 0 {
		return result, nil
	}

	return result,
		fmt.Errorf("delob error: Could not extract arguments from expression: %s", expression)
}
