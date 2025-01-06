package tokenizer

import (
	model "delob/internal/shared/model"
	"errors"
	"fmt"
	"strings"
)

// ADD PLAYER 'Tomek', 'Romek'
// SET WIN FOR 'Tomek' AND LOSE FOR 'Romek' - i dont like it
// SELECT

const addPlayerMethod string = "ADD PLAYER "

func Tokenize(expression string) ([]model.TokenizedExpression, error) {

	if strings.HasPrefix(strings.ToUpper(expression), addPlayerMethod) {
		return addPlayerTokenizer(expression)
	}

	return []model.TokenizedExpression{},
		fmt.Errorf("delob error: Could not parse given expression: %s", expression)
}

func addPlayerTokenizer(expression string) ([]model.TokenizedExpression, error) {
	args, err := extractArgumentsForAddPlayerMethod(expression)
	if err != nil {
		return []model.TokenizedExpression{}, err
	}

	return []model.TokenizedExpression{
			{
				ProcessMethod: model.AddPlayer,
				Arguments:     args,
			},
		},
		nil
}

func extractArgumentsForAddPlayerMethod(expression string) ([]string, error) {
	var result []string
	if expression[len(expression)-1:] != ";" {
		return []string{}, fmt.Errorf("expression should end with semicolon")
	}
	expression = expression[:len(expression)-1]

	expressionArguments := expression[len(addPlayerMethod):]
	rawArgs := strings.Split(expressionArguments, ",")

	for i := range rawArgs {
		str := strings.TrimSpace(rawArgs[i])
		if strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'") {
			result = append(result, str[1:len(str)-1])
		} else {
			return nil, errors.New("delob error: Key not enclosed in parentheses: " + rawArgs[i])
		}
	}

	if len(result) > 0 {
		return result, nil
	}

	return result,
		fmt.Errorf("delob error: Could not extract arguments from expression: %s", expression)
}
