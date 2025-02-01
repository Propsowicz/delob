package processor

import (
	buffer "delob/internal/buffer"
	tokenizer "delob/internal/tokenizer"
	"delob/internal/utils"
	"encoding/json"
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

func (p *Processor) handleOrders(orders interface{}) (string, error) {
	var result string = ""
	var orderError error

	switch v := orders.(type) {
	case tokenizer.AddPlayersToken:
		result, orderError = p.addPlayer(v.Keys)
		if orderError != nil {
			return result, orderError
		}
	case tokenizer.UpdatePlayersToken:
		result, orderError = p.updatePlayers(v.WinKeys, v.LoseKeys)
		if orderError != nil {
			return result, orderError
		}
	case tokenizer.SelectAllToken:
		result, orderError = p.selectAll()
		if orderError != nil {
			return result, orderError
		}
	default:
		fmt.Printf("unexpected type %T", v)
	}

	return result, nil
}

func (p *Processor) selectAll() (string, error) {
	allEntities, pagesCollections, err := p.bufferManager.GetAllPages()
	if err != nil {
		return "", err
	}
	playersCollection := []Player{}

	for i := 0; i < len(allEntities); i++ {
		playersCollection = append(playersCollection,
			newPlayer(allEntities[i], pagesCollections[i]))
	}
	jsonResult, errMarshal := json.Marshal(playersCollection)
	if errMarshal != nil {
		return "", errMarshal
	}

	return string(jsonResult), nil
}

func (p *Processor) updatePlayers(winKeys []string, loseKeys []string) (string, error) {

	winPlayers, losePlayers, err := p.loadPlayersToUpdate(winKeys, loseKeys)

	if err != nil {
		return "", err
	}

	calc := NewCalculator(playerWin, playerLose)

	err1 := p.bufferManager.UpdatePlayer(playerWin.Id, calc.GetWinElo())
	if err1 != nil {
		return "", err1
	}

	err2 := p.bufferManager.UpdatePlayer(playerLose.Id, calc.GetLoseElo())
	if err2 != nil {
		return "", err2
	}

	return affectNumberOfRowsMessage(2), nil
}

func (p *Processor) loadPlayersToUpdate(winKeys []string, loseKeys []string) ([]Player, []Player, error) {
	winPlayers, losePlayers := []Player{}, []Player{}

	for i := range winKeys {
		player, errWin := p.getPlayerById(winKeys[i])
		if errWin != nil {
			return winPlayers, losePlayers, errWin
		}
		winPlayers = append(winPlayers, player)
	}

	for i := range loseKeys {
		player, errLose := p.getPlayerById(loseKeys[i])
		if errLose != nil {
			return winPlayers, losePlayers, errLose
		}
		losePlayers = append(losePlayers, player)
	}
	return winPlayers, losePlayers, nil
}

func (p *Processor) getPlayerById(entityId string) (Player, error) {
	pages, err := p.bufferManager.GetPages(entityId)
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
			playerLose = args[i+1]
			errorChecker++
		}
	}

	if errorChecker != 2 {
		return playerWin, playerLose, fmt.Errorf("at least on of entities does not exists")
	}

	return playerWin, playerLose, nil
}

func (p *Processor) addPlayer(order []string) (string, error) {
	var numberOfAddedPlayers int16
	var isFullySuccessful bool = true
	var invalidEntityIds []string

	for i := range order {
		err := p.bufferManager.AddPlayer(order[i], utils.INITIAL_ELO)
		if err != nil {
			isFullySuccessful = false
			invalidEntityIds = append(invalidEntityIds, order[i])
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
