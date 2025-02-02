package elo

import (
	dto "delob/internal/processor/model"
	"delob/internal/shared"
	"math"
)

const k float64 = 25

type Calculator struct {
	teamOneEloLambda int16
	teamTwoEloLambda int16
}

func NewCalculator(teamOnePlayers []dto.Player, teamTwoPlayers []dto.Player, matchResult shared.MatchResult) Calculator {
	avgTeamOneElo := calculateAvgElo(teamOnePlayers)
	avgTeamTwoElo := calculateAvgElo(teamTwoPlayers)

	// transformed ratings
	R1 := math.Pow(10, avgTeamOneElo/400)
	R2 := math.Pow(10, avgTeamTwoElo/400)

	// expected score
	E1 := R1 / (R1 + R2)
	E2 := R2 / (R1 + R2)

	// score modifiers
	S1, S2 := shared.ScoreModifier(matchResult)

	return Calculator{
		teamOneEloLambda: int16(k * (S1 - E1)),
		teamTwoEloLambda: int16(k * (S2 - E2)),
	}
}

func calculateAvgElo(players []dto.Player) float64 {
	var sum int16
	var numberOfPlayers int = len(players)

	for i := 0; i < int(numberOfPlayers); i++ {
		sum += players[i].Elo
	}
	return float64(sum / int16(numberOfPlayers))
}

func (c *Calculator) TeamOneEloLambda() int16 {
	return c.teamOneEloLambda
}

func (c *Calculator) TeamTwoEloLambda() int16 {
	return c.teamTwoEloLambda
}
