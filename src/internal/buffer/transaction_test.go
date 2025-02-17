package buffer

import "testing"

func Test_IfCannotReadDirtyPageDictionaryEntities(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	transaction := NewTransaction()

	bufferManager, _ := NewBufferManager()
	entityIdMock := "123987"

	bufferManager.AddPlayer(entityIdMock, initElo, nil, &transaction)

	result, _ := bufferManager.GetPages(entityIdMock)

	if len(result) != 0 {
		t.Errorf("Page that is dirty should not be found.")
	}

	transaction.Finish(true)
	result1, _ := bufferManager.GetPages(entityIdMock)

	if len(result1) != 1 {
		t.Errorf("Page that is not dirty should be found.")
	}
}
