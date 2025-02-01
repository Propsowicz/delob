package buffer

type MatchEvent struct {
	id           string
	winIds       []string
	loseId       []string
	AddTimestamp int64
}
