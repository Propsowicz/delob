package buffer

import (
	"delob/internal/utils"
)

type Match struct {
	Key          string
	TeamOneKeys  []string
	TeamTwoKeys  []string
	MatchResult  int8
	AddTimestamp int64
}

func newMatch(teamOneKeys []string, teamTwoKeys []string, matchResult int8) Match {
	return Match{
		Key:          utils.GenerateKey(),
		AddTimestamp: utils.Timestamp(),
		TeamOneKeys:  teamOneKeys,
		TeamTwoKeys:  teamTwoKeys,
		MatchResult:  matchResult,
	}
}
