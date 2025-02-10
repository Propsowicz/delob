package buffer

import (
	"delob/internal/utils"
	"os"
	"testing"
)

func setupSuite(_ *testing.T) func(t *testing.T) {
	return func(t *testing.T) {
		backupManagerPAth := "backup"
		os.RemoveAll(backupManagerPAth)
	}
}

const upodateEloValue int16 = 25
const initElo int16 = 1500

func Test_IfCanAddPlayerToBuffer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"

	err := bufferManager.AddPlayer(entityIdMock, initElo, nil)
	result, getPageErr := bufferManager.GetPages(entityIdMock)

	if err != nil {
		t.Errorf("Should create without any errors.")
	}
	if getPageErr != nil {
		t.Errorf("Should get without any errors.")
	}
	if len(result) != 1 {
		t.Errorf("Should create one page.")
	}
	if result[0].Header.entityId != entityIdMock {
		t.Errorf("Wrong EntityId.")
	}
	if result[0].Header.isLocked != false {
		t.Errorf("Page should be unlocked.")
	}
	if result[0].Header.lastUsedIndex != 0 {
		t.Errorf("Last used index should be 0.")
	}
	if len(result[0].Body) != int(utils.PAGE_SIZE) {
		t.Errorf("Body has wrong number of records.")
	}
	if result[0].Body[0].Value != initElo {
		t.Errorf("First record should contain proper initial elo value.")
	}
}

func Test_IfGetErrorWhenTryToAddExistingPlayerToBuffer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"

	err1 := bufferManager.AddPlayer(entityIdMock, initElo, nil)
	err2 := bufferManager.AddPlayer(entityIdMock, initElo, nil)

	if err1 != nil {
		t.Errorf("Should create player without any errors.")
	}
	if err2 == nil {
		t.Errorf("Should NOT create player.")
	}
}

func Test_IfGetErrorWhenTryToUpdateNotExistingPlayerToBuffer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"

	err := bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)

	if err == nil {
		t.Errorf("Should not update player.")
	}
}

func Test_IfCanUpdatePlayerWhenCanAddToExistingPage(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	bufferManager.AddPlayer(entityIdMock, initElo, nil)

	err := bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(entityIdMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Value != initElo {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[1].Value != upodateEloValue {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenThereIsOnlyOneSlotInPage(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	bufferManager.AddPlayer(entityIdMock, initElo, nil)

	for i := 0; i < int(utils.PAGE_SIZE)-2; i++ {
		bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)
	}

	err := bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(entityIdMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Value != initElo {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[utils.PAGE_SIZE-1].Value != upodateEloValue {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenPageIsFull(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	bufferManager.AddPlayer(entityIdMock, initElo, nil)

	for i := 0; i < int(utils.PAGE_SIZE)-1; i++ {
		bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)
	}

	err := bufferManager.UpdatePlayer(entityIdMock, upodateEloValue, nil)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(entityIdMock)

	if len(result) != 2 {
		t.Errorf("There should be two pages.")
	}
	if result[1].Body[0].Value != upodateEloValue {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanAppendMatchToBuffer(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager := NewBufferManager()
	teamOneKeys := []string{"1", "2"}
	teamTwoKeys := []string{"3", "4"}
	matchResult := 0

	result := bufferManager.AddMatchEvent(teamOneKeys, teamTwoKeys, int8(matchResult))

	if len(result.TeamOneKeys) != 2 {
		t.Errorf("There should be two players in team one.")
	}
	if len(result.TeamOneKeys) != 2 {
		t.Errorf("There should be two players in team one.")
	}
}
