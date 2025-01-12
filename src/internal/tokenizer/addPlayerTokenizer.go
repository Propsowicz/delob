package tokenizer

import (
	"fmt"
	"strings"
)

func addPlayerTokenizer(expression string) ([]TokenizedExpression, error) {
	args, err := extractArgumentsForAddPlayerMethod(expression)
	if err != nil {
		return []TokenizedExpression{}, err
	}

	return []TokenizedExpression{
			{
				ProcessMethod: AddPlayer,
				Arguments:     args,
			},
		},
		nil
}

func extractArgumentsForAddPlayerMethod(expression string) ([]string, error) {
	var result []string

	expressionArguments := expression[len(addPlayerMethod):]
	rawArgs := strings.Split(expressionArguments, ",")

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
