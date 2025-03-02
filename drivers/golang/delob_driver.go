package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type DelobContext struct {
	tcpHandler *tcpHandler
}

type Player struct {
	Key string
	Elo int
}

type OrderKey string

const (
	Key OrderKey = "Key"
	Elo OrderKey = "Elo"
)

type OrderDirection string

const (
	Ascending  OrderDirection = "ASC"
	Descending OrderDirection = "DESC"
)

func NewContext(connectionString string) (DelobContext, error) {
	tcpHandler, err := newTcpHandler(connectionString)
	if err != nil {
		return DelobContext{}, err
	}

	return DelobContext{
		tcpHandler: tcpHandler,
	}, nil
}

func (c *DelobContext) AddPlayer(playerKey string) error {
	expression := fmt.Sprintf("ADD PLAYER '%s';", playerKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) AddPlayers(playerKeys []string) error {
	expression := fmt.Sprintf("ADD PLAYERS %s;", creteCollectionFromArray(playerKeys))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDecisiveTeamMatch(teamWinKeys, teamLoseKeys []string) error {
	expression := fmt.Sprintf("SET WIN FOR %s AND LOSE FOR %s;", creteCollectionFromArray(teamWinKeys), creteCollectionFromArray(teamLoseKeys))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDecisiveMatch(playerOneKey, playerTwoKey string) error {
	expression := fmt.Sprintf("SET WIN FOR '%s' AND LOSE FOR '%s';", playerOneKey, playerTwoKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDrawTeamMatch(teamOne, teamTwo []string) error {
	expression := fmt.Sprintf("SET DRAW BETWEEN %s AND %s;", creteCollectionFromArray(teamOne), creteCollectionFromArray(teamTwo))

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) SetDrawMatch(playerOneKey, playerTwoKey string) error {
	expression := fmt.Sprintf("SET DRAW BETWEEN '%s' AND '%s';", playerOneKey, playerTwoKey)

	_, result := c.sendMessage(expression)
	return result
}

func (c *DelobContext) GetPlayers() ([]Player, error) {
	expression := "SELECT Players;"

	return c.getPlayersQuery(expression)
}

func (c *DelobContext) GetPlayersOrderBy(orderKey OrderKey, orderDirection OrderDirection) ([]Player, error) {
	expression := fmt.Sprintf("SELECT Players ORDER BY %s %s;", orderKey, orderDirection)

	return c.getPlayersQuery(expression)
}

func (c *DelobContext) getPlayersQuery(expression string) ([]Player, error) {
	jsonResponse, errDelob := c.sendMessage(expression)
	if errDelob != nil {
		return nil, errDelob
	}

	result := []Player{}
	err := json.Unmarshal([]byte(jsonResponse), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func creteCollectionFromArray(arr []string) string {
	result := []string{}
	for i := range arr {
		result = append(result, fmt.Sprintf("'%s'", arr[i]))
	}

	return fmt.Sprintf("(%s)", strings.Join(result, ","))
}

func (c *DelobContext) sendMessage(expression string) (string, error) {
	response, err := c.tcpHandler.sendMessage(expression)
	if err != nil {
		return "", err
	}

	return response, nil
}

// tcp connection handler

type tcpHandler struct {
	connectionString connectionString
	protocolVersion  string
	conn             net.Conn
	reader           *bufio.Reader
	writer           *bufio.Writer
}

func newTcpHandler(rawConnectionString string) (*tcpHandler, error) {
	connectionString, errConStr := parseConnectionString(rawConnectionString)
	if errConStr != nil {
		return nil, errConStr
	}

	conn, err := net.Dial("tcp", connectionString.adress)
	if err != nil {
		return nil, err
	}

	return &tcpHandler{
		connectionString: connectionString,
		protocolVersion:  "00", // TODO
		conn:             conn,
		reader:           bufio.NewReader(conn),
		writer:           bufio.NewWriter(conn),
	}, nil
}

func (h *tcpHandler) sendMessage(message string) (string, error) {
	_, err := h.writer.WriteString(message + " \n")
	if err != nil {
		return "", err
	}
	if err := h.writer.Flush(); err != nil {
		return "", err
	}

	response, err := h.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	rawResponse := strings.TrimSpace(response)
	protocolVersion := rawResponse[0:2]
	executionState := rawResponse[2:3]

	if protocolVersion != h.protocolVersion {
		return "", fmt.Errorf("wrong protocol version")
	}
	if executionState == "0" {
		return "", fmt.Errorf("expression execution was not successful: %s", strings.TrimSpace(response)[3:])
	}

	return strings.TrimSpace(response)[3:], nil
}

// connection string

type connectionString struct {
	server   string
	port     string
	adress   string
	username string
	password string
}

func parseConnectionString(rawConnectionString string) (connectionString, error) {
	const defaultPort string = "5678"
	const serverKey string = "server"
	const portKey string = "port"
	const uidKey string = "uid"
	const pwdKey string = "pwd"

	tokens := strings.Split(rawConnectionString, ";")
	connectionString := connectionString{}

	for i := range tokens {
		tokenKeyValue := strings.Split(tokens[i], "=")
		switch strings.ToLower(tokenKeyValue[0]) {
		case serverKey:
			connectionString.server = tokenKeyValue[1]
		case portKey:
			connectionString.port = tokenKeyValue[1]
		case uidKey:
			connectionString.username = tokenKeyValue[1]
		case pwdKey:
			connectionString.password = tokenKeyValue[1]
		}
	}
	if connectionString.port == "" {
		connectionString.port = defaultPort
	}

	if err := validateConnectionStringElement(connectionString.server, serverKey); err != nil {
		return connectionString, err
	}
	if err := validateConnectionStringElement(connectionString.username, uidKey); err != nil {
		return connectionString, err
	}
	if err := validateConnectionStringElement(connectionString.password, pwdKey); err != nil {
		return connectionString, err
	}

	connectionString.adress = connectionString.server + ":" + connectionString.port

	return connectionString, nil
}

func validateConnectionStringElement(element, key string) error {
	if element == "" {
		return fmt.Errorf("cannot find %s element in connection string", key)
	}
	return nil
}
