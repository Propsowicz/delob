package buffer

import (
	"testing"
)

func Test_IfCannotReadDirtyPageDictionaryEntities(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	transaction := NewTransaction()
	transaction.Start()
	bufferManager, _ := NewBufferManager()
	entityIdMock := "456789"

	err := bufferManager.AddPlayer(entityIdMock, initElo, nil, &transaction)
	result, _ := bufferManager.GetPages(entityIdMock)

	if len(result) != 0 {
		t.Errorf("Page that is dirty should not be found.")
	}

	transaction.EvaluateTransactionSuccess(err)
	transaction.Finish()

	result1, _ := bufferManager.GetPages(entityIdMock)

	if len(result1) != 1 {
		t.Errorf("Page that is not dirty should be found.")
	}
}
