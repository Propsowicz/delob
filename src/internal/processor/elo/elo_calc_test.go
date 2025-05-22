package elo

import (
	dto "delob/internal/processor/model"
	"delob/internal/shared"
	"testing"
)

const k_mock int16 = 25

func are_not_equal(lambda int16, expected int16) bool {
	return lambda != expected
}

// NOTE: expected values are hardcoded based on manual calculations

func Test_IfCanCalculateCorrectEloForDrawBetweenPlayersWithTheSameElo(t *testing.T) {
	teamOne := []dto.Player{{Elo: 1000}}
	teamTwo := []dto.Player{{Elo: 1000}}
	result := shared.Draw

	calc := NewCalculator(teamOne, teamTwo, result)

	var expectedTeamOneLambda int16 = 0
	var expectedTeamTwoLambda int16 = 0

	if are_not_equal(calc.TeamOneEloLambda(), int16(expectedTeamOneLambda)) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamOneLambda, calc.TeamOneEloLambda())
	}
	if are_not_equal(calc.TeamTwoEloLambda(), int16(expectedTeamTwoLambda)) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamTwoLambda, calc.TeamTwoEloLambda())
	}
}

func Test_IfCanCalculateCorrectEloForDecisiveResultAndTheSameStartingElo(t *testing.T) {
	teamOne := []dto.Player{{Elo: 1000}}
	teamTwo := []dto.Player{{Elo: 1000}}
	result := shared.TeamOneWins

	calc := NewCalculator(teamOne, teamTwo, result)

	var expectedTeamOneLambda int16 = 16
	var expectedTeamTwoLambda int16 = -expectedTeamOneLambda

	if are_not_equal(calc.TeamOneEloLambda(), expectedTeamOneLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamOneLambda, calc.TeamOneEloLambda())
	}
	if are_not_equal(calc.TeamTwoEloLambda(), expectedTeamTwoLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamTwoLambda, calc.TeamTwoEloLambda())
	}
}

func Test_IfCanCalculateCorrectEloForDecisiveResultAndDifferentStartingElo_WinnerHasMoreElo(t *testing.T) {
	teamOne := []dto.Player{{Elo: 1500}}
	teamTwo := []dto.Player{{Elo: 1000}}
	result := shared.TeamOneWins

	calc := NewCalculator(teamOne, teamTwo, result)

	var expectedTeamOneLambda int16 = 2
	var expectedTeamTwoLambda int16 = -expectedTeamOneLambda

	if are_not_equal(calc.TeamOneEloLambda(), expectedTeamOneLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamOneLambda, calc.TeamOneEloLambda())
	}
	if are_not_equal(calc.TeamTwoEloLambda(), expectedTeamTwoLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamTwoLambda, calc.TeamTwoEloLambda())
	}
}

func Test_IfCanCalculateCorrectEloForDecisiveResultAndDifferentStartingElo_WinnerHasLessElo(t *testing.T) {
	teamOne := []dto.Player{{Elo: 1500}}
	teamTwo := []dto.Player{{Elo: 1000}}
	result := shared.TeamTwoWins

	calc := NewCalculator(teamOne, teamTwo, result)

	var expectedTeamOneLambda int16 = -30
	var expectedTeamTwoLambda int16 = -expectedTeamOneLambda

	if are_not_equal(calc.TeamOneEloLambda(), expectedTeamOneLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamOneLambda, calc.TeamOneEloLambda())
	}
	if are_not_equal(calc.TeamTwoEloLambda(), expectedTeamTwoLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedTeamTwoLambda, calc.TeamTwoEloLambda())
	}
}

func Test_IfCanCalculateCorrectEloForDecisiveResultAndDifferentStartingElo_WinnerHasLessElo2(t *testing.T) {
	teamOne := []dto.Player{{Elo: 1379}, {Elo: 1362}, {Elo: 1328}, {Elo: 1300}, {Elo: 1285}, {Elo: 1253}}
	teamTwo := []dto.Player{{Elo: 1362}, {Elo: 1362}, {Elo: 1314}, {Elo: 1313}, {Elo: 1270}, {Elo: 1268}}
	result := shared.TeamOneWins

	calc := NewCalculator(teamOne, teamTwo, result)
	var expectedLambda int16 = 16

	if are_not_equal(calc.TeamOneEloLambda(), expectedLambda) {
		t.Errorf("Wrong elo calculation; expected: %d, got: %d.", expectedLambda, calc.TeamOneEloLambda())
	}
}
