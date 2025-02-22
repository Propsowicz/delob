package processor

import (
	buffer "delob/internal/buffer"
	parser "delob/internal/parser"
	elo "delob/internal/processor/elo"
	dto "delob/internal/processor/model"
	"delob/internal/utils/logger"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type Processor struct {
	bufferManager *buffer.BufferManager
}

func NewProcessor(b *buffer.BufferManager) Processor {
	p := Processor{
		bufferManager: b,
	}

	return p
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

func (p *Processor) Initialize() error {
	dataLogs, err := p.bufferManager.LoadFromDataLogsDictionary()
	if err != nil {
		return err
	}

	for i := range dataLogs {
		transaction := buffer.NewTransaction()
		transaction.Start()
		parsedExpression, err := parser.ParseDataLogJson(dataLogs[i].ParsedExpressionType,
			dataLogs[i].ParsedExpression)

		if err != nil {
			return err
		}

		_, _, errOrder := p.handleOrders(parsedExpression, &transaction)

		if !transaction.EvaluateTransactionSuccess(errOrder) {
			return errOrder
		}

		transaction.Finish()
	}
	return nil
}

func (p *Processor) Execute(
	traceId, expression string) (string, error) {
	logger.Info(traceId, fmt.Sprintf("Processing expression: %s", expression))

	parsedExpression, err := parser.ParseExpression(traceId, expression)

	if err != nil {
		return "", err
	}

	transaction := buffer.NewTransaction()
	transaction.Start()

	result, isWriteOperation, orderError := p.handleOrders(parsedExpression, &transaction)

	return result, p.finishTransaction(
		traceId,
		parsedExpression,
		isWriteOperation,
		orderError,
		&transaction)
}

func (p *Processor) handleOrders(parsedExpression parser.ParsedExpression, transaction *buffer.Transaction) (string, bool, error) {
	var result string
	var orderError error
	var isWriteOperation bool = false

	switch parsedExpression.GetType() {
	case parser.AddPlayersCommandType:
		result, orderError = p.addPlayer(parsedExpression.(parser.AddPlayersCommand).Keys, transaction)
		isWriteOperation = true
		if orderError != nil {
			return result, false, orderError
		}
	case parser.AddMatchCommandType:
		result, orderError = p.updatePlayers(parsedExpression.(parser.AddMatchCommand), transaction)
		isWriteOperation = true
		if orderError != nil {
			return result, false, orderError
		}
	case parser.SelectQueryType:
		result, orderError = p.selectPlayers(parsedExpression.(parser.SelectQuery))
		if orderError != nil {
			return result, false, orderError
		}
	default:
		fmt.Printf("unexpected parsed order type")
	}

	return result, isWriteOperation, nil
}

func (p *Processor) selectPlayers(selectOrder parser.SelectQuery) (string, error) {
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
		case parser.Key:
			return sortComparer(selectOrder.OrderDir == parser.OrderDir(parser.Asc), playersCollection[i].Key, playersCollection[j].Key)
		case parser.Elo:
			return sortComparer(selectOrder.OrderDir == parser.OrderDir(parser.Asc), playersCollection[i].Elo, playersCollection[j].Elo)
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

func (p *Processor) updatePlayers(addMatchOrder parser.AddMatchCommand, transaction *buffer.Transaction) (string, error) {
	teamOnePlayers, teamTwoPlayers, err := p.loadPlayersToUpdate(addMatchOrder.TeamOneKeys, addMatchOrder.TeamTwoKeys)
	if err != nil {
		return "", err
	}
	teamOneKeys, teamTwoKeys := dto.MapPlayerToKeysCollection(teamOnePlayers), dto.MapPlayerToKeysCollection(teamTwoPlayers)

	match := p.bufferManager.AddMatchEvent(teamOneKeys, teamTwoKeys, int8(addMatchOrder.MatchResult), transaction)

	calc := elo.NewCalculator(teamOnePlayers, teamTwoPlayers, addMatchOrder.MatchResult)

	errTeamOneUpdate := p.bufferManager.UpdatePlayersElo(teamOneKeys, calc.TeamOneEloLambda(), match, transaction)
	if errTeamOneUpdate != nil {
		return "", errTeamOneUpdate
	}
	errTeamTwoUpdate := p.bufferManager.UpdatePlayersElo(teamTwoKeys, calc.TeamTwoEloLambda(), match, transaction)
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

func (p *Processor) addPlayer(order []string, transaction *buffer.Transaction) (string, error) {
	var numberOfAddedPlayers int16
	var isFullySuccessful bool = true
	var invalidEntityIds []string

	for i := range order {
		err := p.bufferManager.AddPlayer(order[i], elo.INITIAL_ELO, nil, transaction)
		if err != nil {
			isFullySuccessful = false
			invalidEntityIds = append(invalidEntityIds, order[i])
			continue
		}
		numberOfAddedPlayers++
	}

	if !isFullySuccessful {
		return "", fmt.Errorf("cannot add players with Ids: %s", strings.Join(invalidEntityIds, " | "))
	}
	return affectNumberOfRowsMessage(numberOfAddedPlayers), nil
}

func (p *Processor) finishTransaction(
	traceId string,
	parsedExpression parser.ParsedExpression,
	isWriteOperation bool,
	orderError error,
	transaction *buffer.Transaction) error {

	if !transaction.EvaluateTransactionSuccess(orderError) {
		return orderError
	}
	json, err := parsedExpression.ToJson()
	if err != nil {
		return err
	}

	if isWriteOperation {
		errWriteToLogsDict := p.bufferManager.AppendToDataLogsDictionary(traceId, parsedExpression.GetStringType(), json)
		if errWriteToLogsDict != nil {
			return errWriteToLogsDict
		}
		transaction.Finish()
	}

	return nil
}

func affectNumberOfRowsMessage(numberOfAffectedRows int16) string {
	return fmt.Sprintf("%d row(s) affected", numberOfAffectedRows)
}
