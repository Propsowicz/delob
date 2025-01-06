package processor

import (
	buffer "delob/internal/buffer"
	"testing"
)

func Test_IfCanNotTokenizeExpressionWithoutSemicolonEnds(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	expressionMock := "ADD PLAYER 'Tomek'"

	_, err := Execute(expressionMock, &bufferManager)

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanAddOnePlayer(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	expressionMock := "ADD PLAYER 'Tomek';"

	result, err := Execute(expressionMock, &bufferManager)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "1 row(s) affected" {
		t.Errorf("Adding should affect 1 row.")
	}
}

func Test_IfCanAddTwoPlayers(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	expressionMock := "ADD PLAYER 'Tomek', 'Romek';"

	result, err := Execute(expressionMock, &bufferManager)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "2 row(s) affected" {
		t.Errorf("Adding should affect 2 rows.")
	}
}

func Test_IfCannotAddTheSamePlayerTwicePlayer(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	firstExpressionMock := "ADD PLAYER 'Tomek';"
	secondExpressionMock := "ADD PLAYER 'Tomek', 'Romek';"

	result1, err1 := Execute(firstExpressionMock, &bufferManager)
	result2, err2 := Execute(secondExpressionMock, &bufferManager)

	if err1 != nil {
		t.Errorf("Should not throw error.")
	}
	if result1 != "1 row(s) affected" {
		t.Errorf("Adding should affect 1 row.")
	}

	if err2 == nil {
		t.Errorf("Should throw error.")
	}
	if result2 != "1 row(s) affected" {
		t.Errorf("Adding should affect 1 row.")
	}
}
