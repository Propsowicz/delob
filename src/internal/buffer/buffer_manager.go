package buffer

import (
	"fmt"
	"sync"
)

type transactionStatus int8

const (
	failed     = 0
	success    = 1
	inProgress = 2
)

type BufferManager struct {
	logDataManager *DataLogsDictionaryManager
	syncMutex      sync.Mutex
	pageDictionary PageDictionary
	pages          []*Page
	matches        []*Match
}

func NewBufferManager() (BufferManager, error) {
	logDataManager, err := NewDataLogsDictionaryManager()
	if err != nil {
		return BufferManager{}, err
	}

	return BufferManager{
		logDataManager: logDataManager,
	}, nil
}

// append to persistance file
func (buffer *BufferManager) AppendToDataLogsDictionary(traceId, parsedExpressionType,
	parsedExpression string) error {
	return buffer.logDataManager.Append(NewDataLog(traceId, parsedExpressionType, parsedExpression))
}

// get from persistance file
func (buffer *BufferManager) LoadFromDataLogsDictionary() ([]DataLog, error) {
	if !buffer.logDataManager.IsLogsDictionaryFileExists {
		return []DataLog{}, nil
	}

	result, err := buffer.logDataManager.Read()
	if err != nil {
		return nil, nil
	}
	return result, nil
}

func (buffer *BufferManager) AddPlayer(entityId string, value int16, matchRef *Match, transaction *Transaction) error {
	if _, err := buffer.getPageAdresses(entityId, isInProgressOrSuccess); err == nil {
		return fmt.Errorf("player already exists: %s", entityId)
	}

	pageAdress := buffer.addPage(entityId, value, matchRef, transaction)

	buffer.syncMutex.Lock()

	pageDictAdress, err := buffer.addPageToDictionary(entityId, pageAdress)

	if err != nil {
		return err
	}

	transaction.AddPageDictionaryPointer(pageDictAdress)

	buffer.syncMutex.Unlock()
	return nil
}

func (buffer *BufferManager) UpdatePlayersElo(entityKeys []string, eloLambda int16, matchRef *Match, transaction *Transaction) error {
	for i := range entityKeys {
		err := buffer.UpdatePlayer(entityKeys[i], eloLambda, matchRef, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (buffer *BufferManager) UpdatePlayer(entityId string, value int16, matchRef *Match, transaction *Transaction) error {
	isAddedToExistingPage, errTryToAppend := buffer.tryAppendToPage(entityId, value, matchRef, transaction)
	if errTryToAppend != nil {
		return errTryToAppend
	}

	if !isAddedToExistingPage {
		pageAdress := buffer.addPage(entityId, value, matchRef, transaction)

		buffer.syncMutex.Lock()

		errAppendPageToExistingKey := buffer.appendPageToExistingId(entityId, pageAdress)
		if errAppendPageToExistingKey != nil {
			return errAppendPageToExistingKey
		}

		buffer.syncMutex.Unlock()
	}

	return nil
}

func (buffer *BufferManager) AddMatchEvent(teamOneKeys []string, teamTwoKeys []string, matchResult int8, transaction *Transaction) *Match {
	newMatch := newMatch(teamOneKeys, teamTwoKeys, matchResult, transaction)

	buffer.matches = append(buffer.matches, newMatch)
	return buffer.matches[len(buffer.matches)-1]
}
