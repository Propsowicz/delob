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
	if result[0].header.entityId != entityIdMock {
		t.Errorf("Wrong EntityId.")
	}
	if result[0].header.isLocked != false {
		t.Errorf("Page should be unlocked.")
	}
	if result[0].header.lastUsedIndex != 0 {
		t.Errorf("Last used index should be 0.")
	}
	if len(result[0].body) != int(utils.PAGE_SIZE) {
		t.Errorf("Body has wrong number of records.")
	}
	if result[0].body[0].method != add {
		t.Errorf("First record should add value.")
	}
	if result[0].body[0].value != float32(utils.INITIAL_ELO) {
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
