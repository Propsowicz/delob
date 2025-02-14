package processor

import (
	buffer "delob/internal/buffer"
	"testing"
)

func TestIfDataIsPersistentBetweenDatabaseRuns(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManagerFirstRun, _ := buffer.NewBufferManager()
	processorFirstRun := Processor{bufferManager: &bufferManagerFirstRun}

	processorFirstRun.Execute("ADD PLAYER 'Tom';")
	processorFirstRun.Execute("ADD PLAYERS ('Joe', 'Bob', 'Jim');")
	processorFirstRun.Execute("SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	processorFirstRun.Execute("SET WIN FOR ('Tom') AND LOSE FOR ('Joe', 'Jim');")
	processorFirstRun.Execute("SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');")

	resultFirstRun, _ := processorFirstRun.Execute("SELECT Players ORDER BY Elo DESC;")

	bufferManagerSecondRun, _ := buffer.NewBufferManager()
	processorSecondRun := Processor{bufferManager: &bufferManagerSecondRun}
	processorSecondRun.Initialize()

	resultSecondRun, _ := processorSecondRun.Execute("SELECT Players ORDER BY Elo DESC;")

	if resultFirstRun != resultSecondRun {
		t.Errorf("Data is not persistent before and after initialization from a log data.")
	}
}
