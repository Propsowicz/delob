package buffer

import (
	p "delob/internal/buffer/persistence"
	"sync"
)

type transactionStatus int8

const (
	failed     = 0
	success    = 1
	inProgress = 2
)

type BufferManager struct {
	logDataManager *p.LogsPersistenceManager
	syncMutex      sync.Mutex
	pageDictionary PageDictionary
	pages          []*Page
	matches        []*Match
}

func NewBufferManager() (BufferManager, error) {
	logDataManager, err := p.NewLogsPersistenceManager()
	if err != nil {
		return BufferManager{}, err
	}

	return BufferManager{
		logDataManager: logDataManager,
	}, nil
}

func (buffer *BufferManager) AppendLogsToPersistenceFile(traceId, parsedExpressionType,
	parsedExpression string) error {
	return buffer.logDataManager.Append(p.NewDataLog(traceId, parsedExpressionType, parsedExpression))
}

func (buffer *BufferManager) LoadLogsFromPersistenceFile() ([]p.Log, error) {
	if !buffer.logDataManager.LogsFileExists {
		return []p.Log{}, nil
	}

	result, err := buffer.logDataManager.Read()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (buffer *BufferManager) AddPlayer(key string, value int16, matchRef *Match, transaction *Transaction) error {
	if _, err := buffer.getPageAdresses(key, isInProgressOrSuccess); err == nil {
		return errorPlayerAlreadyExists(key)
	}

	pageAdress := buffer.addPage(key, value, matchRef, transaction)

	buffer.syncMutex.Lock()

	pageDictAdress, err := buffer.addPageToDictionary(key, pageAdress)

	if err != nil {
		return err
	}

	transaction.AddPageDictionaryPointer(pageDictAdress)

	buffer.syncMutex.Unlock()
	return nil
}

func (buffer *BufferManager) UpdatePlayersElo(keys []string, eloLambda int16, matchRef *Match, transaction *Transaction) error {
	for i := range keys {
		err := buffer.UpdatePlayer(keys[i], eloLambda, matchRef, transaction)
		if err != nil {
			return err
		}
	}
	return nil
}

func (buffer *BufferManager) UpdatePlayer(key string, value int16, matchRef *Match, transaction *Transaction) error {
	isAddedToExistingPage, errTryToAppend := buffer.tryAppendToPage(key, value, matchRef, transaction)
	if errTryToAppend != nil {
		return errTryToAppend
	}

	if !isAddedToExistingPage {
		pageAdress := buffer.addPage(key, value, matchRef, transaction)

		buffer.syncMutex.Lock()

		errAppendPageToExistingKey := buffer.appendPageToExistingId(key, pageAdress)
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
