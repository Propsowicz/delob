package processor

import (
	buffer "delob/internal/buffer"
	tokenizer "delob/internal/tokenizer"
	"fmt"
	"strings"
)

type Processor struct {
	bufferManager *buffer.BufferManager
}

func NewProcessor(b *buffer.BufferManager) Processor {
	return Processor{
		bufferManager: b,
	}
}

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

func (p *Processor) Execute(
	expression string) (string, error) {
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
	isTransactionSuccessful, transactionStepsTable := p.startTansaction()

	result, ordersError := p.handleOrders(tokenizedExpression)
	if ordersError == nil {
		isTransactionSuccessful = true
	}

	return result, p.finishTransaction(isTransactionSuccessful, transactionStepsTable, ordersError)

}

func (p *Processor) handleOrders(orders []tokenizer.TokenizedExpression) (string, error) {
	var result string = ""
	var orderError error

	for _, order := range orders {
		if order.ProcessMethod == tokenizer.AddPlayer {
			result, orderError = p.addPlayer(order)
			if orderError != nil {
				return result, orderError
			}
		}
		if order.ProcessMethod == tokenizer.UpdatePlayers {
			result, orderError = p.updatePlayers(order)
			if orderError != nil {
				return result, orderError
			}
		}
	}

	return result, nil
}

func (p *Processor) updatePlayers(order tokenizer.TokenizedExpression) (string, error) {
	// get players
	// calculate
	// update
	playerWin, playerLose, err := p.newPlayersPair(order.Arguments)
	if err != nil {
		return "", err
	}

	// playerWinElo, errWinElo := p.getPlayerById(playerWin)

	return playerLose.Id + playerWin.Id, nil

}

func (p *Processor) newPlayersPair(ids []string) (Player, Player, error) {
	playerWin, playerLose := Player{}, Player{}

	playerWinId, playerLoseId, err := p.extractPlayerIds(ids)
	if err != nil {
		return Player{}, Player{}, err
	}

	playerWin, errWin := p.getPlayerById(playerWinId)
	if errWin != nil {
		return playerWin, playerLose, errWin
	}

	playerLose, errLose := p.getPlayerById(playerLoseId)
	if errLose != nil {
		return playerWin, playerLose, errLose
	}

	return playerWin, playerLose, nil
}

func (p *Processor) getPlayerById(entityId string) (Player, error) {
	pages, err := p.bufferManager.GetPage(entityId)
	if err != nil {
		return Player{}, err
	}

	return newPlayer(entityId, pages), nil

}

func (p *Processor) extractPlayerIds(args []string) (string, string, error) {
	var playerWin string
	var playerLose string

	errorChecker := 0
	for i := 0; i < len(args); i++ {
		if args[i] == "WIN" {
			playerWin = args[i+1]
			errorChecker++
		}
		if args[i] == "LOSE" {
			playerWin = args[i+1]
			errorChecker++
		}
	}

	if errorChecker != 2 {
		return playerWin, playerLose, fmt.Errorf("at least on of entities does not exists")
	}

	return playerWin, playerLose, nil
}

func (p *Processor) addPlayer(order tokenizer.TokenizedExpression) (string, error) {
	var numberOfAddedPlayers int16
	var isFullySuccessful bool = true
	var invalidEntityIds []string

	for i := range order.Arguments {
		err := p.bufferManager.AddPlayer(order.Arguments[i])
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

func (p *Processor) startTansaction() (bool, *transactionSteps) {
	return false, &transactionSteps{}
}

func (p *Processor) finishTransaction(isTransactionSuccessful bool, tratransactionStepsTable *transactionSteps, err error) error {
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
