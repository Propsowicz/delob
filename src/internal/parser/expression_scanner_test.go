package parser

import (
	"testing"
)

func newExpressionScannerMock() ExpressionScanner {
	return newExpressionScanner("traceIdMock")
}

func Test_IfCannotGetLogicalTokensBeforeDetermineExpressionType(t *testing.T) {
	scanner := newExpressionScannerMock()

	_, _, err := scanner.getLogicalTokens()

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCannotUseExpressionWithoutSemicolonAtTheEnd(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("ADD PLAYER 'Tom'")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanScan_AddPlayer_ExpressionAndGetCorrectTypeAndTokens(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("ADD PLAYER 'Tom';")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddPlayersCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 1 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != AddPlayer || result[0].Value[0] != "Tom" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddPlayers_ExpressionAndGetCorrectTypeAndTokens(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("ADD PLAYERS ('Tom', 'Joe','Jim');")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddPlayersCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 1 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != AddPlayers || result[0].Value[0] != "Tom" || result[0].Value[1] != "Joe" || result[0].Value[2] != "Jim" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForSinglePlayers(t *testing.T) {
	scanner := newExpressionScannerMock()
	err := scanner.scanRawExpression("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddMatchCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SetWin || result[0].Value[0] != "Tom" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != SetLose || result[1].Value[0] != "Joe" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForSinglePlayers_WithReplacedArguments(t *testing.T) {
	firstScanner := newExpressionScannerMock()
	firstScanner.scanRawExpression("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	_, firstResult, _ := firstScanner.getLogicalTokens()

	secondScanner := newExpressionScannerMock()
	secondScanner.scanRawExpression("SET LOSE FOR 'Joe' AND WIN FOR 'Tom';")
	_, secondResult, _ := secondScanner.getLogicalTokens()

	if len(firstResult) != len(secondResult) {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[0].Value[0] != secondResult[0].Value[0] {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[1].Value[0] != secondResult[1].Value[0] {
		t.Errorf("Tokens that should be identical")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForTeams(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddMatchCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SetWin {
		t.Errorf("Invalid logical token")
	}
	if result[0].Value[0] != "Tom" || result[0].Value[1] != "Bob" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != SetLose {
		t.Errorf("Invalid logical token")
	}
	if result[1].Value[0] != "Joe" || result[1].Value[1] != "Jim" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForTeams_WithReplacedArguments(t *testing.T) {
	firstScanner := newExpressionScannerMock()
	firstScanner.scanRawExpression("SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');")
	_, firstResult, _ := firstScanner.getLogicalTokens()

	secondScanner := newExpressionScannerMock()
	secondScanner.scanRawExpression("SET LOSE FOR ('Joe', 'Jim') AND WIN FOR ('Tom', 'Bob');")
	_, secondResult, _ := secondScanner.getLogicalTokens()

	if len(firstResult) != len(secondResult) {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[0].Value[0] != secondResult[0].Value[0] {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[0].Value[1] != secondResult[0].Value[1] {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[1].Value[0] != secondResult[1].Value[0] {
		t.Errorf("Tokens that should be identical")
	}
	if firstResult[1].Value[1] != secondResult[1].Value[1] {
		t.Errorf("Tokens that should be identical")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForSinglePlayers_Draw(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SET DRAW BETWEEN 'Joe' AND 'Tom';")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddMatchCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SetDraw {
		t.Errorf("Invalid logical token")
	}
	if result[0].Value[0] != "Joe" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != SetDraw {
		t.Errorf("Invalid logical token")
	}
	if result[1].Value[0] != "Tom" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddMatch_ExpressionAndGetCorrectTypeAndTokens_ForTeams_Draw(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SET DRAW BETWEEN ('Joe','Jim') AND ('Tom', 'Bob');")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != AddMatchCommandType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SetDraw {
		t.Errorf("Invalid logical token")
	}
	if result[0].Value[0] != "Joe" || result[0].Value[1] != "Jim" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != SetDraw {
		t.Errorf("Invalid logical token")
	}
	if result[1].Value[0] != "Tom" || result[1].Value[1] != "Bob" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_SelectPlayer_ExpressionAndGetCorrectTypeAndTokens(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SELECT Players;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != SelectQueryType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 1 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SelectPlayers {
		t.Errorf("Invalid logical token")
	}
	if result[0].Value[0] != "*" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_SelectPlayer_ExpressionAndGetCorrectTypeAndTokens_OrderByEloAsc(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SELECT Players ORDER BY Elo ASC;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != SelectQueryType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SelectPlayers || result[0].Value[0] != "*" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != OrderByAsc || result[1].Value[0] != "Elo" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_SelectPlayer_ExpressionAndGetCorrectTypeAndTokens_OrderByKeyDesc(t *testing.T) {
	scanner := newExpressionScannerMock()

	err := scanner.scanRawExpression("SELECT Players ORDER BY Key DESC;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	expressionType, result, errTokens := scanner.getLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if expressionType != SelectQueryType {
		t.Errorf("Wrong expression type.")
	}
	if len(result) != 2 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != SelectPlayers || result[0].Value[0] != "*" {
		t.Errorf("Invalid logical token")
	}
	if result[1].Token != OrderByDesc || result[1].Value[0] != "Key" {
		t.Errorf("Invalid logical token")
	}
}
