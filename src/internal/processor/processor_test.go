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

func Test_IfCannotUpdateWhenIdDoesnNotExists(t *testing.T) {
	bufferManager := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tomek';"
	p.Execute(expressionMock)

	_, err := p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")

	if err == nil {
		t.Errorf("Should throw error.")
	}

	bufferManager = buffer.NewBufferManager()
	p = Processor{bufferManager: &bufferManager}
	expressionMock = "ADD PLAYER 'Romek';"
	p.Execute(expressionMock)

	_, err = p.Execute("SET WIN FOR 'Tomek' AND LOSE FOR 'Romek';")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_CanCorrectlyCalculateEloForOnePage(t *testing.T) {
	var startValue int16 = 1300
	var addValue int16 = 25
	var substractValue int16 = 35
	player := Player{}

	player.records = append(player.records, addValue)

	player.calculateElo()
	if player.Elo != startValue+addValue {
		t.Errorf("Wrong elo calculation result")
	}

	player.records = append(player.records, substractValue)
	player.records = append(player.records, addValue)

	player.calculateElo()
	if player.Elo != startValue+addValue-substractValue+addValue {
		t.Errorf("Wrong elo calculation result")
	}
}
