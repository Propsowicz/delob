package buffer

import (
	"fmt"
	"sync"
)

type Transaction struct {
	status                               transactionStatus
	shouldSuccessfulyCompleteTransaction bool
	pageDictionaryPointers               []*PageData
	pagePointers                         []*Page
	pageMatchPointers                    []*Match
	syncMutex                            sync.Mutex
}

func NewTransaction() Transaction {
	return Transaction{}
}

func (t *Transaction) Start() {
	t.status = inProgress
	t.shouldSuccessfulyCompleteTransaction = false
	t.pageDictionaryPointers = []*PageData{}
	t.pagePointers = []*Page{}
	t.pageMatchPointers = []*Match{}
}

func (t *Transaction) EvaluateTransactionSuccess(err error) bool {
	if err == nil {
		t.shouldSuccessfulyCompleteTransaction = true
		return true
	}
	return false
}

func (t *Transaction) AddPageDictionaryPointer(pageDict *PageData) {
	t.pageDictionaryPointers = append(t.pageDictionaryPointers, pageDict)
}

func (t *Transaction) Finish() {
	var transactionStatus transactionStatus
	if t.shouldSuccessfulyCompleteTransaction {
		transactionStatus = success
	} else {
		transactionStatus = failed
	}

	t.syncMutex.Lock()

	for i := range t.pageDictionaryPointers {
		fmt.Println(&t.pageDictionaryPointers[i])
		t.pageDictionaryPointers[i].transactionStatus = transactionStatus
	}

	t.syncMutex.Unlock()
}

func (t *Transaction) changeTransationStatus(isSuccessful bool) {

}

func (t *Transaction) Clean() {

}
