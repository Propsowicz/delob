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
	keyMock := "456789"

	err := bufferManager.AddPlayer(keyMock, initElo, nil, &transaction)
	result, _ := bufferManager.GetPages(keyMock)

	if len(result) != 0 {
		t.Errorf("Page that is dirty should not be found.")
	}

	transaction.EvaluateTransactionSuccess(err)
	transaction.Finish()

	result1, _ := bufferManager.GetPages(keyMock)

	if len(result1) != 1 {
		t.Errorf("Page that is not dirty should be found.")
	}
}

func Test_IfCannotReadDirtyPageRecords(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)
	bufferManager, _ := NewBufferManager()
	keyMock := "456789"

	transaction_1 := NewTransaction()
	transaction_1.Start()
	err_1 := bufferManager.AddPlayer(keyMock, initElo, nil, &transaction_1)
	transaction_1.EvaluateTransactionSuccess(err_1)
	transaction_1.Finish()

	result_1, _ := bufferManager.GetPages(keyMock)

	if len(result_1) != 1 {
		t.Errorf("Page should be created, since it is not dirty")
	}
	if result_1[len(result_1)-1].Body[0].transactionStatus != success {
		t.Errorf("Initial record should be visible")
	}

	transaction_2 := NewTransaction()
	transaction_2.Start()
	err_2 := bufferManager.UpdatePlayer(keyMock, 25, nil, &transaction_2)

	result_2, _ := bufferManager.GetPages(keyMock)

	if result_2[len(result_2)-1].Body[1].transactionStatus != inProgress {
		t.Errorf("Only initial record should be visible")
	}

	transaction_2.EvaluateTransactionSuccess(err_2)
	transaction_2.Finish()

	result_3, _ := bufferManager.GetPages(keyMock)

	if result_3[len(result_3)-1].Body[1].transactionStatus != success {
		t.Errorf("New record should also be visible")
	}
}

func Test_IfMatchesChangesStatusAfterTransactionFinish(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := NewBufferManager()
	teamOneKeys := []string{"1", "2"}
	teamTwoKeys := []string{"3", "4"}
	matchResult := 0

	transaction_1 := NewTransaction()
	transaction_1.Start()
	result := bufferManager.AddMatchEvent(teamOneKeys, teamTwoKeys, int8(matchResult), &transaction_1)

	if result.transactionStatus != inProgress {
		t.Errorf("match status should be inProgress.")
	}

	transaction_1.EvaluateTransactionSuccess(nil)
	transaction_1.Finish()

	if result.transactionStatus != success {
		t.Errorf("match status should be success.")
	}
}
