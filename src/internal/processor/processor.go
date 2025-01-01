package processor

import (
	"delob/internal/shared/model"
	tokenizer "delob/internal/tokenizer"
)

type transactionSteps struct {
	addToPageManager addToPageManager
	addToPage        addToPage
	addToLogs        addToLogs
}

type addToPageManager struct {
}

type addToPage struct {
}

type addToLogs struct {
}

func Execute(expression string) (string, error) {
	tokenizedExpression, err := tokenizer.Tokenize(expression)
	if err != nil {
		return "", err
	}

	// transactionId -> create file that is going to collect transactions steps in case of reverting it
	// steps:
	// add to page-manager
	// save to in-memory buffer

	// should be smth different
	// transactionId := time.Now().Nanosecond()
	isTransactionSuccessful, transactionStepsTable := startTansaction()

	result, ordersError := handleOrders(tokenizedExpression)
	if ordersError == nil {
		isTransactionSuccessful = true
	}

	return result, finishTransaction(isTransactionSuccessful, transactionStepsTable)

}

func handleOrders(orders []model.TokenizedExpression) (string, error) {
	var result string = ""
	var orderError error

	for _, order := range orders {

		if order.ProcessMethod == model.AddPlayer {
			result, orderError = addPlayer(order)
			if orderError != nil {
				return result, orderError
			}
		}
	}

	return result, nil
}

func addPlayer(order model.TokenizedExpression) (string, error) {
	return "", nil
}

func startTansaction() (bool, *transactionSteps) {
	return false, &transactionSteps{}
}

func finishTransaction(isTransactionSuccessful bool, tratransactionStepsTable *transactionSteps) error {
	if !isTransactionSuccessful {
		return revertChanges(1)
	}

	return nil
}

func revertChanges(transactionId int) error {
	// in terms of error the records should be marked - isDirty and deleted?
	return nil
}
