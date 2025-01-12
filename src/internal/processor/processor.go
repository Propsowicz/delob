package processor

import (
	buffer "delob/internal/buffer"
	tokenizer "delob/internal/tokenizer"
	"fmt"
	"strings"
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

func Execute(
	expression string,
	bufferManager *buffer.BufferManager) (string, error) {
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

	result, ordersError := handleOrders(tokenizedExpression, bufferManager)
	if ordersError == nil {
		isTransactionSuccessful = true
	}

	return result, finishTransaction(isTransactionSuccessful, transactionStepsTable, ordersError)

}

func handleOrders(orders []tokenizer.TokenizedExpression, bufferManager *buffer.BufferManager) (string, error) {
	var result string = ""
	var orderError error

	for _, order := range orders {
		if order.ProcessMethod == tokenizer.AddPlayer {
			result, orderError = addPlayer(order, bufferManager)
			if orderError != nil {
				return result, orderError
			}
		}
	}

	return result, nil
}

func addPlayer(order tokenizer.TokenizedExpression, bufferManager *buffer.BufferManager) (string, error) {
	var numberOfAddedPlayers int16
	var isFullySuccessful bool = true
	var invalidEntityIds []string

	for i := range order.Arguments {
		err := bufferManager.AddPlayer(order.Arguments[i])
		if err != nil {
			isFullySuccessful = false
			invalidEntityIds = append(invalidEntityIds, order.Arguments[i])
			continue
		}
		numberOfAddedPlayers++
	}

	if !isFullySuccessful {
		return affectNumberOfRowsMessage(numberOfAddedPlayers),
			fmt.Errorf("cannot add players with Ids: %s", strings.Join(invalidEntityIds, " | "))
	}

	return affectNumberOfRowsMessage(numberOfAddedPlayers), nil
}

func startTansaction() (bool, *transactionSteps) {
	return false, &transactionSteps{}
}

func finishTransaction(isTransactionSuccessful bool, tratransactionStepsTable *transactionSteps, err error) error {
	if !isTransactionSuccessful {
		revertChanges(1)

		return err
	}

	return nil
}

func revertChanges(transactionId int) error {
	// in terms of error the records should be marked - isDirty and deleted?
	return nil
}

func affectNumberOfRowsMessage(numberOfAffectedRows int16) string {
	return fmt.Sprintf("%d row(s) affected", numberOfAffectedRows)
}
