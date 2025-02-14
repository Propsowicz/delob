package tokenizer

import (
	"fmt"
	"strings"
)

const addSinglePlayerExpression string = "ADD PLAYER "
const addMultiplePlayersExpression string = "ADD PLAYERS "

func tokenizeAddPlayersExpression(expression string) (ParsedExpression, error) {
	var result []string
	var err error

	if strings.HasPrefix(strings.ToUpper(expression), addSinglePlayerExpression) {
		result, err = extractArgumentsForAddPlayerMethod(expression, addSinglePlayerExpression)
		if err != nil {
			return nil, err
		}
		if len(result) != 1 {
			return nil, fmt.Errorf("delob: incorrect syntax - tried to add multiple players with 'ADD PLAYER' expression")
		}
	}
	if strings.HasPrefix(strings.ToUpper(expression), addMultiplePlayersExpression) {
		result, err = extractArgumentsForAddPlayerMethod(expression, addMultiplePlayersExpression)
		if err != nil {
			return nil, err
		}
	}

	return AddPlayersCommand{
		Keys: result,
	}, nil
}

func extractArgumentsForAddPlayerMethod(expression string, prefix string) ([]string, error) {
	var result []string

	expressionArguments := expression[len(prefix):]
	rawArgs := strings.Split(tryExtractExpressionFromBrackets(expressionArguments), ",")

	for i := range rawArgs {
		id, err := extractIdFromString(rawArgs[i])
		if err != nil {
			return result, err
		}
		result = append(result, id)
	}

	if len(result) > 0 {
		return result, nil
	}

	return result,
		fmt.Errorf("delob error: Could not extract arguments from expression: %s", expression)
}
