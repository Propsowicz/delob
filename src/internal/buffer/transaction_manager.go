package buffer

import (
	"sync"
)

type Transaction struct {
	status                               transactionStatus
	shouldSuccessfulyCompleteTransaction bool
	pageDictionaryPointers               []*PageData
	recordPointers                       []*Record
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
	t.recordPointers = []*Record{}
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

func (t *Transaction) AddPageRecordPointer(record *Record) {
	t.recordPointers = append(t.recordPointers, record)
}

func (t *Transaction) Finish() {
	var transactionStatus transactionStatus
	if t.shouldSuccessfulyCompleteTransaction {
		transactionStatus = success
	} else {
		transactionStatus = failed
	}

	for i := range t.pageDictionaryPointers {
		t.syncMutex.Lock()
		t.pageDictionaryPointers[i].transactionStatus = transactionStatus
		t.syncMutex.Unlock()
	}

	for i := range t.recordPointers {

		t.recordPointers[i].transactionStatus = transactionStatus
	}
}

func (t *Transaction) changeTransationStatus(isSuccessful bool) {

}

func (t *Transaction) Clean() {

}
