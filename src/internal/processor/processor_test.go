package processor

import (
	buffer "delob/internal/buffer"
	"testing"
)

func Test_IfCanNotTokenizeExpressionWithoutSemicolonEnds(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tomek'"

	_, err := p.Execute(expressionMock)

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanAddOnePlayer(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tomek';"

	result, err := p.Execute(expressionMock)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "1 row(s) affected" {
		t.Errorf("Adding should affect 1 row.")
	}
}

func Test_IfCanAddTwoPlayers(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tomek', 'Romek';"

	result, err := p.Execute(expressionMock)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "2 row(s) affected" {
		t.Errorf("Adding should affect 2 rows.")
	}
}

func Test_IfCannotAddTheSamePlayerTwicePlayer(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	firstExpressionMock := "ADD PLAYER 'Tomek';"
	secondExpressionMock := "ADD PLAYER 'Tomek', 'Romek';"

	result1, err1 := p.Execute(firstExpressionMock)
	result2, err2 := p.Execute(secondExpressionMock)

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

func Test_IfCannotUpdateWhenIdDoesnNotExists_Case1(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tomek';"
	p.Execute(expressionMock)

	_, err := p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCannotUpdateWhenIdDoesnNotExists_Case2(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Romek';"
	p.Execute(expressionMock)

	_, err := p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanSelectAllWhenThereIsOnePlayer(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tomek';")
	result, err := p.Execute("SELECT ALL;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "[{\"Id\":\"Tomek\",\"Elo\":1300}]" {
		t.Errorf("incorrect result - %s", result)
	}
}

func Test_IfCanSelectTwoPlayersWithoutUpdatingResults(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tomek';")
	p.Execute("ADD PLAYER 'Romek';")
	result, err := p.Execute("SELECT ALL;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "[{\"Id\":\"Tomek\",\"Elo\":1300},{\"Id\":\"Romek\",\"Elo\":1300}]" {
		t.Errorf("incorrect result - %s", result)
	}
}

func Test_IfCanSelectTwoPlayersWithUpdatingResults(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tomek';")
	p.Execute("ADD PLAYER 'Romek';")
	p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")
	p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")
	p.Execute("SET WIN FOR 'Romek' AND LOSE FOR 'Tomek';")

	result, err := p.Execute("SELECT ALL;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "[{\"Id\":\"Tomek\",\"Elo\":1325},{\"Id\":\"Romek\",\"Elo\":1275}]" {
		t.Errorf("incorrect result - %s", result)
	}
}
