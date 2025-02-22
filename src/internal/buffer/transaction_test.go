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

func Test_IfCannotReadDirtyPageRecords(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	bufferManager, _ := NewBufferManager()
	entityIdMock := "456789"

	transaction_1 := NewTransaction()
	transaction_1.Start()
	err_1 := bufferManager.AddPlayer(entityIdMock, initElo, nil, &transaction_1)
	transaction_1.EvaluateTransactionSuccess(err_1)
	transaction_1.Finish()

	result_1, _ := bufferManager.GetPages(entityIdMock)

	if len(result_1) != 1 {
		t.Errorf("Page should be created, since it is not dirty")
	}
	if result_1[len(result_1)-1].Body[0].transactionStatus != success {
		t.Errorf("Initial record should be visible")
	}

	transaction_2 := NewTransaction()
	transaction_2.Start()
	err_2 := bufferManager.UpdatePlayer(entityIdMock, 25, nil, &transaction_2)

	result_2, _ := bufferManager.GetPages(entityIdMock)

	if result_2[len(result_2)-1].Body[1].transactionStatus != inProgress {
		t.Errorf("Only initial record should be visible")
	}

	transaction_2.EvaluateTransactionSuccess(err_2)
	transaction_2.Finish()

	result_3, _ := bufferManager.GetPages(entityIdMock)

	if result_3[len(result_3)-1].Body[1].transactionStatus != success {
		t.Errorf("New record should also be visible")
	}

}
