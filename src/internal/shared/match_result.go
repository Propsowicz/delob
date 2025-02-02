package shared

type MatchResult int8

const (
	Unknown MatchResult = iota
	Draw
	TeamOneWins
	TeamTwoWins
)

func ScoreModifier(m MatchResult) (float64, float64) {
	if m == TeamOneWins {
		return 1, 0
	}
	if m == TeamTwoWins {
		return 0, 1
	}
	return 0.5, 0.5
}
