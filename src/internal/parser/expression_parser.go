package parser

import (
	"encoding/json"
	"fmt"
)

type ParsedExpression interface {
	GetType() ExpressionType
	GetStringType() string
	ToJson() (string, error)
}

func ParseDataLogJson(parsedExpression, jsonObj string) (ParsedExpression, error) {
	if parsedExpression == string(AddPlayersCommandType) {
		obj := AddPlayersCommand{}
		err := json.Unmarshal([]byte(jsonObj), &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	if parsedExpression == string(AddMatchCommandType) {
		obj := AddMatchCommand{}
		err := json.Unmarshal([]byte(jsonObj), &obj)
		if err != nil {
			return nil, err
		}
		return obj, nil
	}
	return nil, fmt.Errorf("unhandled type to parse from json to parsed expression")
}

func ParseExpression(traceId, expression string) (ParsedExpression, error) {
	scanner := newExpressionScanner(traceId)
	err := scanner.scanRawExpression(expression)
	if err != nil {
		return nil, err
	}

	expressionType, tokens, err := scanner.getLogicalTokens()
	if err != nil {
		return nil, err
	}

	switch expressionType {
	case AddPlayersCommandType:
		return newAddPlayersCommand(traceId, tokens)
	case AddMatchCommandType:
		return newAddMatchCommand(traceId, tokens)
	case SelectQueryType:
		return newSelectQuery(traceId, tokens)
	default:
		return nil, errorCannotParseExpression(traceId, expression)
	}
}
