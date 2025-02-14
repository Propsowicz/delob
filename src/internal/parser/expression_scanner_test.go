package parser

import (
	"testing"
)

func Test_IfCannotGetLogicalTokensBeforeDetermineExpressionType(t *testing.T) {
	scanner := NewExpressionScanner()

	_, err := scanner.GetLogicalTokens()

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanScan_AddPlayer_ExpressionAndGetCorrectTypeAndTokens(t *testing.T) {
	scanner := NewExpressionScanner()

	err := scanner.ScanRawExpression("ADD PLAYER 'Tom';")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	result, errTokens := scanner.GetLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if len(result) != 1 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != AddPlayer || result[0].Value[0] != "Tom" {
		t.Errorf("Invalid logical token")
	}
}

func Test_IfCanScan_AddPlayers_ExpressionAndGetCorrectTypeAndTokens(t *testing.T) {
	scanner := NewExpressionScanner()

	err := scanner.ScanRawExpression("ADD PLAYERS ('Tom', 'Joe','Jim');")

	if err != nil {
		t.Errorf("Should not throw error.")
	}

	result, errTokens := scanner.GetLogicalTokens()

	if errTokens != nil {
		t.Errorf("Should not throw error.")
	}
	if len(result) != 1 {
		t.Errorf("Wrong numbers of tokens.")
	}
	if result[0].Token != AddPlayers || result[0].Value[0] != "Tom" || result[0].Value[1] != "Joe" || result[0].Value[2] != "Jim" {
		t.Errorf("Invalid logical token")
	}
}
