package tokenizer

import (
	"testing"
)

func Test_IfRandomStringReturnsError(t *testing.T) {
	callStringMock := "random string"
	_, err := Tokenize(callStringMock)

	if err == nil {
		t.Errorf("Expression should return error.")
	}
}

func Test_IfCanAddPlayer(t *testing.T) {
	callStringMock := "ADD PLAYER 'Tomek';"
	result, err := Tokenize(callStringMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}
	if len(result) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(result))
	}
	if result[0].ProcessMethod != AddPlayer {
		t.Errorf("Expected: %d, got: %d.", AddPlayer, result[0].ProcessMethod)
	}
	if len(result[0].Arguments) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(result[0].Arguments))
	}
	if result[0].Arguments[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", result[0].Arguments[0])
	}
}

func Test_IfCanAddTwoPlayers(t *testing.T) {
	callStringMock := "ADD PLAYER 'Tomek', 'Romek';"
	result, err := Tokenize(callStringMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}
	if len(result) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(result))
	}
	if result[0].ProcessMethod != AddPlayer {
		t.Errorf("Expected: %d, got: %d.", AddPlayer, result[0].ProcessMethod)
	}
	if len(result[0].Arguments) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(result[0].Arguments))
	}
	if result[0].Arguments[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", result[0].Arguments[0])
	}
	if result[0].Arguments[1] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", result[0].Arguments[0])
	}
}

func Test_IfWrongFormatDoesReturnError(t *testing.T) {
	callStringMock := "ADD PLAYER 'Tomek', 'Romek"
	_, err := Tokenize(callStringMock)

	if err == nil {
		t.Errorf("Expression should return error.")
	}
}

func Test_IfCanNotUpdateIndividualPlayerData(t *testing.T) {
	expressionMock := "SET WIN FOR 'Tomek';"
	_, err := Tokenize(expressionMock)

	if err == nil {
		t.Errorf("Expression should return error.")
	}
}

func Test_IfCanUpdatePlayersData(t *testing.T) {
	expressionMock := "SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}
	if len(result) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(result))
	}
	if result[0].ProcessMethod != UpdatePlayers {
		t.Errorf("Expected: %d, got: %d.", UpdatePlayers, result[0].ProcessMethod)
	}
	if len(result[0].Arguments) != 4 {
		t.Errorf("Expected: %d, got: %d.", 1, len(result[0].Arguments))
	}
	if result[0].Arguments[0] != "WIN" {
		t.Errorf("Expected: %s, got: %s.", "WIN", result[0].Arguments[0])
	}
	if result[0].Arguments[1] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", result[0].Arguments[1])
	}
	if result[0].Arguments[2] != "LOSE" {
		t.Errorf("Expected: %s, got: %s.", "LOSE", result[0].Arguments[2])
	}
	if result[0].Arguments[3] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", result[0].Arguments[3])
	}
}
