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
	err := checkForDuplicatedKeys(tokens[0].Value, tokens[1].Value)
	if err != nil {
		return nil, err
	}

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

func checkForDuplicatedKeys(t1, t2 []string) error {
	seen := make(map[string]bool, len(t1)+len(t2))

	for i := range t1 {
		if _, exists := seen[t1[i]]; exists {
			return errorCannotUseUserKeyMoreThanOnce(t1[i])
		}
		seen[t1[i]] = true
	}

	for j := range t2 {
		if _, exists := seen[t2[j]]; exists {
			return errorCannotUseUserKeyMoreThanOnce(t2[j])
		}
		seen[t2[j]] = true
	}
	return nil

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
