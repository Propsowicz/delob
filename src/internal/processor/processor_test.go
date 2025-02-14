package processor

import (
	buffer "delob/internal/buffer"
	"delob/internal/processor/model"
	"encoding/json"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func assertCorrectKeyOrder(t *testing.T, record model.Player, expectedKey string) {
	if record.Key != expectedKey {
		t.Errorf("wrong order: expected %s, got %s", expectedKey, record.Key)
	}
}

func setupSuite(_ *testing.T) func(t *testing.T) {
	backupManagerPath := "log_data"
	os.RemoveAll(backupManagerPath)
	return func(t *testing.T) {
		backupManagerPath := "log_data"
		os.RemoveAll(backupManagerPath)
	}
}

func Test_IfCanNotTokenizeExpressionWithoutSemicolonEnds(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tom'"

	_, err := p.Execute(expressionMock)

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanAddOnePlayer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tom';"

	result, err := p.Execute(expressionMock)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "1 row(s) affected" {
		t.Errorf("Adding should affect 1 row.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCanAddTwoPlayers(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYERS ('Tom', 'Joe');"

	result, err := p.Execute(expressionMock)

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	if result != "2 row(s) affected" {
		t.Errorf("Adding should affect 2 rows.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCannotAddTheSamePlayerTwicePlayer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	firstExpressionMock := "ADD PLAYER 'Tom';"
	secondExpressionMock := "ADD PLAYERS ('Tom', 'Joe' );"

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
	snaps.MatchSnapshot(t, result1)
	snaps.MatchSnapshot(t, result2)
}

func Test_IfCannotUpdateWhenIdDoesnNotExists_Case1(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Tom';"
	p.Execute(expressionMock)

	_, err := p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCannotUpdateWhenIdDoesnNotExists_Case2(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}
	expressionMock := "ADD PLAYER 'Joe';"
	p.Execute(expressionMock)

	_, err := p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")

	if err == nil {
		t.Errorf("Should throw error.")
	}
}

func Test_IfCanSelectAllWhenThereIsOnePlayer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tom';")
	result, err := p.Execute("SELECT Players;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCanSelectTwoPlayersWithoutUpdatingResults(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tom';")
	p.Execute("ADD PLAYER 'Joe';")
	result, err := p.Execute("SELECT Players;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCanSelectTwoPlayersWithUpdatingResults(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tom';")
	p.Execute("ADD PLAYER 'Joe';")
	p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	p.Execute("SET WIN FOR 'Joe' AND LOSE FOR 'Tom';")

	result, err := p.Execute("SELECT Players;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCanSelectTwoPlayersWithUpdatingResultsWithDrawResult(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYER 'Tom';")
	p.Execute("ADD PLAYER 'Joe';")
	p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	p.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	p.Execute("SET WIN FOR 'Joe' AND LOSE FOR 'Tom';")
	p.Execute("SET WIN FOR 'Joe' AND LOSE FOR 'Tom';")
	p.Execute("SET DRAW BETWEEN 'Joe' AND 'Tom';")

	result, err := p.Execute("SELECT Players;")

	if err != nil {
		t.Errorf("Should not throw error.")
	}
	snaps.MatchSnapshot(t, result)
}

func Test_IfCanSortAscendingByPlayerKey(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYERS ('A', 'C', 'E');")
	p.Execute("ADD PLAYERS ('B', 'D');")

	result, _ := p.Execute("SELECT Players ORDER BY Key ASC;")

	data := []model.Player{}
	json.Unmarshal([]byte(result), &data)
	assertCorrectKeyOrder(t, data[0], "A")
	assertCorrectKeyOrder(t, data[1], "B")
	assertCorrectKeyOrder(t, data[2], "C")
	assertCorrectKeyOrder(t, data[3], "D")
	assertCorrectKeyOrder(t, data[4], "E")
}

func Test_IfCanSortDescendingByPlayerKey(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYERS ('A', 'C', 'E');")
	p.Execute("ADD PLAYERS ('B', 'D');")

	result, _ := p.Execute("SELECT Players ORDER BY Key DESC;")

	data := []model.Player{}
	json.Unmarshal([]byte(result), &data)
	assertCorrectKeyOrder(t, data[4], "A")
	assertCorrectKeyOrder(t, data[3], "B")
	assertCorrectKeyOrder(t, data[2], "C")
	assertCorrectKeyOrder(t, data[1], "D")
	assertCorrectKeyOrder(t, data[0], "E")
}

func Test_IfCanSortDescendingByPlayerElo(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := buffer.NewBufferManager()
	p := Processor{bufferManager: &bufferManager}

	p.Execute("ADD PLAYERS ('A', 'B', 'X');")

	p.Execute("SET WIN FOR 'B' AND LOSE FOR 'X';")
	result, _ := p.Execute("SELECT Players ORDER BY Elo DESC;")
	p.Execute("SET WIN FOR 'B' AND LOSE FOR 'X';")
	p.Execute("SET WIN FOR 'A' AND LOSE FOR 'X';")

	data := []model.Player{}
	json.Unmarshal([]byte(result), &data)
	assertCorrectKeyOrder(t, data[0], "B")
	assertCorrectKeyOrder(t, data[1], "A")
	assertCorrectKeyOrder(t, data[2], "X")
}
