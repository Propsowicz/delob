package tokenizer

import (
	"delob/internal/shared"
	"encoding/json"
	"fmt"
	"strings"
)

// TODO it needs to be refactored: see pgsql and how it is done there
type ParsedExpression interface {
	GetType() string
	ToJson() (string, error)
}

const AddPlayersCommandType string = "AddPlayersCommand"
const AddMatchCommandType string = "AddMatchCommand"
const SelectQueryType string = "SelectQuery"

type AddPlayersCommand struct {
	Keys []string
}

func (a AddPlayersCommand) GetType() string {
	return AddPlayersCommandType
}

func (a AddPlayersCommand) ToJson() (string, error) {
	result, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

type AddMatchCommand struct {
	TeamOneKeys []string
	TeamTwoKeys []string
	MatchResult shared.MatchResult
}

func (a AddMatchCommand) GetType() string {
	return AddMatchCommandType
}

func (a AddMatchCommand) ToJson() (string, error) {
	result, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

const ADD_PLAYER_EXPRESSION string = "ADD PLAYER"
const SET_WIN_EXPRESSION string = "SET WIN FOR "
const SET_LOSE_EXPRESSION string = "SET LOSE FOR "
const DRAW_EXPRESSION string = "SET DRAW BETWEEN "
const SELECT_EXPRESSION string = "SELECT "

func ParseFromJson(parsedExpression, jsonObj string) (ParsedExpression, error) {
	if parsedExpression == AddPlayersCommandType {
		obj := AddPlayersCommand{}
		err := json.Unmarshal([]byte(jsonObj), &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	if parsedExpression == AddMatchCommandType {
		obj := AddMatchCommand{}
		err := json.Unmarshal([]byte(jsonObj), &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return nil, fmt.Errorf("unhandled type to parse from json to parsed expression")
}

func Tokenize(expression string) (ParsedExpression, error) {
	sanitazedExpression, err := sanitazeExpression(expression)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), ADD_PLAYER_EXPRESSION) {
		return tokenizeAddPlayersExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), SET_WIN_EXPRESSION) ||
		strings.HasPrefix(strings.ToUpper(sanitazedExpression), SET_LOSE_EXPRESSION) {
		return tokenizeDecisiveMatchResultExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), DRAW_EXPRESSION) {
		return tokenizeDrawMatchResultExpression(sanitazedExpression)
	}

	if strings.HasPrefix(strings.ToUpper(sanitazedExpression), SELECT_EXPRESSION) {
		return tokenizeSelect(sanitazedExpression)
	}

	return nil,
		fmt.Errorf("delob error: Could not parse given expression: %s", sanitazedExpression)
}
