package processor

import (
	buffer "delob/internal/buffer"
	elo "delob/internal/processor/elo"
	dto "delob/internal/processor/model"
	tokenizer "delob/internal/tokenizer"
	"encoding/json"
	"fmt"
	"sort"
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
	case tokenizer.AddPlayersOrder:
		result, orderError = p.addPlayer(v.Keys)
		if orderError != nil {
			return result, orderError
		}
	case tokenizer.AddMatchOrder:
		result, orderError = p.updatePlayers(v)
		if orderError != nil {
			return result, orderError
		}
	case tokenizer.SelectOrder:
		result, orderError = p.selectPlayers(v)
		if orderError != nil {
			return result, orderError
		}
	default:
		fmt.Printf("unexpected token type %T", v)
	}

	return result, nil
}

func (p *Processor) selectPlayers(selectOrder tokenizer.SelectOrder) (string, error) {
	allEntities, pagesCollections, err := p.bufferManager.GetAllPages()
	if err != nil {
		return "", err
	}
	playersCollection := []dto.Player{}

	for i := 0; i < len(allEntities); i++ {
		playersCollection = append(playersCollection, dto.NewPlayer(allEntities[i], pagesCollections[i]))
	}

	sort.Slice(playersCollection, func(i, j int) bool {
		switch selectOrder.OrderBy {
		case tokenizer.Key:
			return sortComparer(selectOrder.OrderDir == tokenizer.OrderDir(tokenizer.Asc), playersCollection[i].Key, playersCollection[j].Key)
		case tokenizer.Elo:
			return sortComparer(selectOrder.OrderDir == tokenizer.OrderDir(tokenizer.Asc), playersCollection[i].Elo, playersCollection[j].Elo)
		}
		return sortComparer(true, playersCollection[i].Key, playersCollection[j].Key)
	})

	jsonResult, errMarshal := json.Marshal(playersCollection)
	if errMarshal != nil {
		return "", errMarshal
	}

	return string(jsonResult), nil
}

func sortComparer[T int16 | string](isAsc bool, leftOperand, rightOperand T) bool {
	if isAsc {
		return leftOperand < rightOperand
	}
	return leftOperand > rightOperand
}

func (p *Processor) updatePlayers(addMatchOrder tokenizer.AddMatchOrder) (string, error) {
	teamOnePlayers, teamTwoPlayers, err := p.loadPlayersToUpdate(addMatchOrder.TeamOneKeys, addMatchOrder.TeamTwoKeys)
	if err != nil {
		return "", err
	}
	teamOneKeys, teamTwoKeys := dto.MapPlayerToKeysCollection(teamOnePlayers), dto.MapPlayerToKeysCollection(teamTwoPlayers)

	match := p.bufferManager.AddMatchEvent(teamOneKeys, teamTwoKeys, int8(addMatchOrder.MatchResult))

	calc := elo.NewCalculator(teamOnePlayers, teamTwoPlayers, addMatchOrder.MatchResult)

	errTeamOneUpdate := p.bufferManager.UpdatePlayersElo(teamOneKeys, calc.TeamOneEloLambda(), match)
	if errTeamOneUpdate != nil {
		return "", errTeamOneUpdate
	}
	errTeamTwoUpdate := p.bufferManager.UpdatePlayersElo(teamTwoKeys, calc.TeamTwoEloLambda(), match)
	if errTeamTwoUpdate != nil {
		return "", errTeamTwoUpdate
	}

	return affectNumberOfRowsMessage(int16(len(teamOnePlayers) + len(teamTwoPlayers))), nil
}

func (p *Processor) loadPlayersToUpdate(teamOneKeys []string, teamTwoKeys []string) ([]dto.Player, []dto.Player, error) {
	teamOnePlayers, teamTwoPlayers := []dto.Player{}, []dto.Player{}

	for i := range teamOneKeys {
		player, errWin := p.getPlayerByKey(teamOneKeys[i])
		if errWin != nil {
			return teamOnePlayers, teamTwoPlayers, errWin
		}
		teamOnePlayers = append(teamOnePlayers, player)
	}

	for i := range teamTwoKeys {
		player, errLose := p.getPlayerByKey(teamTwoKeys[i])
		if errLose != nil {
			return teamOnePlayers, teamTwoPlayers, errLose
		}
		teamTwoPlayers = append(teamTwoPlayers, player)
	}
	return teamOnePlayers, teamTwoPlayers, nil
}

func (p *Processor) getPlayerByKey(entityId string) (dto.Player, error) {
	pages, err := p.bufferManager.GetPages(entityId)
	if err != nil {
		return dto.Player{}, err
	}
	return dto.NewPlayer(entityId, pages), nil
}

func (p *Processor) addPlayer(order []string) (string, error) {
	var numberOfAddedPlayers int16
	var isFullySuccessful bool = true
	var invalidEntityIds []string

	for i := range order {
		err := p.bufferManager.AddPlayer(order[i], elo.INITIAL_ELO, nil)
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
