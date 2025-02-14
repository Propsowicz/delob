package tokenizer

import (
	"delob/internal/shared"
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

	token := result.(AddPlayersCommand)

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

	token := result.(AddPlayersCommand)

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

	token := result.(AddMatchCommand)

	if token.MatchResult != shared.TeamOneWins {
		t.Errorf("Expected: %d, got: %d.", shared.TeamOneWins, token.MatchResult)
	}
	if len(token.TeamOneKeys) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(token.TeamOneKeys))
	}
	if len(token.TeamTwoKeys) != 1 {
		t.Errorf("Expected: %d, got: %d.", 1, len(token.TeamTwoKeys))
	}
	if token.TeamOneKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.TeamOneKeys[0])
	}
	if token.TeamTwoKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.TeamTwoKeys[0])
	}
}

func Test_IfCanAddMatchEventForMultiplePlayers(t *testing.T) {
	expressionMock := "SET WIN FOR ('Tomek', 'Joe') AND LOSE FOR ('Romek','John');"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(AddMatchCommand)

	if token.MatchResult != shared.TeamOneWins {
		t.Errorf("Expected: %d, got: %d.", shared.TeamOneWins, token.MatchResult)
	}
	if len(token.TeamOneKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamOneKeys))
	}
	if len(token.TeamTwoKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamTwoKeys))
	}
	if token.TeamOneKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.TeamOneKeys[0])
	}
	if token.TeamOneKeys[1] != "Joe" {
		t.Errorf("Expected: %s, got: %s.", "Joe", token.TeamOneKeys[1])
	}
	if token.TeamTwoKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.TeamTwoKeys[0])
	}
	if token.TeamTwoKeys[1] != "John" {
		t.Errorf("Expected: %s, got: %s.", "John", token.TeamTwoKeys[1])
	}
}

func Test_IfCanAddMatchEventForMultiplePlayersWithReverterOrder(t *testing.T) {
	expressionMock := "SET LOSE FOR ('Tomek', 'Joe') AND WIN FOR ('Romek','John');"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(AddMatchCommand)

	if token.MatchResult != shared.TeamTwoWins {
		t.Errorf("Expected: %d, got: %d.", shared.TeamTwoWins, token.MatchResult)
	}
	if len(token.TeamOneKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamOneKeys))
	}
	if len(token.TeamTwoKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamTwoKeys))
	}
	if token.TeamOneKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.TeamOneKeys[0])
	}
	if token.TeamOneKeys[1] != "Joe" {
		t.Errorf("Expected: %s, got: %s.", "Joe", token.TeamOneKeys[1])
	}
	if token.TeamTwoKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.TeamTwoKeys[0])
	}
	if token.TeamTwoKeys[1] != "John" {
		t.Errorf("Expected: %s, got: %s.", "John", token.TeamTwoKeys[1])
	}
}

func Test_IfCanAddMatchEventForMultiplePlayersWithDrawBetweenThem(t *testing.T) {
	expressionMock := "SET DRAW BETWEEN ('Tomek', 'Joe') AND ('Romek','John');"
	result, err := Tokenize(expressionMock)

	if err != nil {
		t.Errorf("Expression should not return error.")
	}

	token := result.(AddMatchCommand)

	if token.MatchResult != shared.Draw {
		t.Errorf("Expected: %d, got: %d.", shared.Draw, token.MatchResult)
	}
	if len(token.TeamOneKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamOneKeys))
	}
	if len(token.TeamTwoKeys) != 2 {
		t.Errorf("Expected: %d, got: %d.", 2, len(token.TeamTwoKeys))
	}
	if token.TeamOneKeys[0] != "Tomek" {
		t.Errorf("Expected: %s, got: %s.", "Tomek", token.TeamOneKeys[0])
	}
	if token.TeamOneKeys[1] != "Joe" {
		t.Errorf("Expected: %s, got: %s.", "Joe", token.TeamOneKeys[1])
	}
	if token.TeamTwoKeys[0] != "Romek" {
		t.Errorf("Expected: %s, got: %s.", "Romek", token.TeamTwoKeys[0])
	}
	if token.TeamTwoKeys[1] != "John" {
		t.Errorf("Expected: %s, got: %s.", "John", token.TeamTwoKeys[1])
	}
}

func Test_IfCanGetSelecetExpressionComponents(t *testing.T) {
	expressionMock := "SELECT Players JOIN Matches WHERE Key = 'zxc' AND Elo > 2500 ORDER BY Elo ASC"

	selectToken, joinToken, whereToken, orderToken := getSelectExpressionTokens(expressionMock)

	if selectToken != "Players" {
		t.Errorf("wrong token: expected %s, got %s", "Players", selectToken)
	}
	if joinToken != "Matches" {
		t.Errorf("wrong token: expected %s, got %s", "Matches", joinToken)
	}
	if whereToken != "Key = 'zxc' AND Elo > 2500" {
		t.Errorf("wrong token: expected %s, got %s", "Key = 'zxc' AND Elo > 2500", whereToken)
	}
	if orderToken != "Elo ASC" {
		t.Errorf("wrong token: expected %s, got %s", "Elo ASC", orderToken)
	}
}

func Test_IfCanGetSelecetExpressionComponentsWithOnlyTwoComponents(t *testing.T) {
	expressionMock := "SELECT Players ORDER BY Elo ASC"

	selectToken, joinToken, whereToken, orderToken := getSelectExpressionTokens(expressionMock)

	if selectToken != "Players" {
		t.Errorf("wrong token: expected %s, got %s", "Players", selectToken)
	}
	if joinToken != "" {
		t.Errorf("wrong token: expected %s, got %s", "", joinToken)
	}
	if whereToken != "" {
		t.Errorf("wrong token: expected %s, got %s", "", whereToken)
	}
	if orderToken != "Elo ASC" {
		t.Errorf("wrong token: expected %s, got %s", "Elo ASC", orderToken)
	}
}
