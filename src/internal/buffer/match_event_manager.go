package buffer

import (
	"delob/internal/utils"
)

type Match struct {
	transactionStatus transactionStatus
	Key               string
	TeamOneKeys       []string
	TeamTwoKeys       []string
	MatchResult       int8
	AddTimestamp      int64
}

func newMatch(teamOneKeys []string, teamTwoKeys []string, matchResult int8, transaction *Transaction) *Match {
	match := &Match{
		transactionStatus: inProgress,
		Key:               utils.GenerateKey(),
		AddTimestamp:      utils.Timestamp(),
		TeamOneKeys:       teamOneKeys,
		TeamTwoKeys:       teamTwoKeys,
		MatchResult:       matchResult,
	}
	transaction.AddMatchPointer(match)
	return match
}
