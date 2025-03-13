package buffer

import (
	"delob/internal/utils"
	"os"
	"testing"
)

func setupSuite(_ *testing.T) func(t *testing.T) {
	backupManagerPath := "log_data"
	os.RemoveAll(backupManagerPath)
	return func(t *testing.T) {
		backupManagerPath := "log_data"
		os.RemoveAll(backupManagerPath)
	}
}

const updateEloValue int16 = 25
const initElo int16 = 1500

var transactionMock Transaction = NewTransaction()

func Test_IfCanUpdatePlayerWhenCanAddToExistingPage(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := NewBufferManager()
	keyMock := "123987"

	transaction_1 := NewTransaction()
	transaction_1.Start()
	err := bufferManager.AddPlayer(keyMock, initElo, nil, &transaction_1)
	transaction_1.EvaluateTransactionSuccess(err)
	transaction_1.Finish()

	transaction_2 := NewTransaction()
	transaction_2.Start()
	err1 := bufferManager.UpdatePlayer(keyMock, updateEloValue, nil, &transaction_2)
	transaction_2.EvaluateTransactionSuccess(err1)
	transaction_2.Finish()

	if err1 != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(keyMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Value != initElo {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[1].Value != updateEloValue {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenThereIsOnlyOneSlotInPage(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := NewBufferManager()
	keyMock := "123987"

	transaction_1 := NewTransaction()
	transaction_1.Start()
	err_1 := bufferManager.AddPlayer(keyMock, initElo, nil, &transaction_1)
	transaction_1.EvaluateTransactionSuccess(err_1)
	transaction_1.Finish()

	for i := 0; i < int(utils.PAGE_SIZE)-2; i++ {
		transaction_1 = NewTransaction()
		transaction_1.Start()
		err_1 = bufferManager.UpdatePlayer(keyMock, updateEloValue, nil, &transaction_1)
		transaction_1.EvaluateTransactionSuccess(err_1)
		transaction_1.Finish()
	}

	transaction_2 := NewTransaction()
	transaction_2.Start()
	err_2 := bufferManager.UpdatePlayer(keyMock, updateEloValue, nil, &transaction_2)
	transaction_2.EvaluateTransactionSuccess(err_2)
	transaction_2.Finish()

	if err_2 != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(keyMock)

	if len(result) != 1 {
		t.Errorf("There should be one page.")
	}
	if result[0].Body[0].Value != initElo {
		t.Errorf("First record should stay the same.")
	}
	if result[0].Body[utils.PAGE_SIZE-1].Value != updateEloValue {
		t.Errorf("New record should be added.")
	}
}

func Test_IfCanUpdatePlayerWhenPageIsFull(t *testing.T) {
	teardownSuite := setupSuite(t)
	defer teardownSuite(t)

	bufferManager, _ := NewBufferManager()
	keyMock := "123987"

	transaction_1 := NewTransaction()
	transaction_1.Start()
	err_1 := bufferManager.AddPlayer(keyMock, initElo, nil, &transaction_1)
	transaction_1.EvaluateTransactionSuccess(err_1)
	transaction_1.Finish()

	for i := 0; i < int(utils.PAGE_SIZE)-1; i++ {
		transaction_1 = NewTransaction()
		transaction_1.Start()
		err_1 = bufferManager.UpdatePlayer(keyMock, updateEloValue, nil, &transaction_1)
		transaction_1.EvaluateTransactionSuccess(err_1)
		transaction_1.Finish()
	}

	transaction_2 := NewTransaction()
	transaction_2.Start()
	err_2 := bufferManager.UpdatePlayer(keyMock, updateEloValue, nil, &transaction_2)
	transaction_2.EvaluateTransactionSuccess(err_2)
	transaction_2.Finish()

	if err_2 != nil {
		t.Errorf("Should update player.")
	}

	result, _ := bufferManager.GetPages(keyMock)

	if len(result) != 2 {
		t.Errorf("There should be two pages.")
	}
	if result[1].Body[0].Value != updateEloValue {
		t.Errorf("New record should be added.")
	}
}
