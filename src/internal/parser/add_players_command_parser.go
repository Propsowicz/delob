package parser

import (
	"encoding/json"
)

type AddPlayersCommand struct {
	Keys []string
}

func newAddPlayersCommand(traceId string, tokens []Token) (ParsedExpression, error) {
	return AddPlayersCommand{
		Keys: tokens[0].Value,
	}, nil
}

func (a AddPlayersCommand) GetType() ExpressionType {
	return AddPlayersCommandType
}

func (a AddPlayersCommand) GetStringType() string {
	return string(AddPlayersCommandType)
}

func (a AddPlayersCommand) ToJson() (string, error) {
	result, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
