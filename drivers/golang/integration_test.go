package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func assertCorrectKeyOrder(t *testing.T, player Player, expectedKey string) {
	if player.Key != expectedKey {
		t.Errorf("wrong order: expected %s, got %s", expectedKey, player.Key)
	}
}

type TestCase struct {
	Context []TestRecord
}

type TestRecord struct {
	Type      string
	Arguments []string
}

func loadTestCase(fileName string) TestCase {
	result := TestCase{}
	f, err := os.ReadFile(fileName)
	if err != nil {
		return result
	}

	json.Unmarshal(f, &result)
	return result
}

func executeTestCase(context *DelobContext, record TestRecord) error {
	if record.Type == "addPlayer" {
		return context.AddPlayer(record.Arguments[0])
	}
	if record.Type == "addPlayers" {
		return context.AddPlayers(record.Arguments)
	}
	if record.Type == "setDecisiveMatch" {
		return context.SetDecisiveMatch(record.Arguments[0], record.Arguments[1])
	}
	if record.Type == "setDecisiveTeamMatch" {
		len := len(record.Arguments)
		return context.SetDecisiveTeamMatch(record.Arguments[:len/2], record.Arguments[len/2:])
	}
	if record.Type == "setDrawMatch" {
		return context.SetDrawMatch(record.Arguments[0], record.Arguments[1])
	}
	if record.Type == "setDrawTeamMatch" {
		len := len(record.Arguments)
		return context.SetDrawTeamMatch(record.Arguments[:len/2], record.Arguments[len/2:])
	}

	return nil
}

func Test_TestCase_1(t *testing.T) {
	context, err := NewContext("")
	if err != nil {
		t.Errorf("Should be able to create delob context")
	}

	testCases := loadTestCase("testcase_1.json")

	for _, testCase := range testCases.Context {
		execErr := executeTestCase(&context, testCase)

		if execErr != nil {
			t.Errorf("Error during execution: %s", execErr.Error())
		}
	}

	result, errResult := context.GetPlayersOrderBy(Elo, Descending)
	if errResult != nil {
		t.Errorf("Should be able to create delob context")
	}
	if len(result) != 10 {
		t.Errorf("Created wrong number of players: expected 10. got %d.", len(result))
	}
	assertCorrectKeyOrder(t, result[0], "qJWWl")
	assertCorrectKeyOrder(t, result[1], "6LTnq")
	assertCorrectKeyOrder(t, result[2], "0cLEk")
	assertCorrectKeyOrder(t, result[3], "Nr7tf")
	assertCorrectKeyOrder(t, result[4], "ts58a")
	assertCorrectKeyOrder(t, result[5], "XnF1a")
	assertCorrectKeyOrder(t, result[6], "FAJur")
	assertCorrectKeyOrder(t, result[7], "FoEqy")
	assertCorrectKeyOrder(t, result[8], "caka4")
	assertCorrectKeyOrder(t, result[9], "nL7QC")
}

func Test_TestCase_1_SELECT_PerformanceTest(t *testing.T) {
	const expectedExecutionTimeInMiliseconds int64 = 10
	context, err := NewContext("")
	if err != nil {
		t.Errorf("Should be able to create delob context")
	}

	start := time.Now()

	_, errResult := context.GetPlayersOrderBy(Elo, Descending)
	elapsed := time.Since(start).Milliseconds()

	if errResult != nil {
		t.Errorf("Should be able to create delob context")
	}
	if elapsed >= expectedExecutionTimeInMiliseconds {
		t.Errorf("GetPlayersOrderBy performance test not passed: expected less than %dms. got %dms.",
			expectedExecutionTimeInMiliseconds, elapsed)
	}
}
