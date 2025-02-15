package parser

import (
	"delob/internal/shared"
	"encoding/json"
)

type AddMatchCommand struct {
	TeamOneKeys []string
	TeamTwoKeys []string
	MatchResult shared.MatchResult
}

func newAddMatchCommand(traceId string, tokens []Token) (ParsedExpression, error) {
	if tokens[0].Token == SetDraw {
		return AddMatchCommand{
			MatchResult: shared.Draw,
			TeamOneKeys: tokens[0].Value,
			TeamTwoKeys: tokens[1].Value,
		}, nil
	}
	if tokens[0].Token == SetWin {
		return AddMatchCommand{
			MatchResult: shared.TeamOneWins,
			TeamOneKeys: tokens[0].Value,
			TeamTwoKeys: tokens[1].Value,
		}, nil
	}
	if tokens[1].Token == SetWin {
		return AddMatchCommand{
			MatchResult: shared.TeamTwoWins,
			TeamOneKeys: tokens[0].Value,
			TeamTwoKeys: tokens[1].Value,
		}, nil
	}

	return nil, errorCannotGenerateParsedExpression(traceId)
}

func (a AddMatchCommand) GetType() ExpressionType {
	return AddMatchCommandType
}

func (a AddMatchCommand) GetStringType() string {
	return string(AddMatchCommandType)
}

func (a AddMatchCommand) ToJson() (string, error) {
	result, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
