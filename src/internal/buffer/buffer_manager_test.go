package buffer

import (
	"delob/internal/utils"
	"testing"
)

func Test_IfCanAddPlayerToBuffer(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"

	err := bufferManager.AddPlayer(entityIdMock)
	result, getPageErr := bufferManager.GetPage(entityIdMock)

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
	if result[0].Body[0].Method != Add {
		t.Errorf("First record should add value.")
	}
	if result[0].Body[0].Value != utils.INITIAL_ELO {
		t.Errorf("First record should contain proper initial elo value.")
	}
}

func Test_IfGetErrorWhenTryToAddExistingPlayerToBuffer(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"

	err1 := bufferManager.AddPlayer(entityIdMock)
	err2 := bufferManager.AddPlayer(entityIdMock)

	if err1 != nil {
		t.Errorf("Should create player without any errors.")
	}
	if err2 == nil {
		t.Errorf("Should NOT create player.")
	}
}

func Test_IfGetErrorWhenTryToUpdateNotExistingPlayerToBuffer(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	recordMock := Record{
		Method: Add,
		Value:  25,
	}

	err := bufferManager.UpdatePlayer(entityIdMock, recordMock)

	if err == nil {
		t.Errorf("Should not update player.")
	}
}

func Test_IfCanUpdatePlayerWhenCanAddToExistingPage(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	recordMock := Record{
		Method: Add,
		Value:  25,
	}
	bufferManager.AddPlayer(entityIdMock)

	err := bufferManager.UpdatePlayer(entityIdMock, recordMock)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPage(entityIdMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Method != 0 {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[0].Value != utils.INITIAL_ELO {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[1].Method != recordMock.Method {
		t.Errorf("New record should be added.")
	}
	if result[0].Body[1].Value != recordMock.Value {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenThereIsOnlyOneSlotInPage(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	recordMock := Record{Method: Add, Value: 25}
	bufferManager.AddPlayer(entityIdMock)

	for i := 0; i < int(utils.PAGE_SIZE)-2; i++ {
		bufferManager.UpdatePlayer(entityIdMock, recordMock)
	}

	err := bufferManager.UpdatePlayer(entityIdMock, recordMock)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPage(entityIdMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Method != 0 {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[0].Value != utils.INITIAL_ELO {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[utils.PAGE_SIZE-1].Method != recordMock.Method {
		t.Errorf("New record should be added.")
	}
	if result[0].Body[utils.PAGE_SIZE-1].Value != recordMock.Value {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenPageIsFull(t *testing.T) {
	bufferManager := NewBufferManager()
	entityIdMock := "123987"
	recordMock := Record{Method: Add, Value: 25}
	bufferManager.AddPlayer(entityIdMock)

	for i := 0; i < int(utils.PAGE_SIZE)-1; i++ {
		bufferManager.UpdatePlayer(entityIdMock, recordMock)
	}

	err := bufferManager.UpdatePlayer(entityIdMock, recordMock)

	if err != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPage(entityIdMock)

	if len(result) != 2 {
		t.Errorf("There should be two pages.")
	}
	if result[1].Body[0].Method != recordMock.Method {
		t.Errorf("New record should be added.")
	}
	if result[1].Body[0].Value != recordMock.Value {
		t.Errorf("New record should be added.")
	}
}
