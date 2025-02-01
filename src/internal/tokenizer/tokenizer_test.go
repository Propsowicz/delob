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

	token := result.(AddPlayersToken)

	if len(token.Keys) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(token.Keys))
	}
	if token.Keys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.Keys[0])
	}
}

func Test_IfSyntaxForAddingMutlipleUsersIsCorrect(t *testing.T) {
	callStringMock := "ADD PLAYER 'Tomek', 'Romek';"
	_, err := Tokenize(callStringMock)

	if err == nil {
		t.Errorf("Expression should return error.")
	}
}

func Test_IfCanAddTwoPlayers(t *testing.T) {
	callStringMock := "ADD PLAYERS ('Tomek', 'Romek');"
	result, err := Tokenize(callStringMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(AddPlayersToken)

	if len(token.Keys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.Keys))
	}
	if token.Keys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.Keys[0])
	}
	if token.Keys[1] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.Keys[1])
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

func Test_IfCanAddMatchEventForTwoPlayers(t *testing.T) {
	expressionMock := "SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(UpdatePlayersToken)

	if len(token.WinKeys) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(token.WinKeys))
	}
	if len(token.LoseKeys) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(token.LoseKeys))
	}
	if token.WinKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.WinKeys[0])
	}
	if token.LoseKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.LoseKeys[0])
	}
}

func Test_IfCanAddMatchEventForMultiplePlayers(t *testing.T) {
	expressionMock := "SET WIN FOR ('Tomek', 'Joe') AND LOSE FOR ('Romek','John');"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(UpdatePlayersToken)

	if len(token.WinKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.WinKeys))
	}
	if len(token.LoseKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.LoseKeys))
	}
	if token.WinKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.WinKeys[0])
	}
	if token.WinKeys[1] != "Joe" {
		t.Errorf("Expected: %s, got: %s.", "Joe", token.WinKeys[1])
	}
	if token.LoseKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.LoseKeys[0])
	}
	if token.LoseKeys[1] != "John" {
		t.Errorf("Expected: %s, got: %s.", "John", token.LoseKeys[1])
	}
}
