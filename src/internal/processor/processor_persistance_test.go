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

	processorFirstRun.Execute("traceId", "ADD PLAYER 'Tom';")
	processorFirstRun.Execute("traceId", "ADD PLAYERS ('Joe', 'Bob', 'Jim');")
	processorFirstRun.Execute("traceId", "SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	processorFirstRun.Execute("traceId", "SET WIN FOR ('Tom') AND LOSE FOR ('Joe', 'Jim');")
	processorFirstRun.Execute("traceId", "SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');")

	resultFirstRun, _ := processorFirstRun.Execute("traceId", "SELECT Players ORDER BY Elo DESC;")

	bufferManagerSecondRun, _ := buffer.NewBufferManager()
	processorSecondRun := Processor{bufferManager: &bufferManagerSecondRun}
	processorSecondRun.Initialize()

	resultSecondRun, _ := processorSecondRun.Execute("traceId", "SELECT Players ORDER BY Elo DESC;")

	if resultFirstRun != resultSecondRun {
		t.Errorf("Data is not persistent before and after initialization from a log data.")
	}
}

func TestIfDataIsPersistentBetweenDatabaseRunsWithTransactionalData(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManagerFirstRun, _ := buffer.NewBufferManager()
	processorFirstRun := Processor{bufferManager: &bufferManagerFirstRun}

	processorFirstRun.Execute("traceId", "ADD PLAYER 'Tom';")

	// this one should fail and I should be able to add 'Joe', 'Bob', 'Jim' in the next run
	processorFirstRun.Execute("traceId", "ADD PLAYERS ('Tom', 'Joe', 'Bob', 'Jim');")

	processorFirstRun.Execute("traceId", "ADD PLAYERS ('Joe', 'Bob', 'Jim');")
	processorFirstRun.Execute("traceId", "SET WIN FOR 'Tom' AND LOSE FOR 'Joe';")
	processorFirstRun.Execute("traceId", "SET WIN FOR ('Tom') AND LOSE FOR ('Joe', 'Jim');")
	processorFirstRun.Execute("traceId", "SET WIN FOR ('Tom', 'Bob') AND LOSE FOR ('Joe', 'Jim');")

	resultFirstRun, _ := processorFirstRun.Execute("traceId", "SELECT Players ORDER BY Elo DESC;")

	if resultFirstRun != "[{\"Key\":\"Tom\",\"Elo\":1331},{\"Key\":\"Bob\",\"Elo\":1315},{\"Key\":\"Jim\",\"Elo\":1285},{\"Key\":\"Joe\",\"Elo\":1269}]" {
		t.Errorf("Data is not correct - %s.", resultFirstRun)
	}

	bufferManagerSecondRun, _ := buffer.NewBufferManager()
	processorSecondRun := Processor{bufferManager: &bufferManagerSecondRun}
	processorSecondRun.Initialize()

	resultSecondRun, _ := processorSecondRun.Execute("traceId", "SELECT Players ORDER BY Elo DESC;")

	if resultFirstRun != resultSecondRun {
		t.Errorf("Data is not persistent before and after initialization from a log data.")
	}
}
