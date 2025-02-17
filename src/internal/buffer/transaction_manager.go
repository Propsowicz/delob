package buffer

type Transaction struct {
	status                 transactionStatus
	pageDictionaryPointers []*PageData
	pagePointers           []*Page
	pagematchPointers      []*Match
}

func NewTransaction() Transaction {
	return Transaction{}
}

func (t *Transaction) Start() {
	t.status = inProgress
}

func (t *Transaction) AddPageDictionaryPointer(pageDict *PageData) {
	t.pageDictionaryPointers = append(t.pageDictionaryPointers, pageDict)
}

func (t *Transaction) Finish(isSuccessful bool) {
	for i := range t.pageDictionaryPointers {
		t.pageDictionaryPointers[i].transactionStatus = success
	}
}

func (t *Transaction) changeTransationStatus(isSuccessful bool) {

}

func (t *Transaction) clean(isSuccessful bool) {

}
