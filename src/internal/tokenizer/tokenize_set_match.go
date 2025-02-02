package tokenizer

import (
	"delob/internal/shared"
	"fmt"
	"strings"
)

func tokenizeDecisiveMatchResultExpression(expression string) (interface{}, error) {
	splitedUpdateExpression, err := splitUpdateExpression(expression)
	if err != nil {
		return nil, err
	}

	teamOneKeys, teamTwoKeys, matchResult, errExtract := extractArgumentsForUpdatePlayerMethod(splitedUpdateExpression)
	if errExtract != nil {
		return nil, errExtract
	}

	return AddMatchToken{
		TeamOneKeys: teamOneKeys,
		TeamTwoKeys: teamTwoKeys,
		MatchResult: matchResult,
	}, nil
}

func tokenizeDrawMatchResultExpression(expression string) (interface{}, error) {
	splitedUpdateExpression, err := splitUpdateExpression(expression)
	if err != nil {
		return nil, err
	}

	teamOneKeys, errTeamOneKeys := extractKeysFromPartialExpression(splitedUpdateExpression[0][len(DRAW_EXPRESSION) : len(splitedUpdateExpression[0])-1])
	if errTeamOneKeys != nil {
		return nil, errTeamOneKeys
	}

	teamTwoKeys, errTeamTwoKeys := extractKeysFromPartialExpression(splitedUpdateExpression[1])
	if errTeamTwoKeys != nil {
		return nil, errTeamTwoKeys
	}

	return AddMatchToken{
		TeamOneKeys: teamOneKeys,
		TeamTwoKeys: teamTwoKeys,
		MatchResult: shared.Draw,
	}, nil
}

func splitUpdateExpression(expression string) ([]string, error) {
	splittedExpression := strings.Split(expression, "AND")
	if len(splittedExpression) != 2 {
		return splittedExpression, fmt.Errorf("update expression should contain exatcly two paramaters - win and lose")
	}

	return splittedExpression, nil
}

func extractArgumentsForUpdatePlayerMethod(splitedUpdateExpression []string) ([]string, []string, shared.MatchResult, error) {
	var teamOneKeys, teamTwoKeys []string
	var errTeamOneKeys, errTeamTwoKeys error
	var matchResult shared.MatchResult

	for i := range splitedUpdateExpression {
		expressionArgs := strings.Split(strings.TrimSpace(splitedUpdateExpression[i]), "FOR")

		if i == 0 {
			teamOneKeys, errTeamOneKeys = extractKeysFromPartialExpression(expressionArgs[1])
			if errTeamOneKeys != nil {
				return teamOneKeys, teamTwoKeys, shared.Unknown, errTeamOneKeys
			}
			if strings.ToUpper(expressionArgs[0]) == "SET WIN " {
				matchResult = shared.TeamOneWins
			}
		} else {
			teamTwoKeys, errTeamTwoKeys = extractKeysFromPartialExpression(expressionArgs[1])
			if errTeamTwoKeys != nil {
				return teamOneKeys, teamTwoKeys, shared.Unknown, errTeamTwoKeys
			}
			if strings.ToUpper(expressionArgs[0]) == "WIN " {
				matchResult = shared.TeamTwoWins
			}
		}
	}
	return teamOneKeys, teamTwoKeys, matchResult, nil
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
