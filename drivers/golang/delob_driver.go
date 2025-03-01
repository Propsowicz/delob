package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type DelobContext struct {
	address    string
	tcpHandler *TcpHandler
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
	tcpHandler, err := newTcpHandler(5678)
	if err != nil {
		return DelobContext{}, err
	}

	return DelobContext{
		address:    "",
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

// tcp connection

type TcpHandler struct {
	port            int
	serverAddress   string
	protocolVersion string
	conn            net.Conn
	reader          *bufio.Reader
	writer          *bufio.Writer
}

func newTcpHandler(port int) (*TcpHandler, error) {
	buildEnv := os.Getenv("BUILD_ENV")
	hostAdress := "127.0.0.1"
	if buildEnv == "docker" {
		hostAdress = "0.0.0.0"
	}
	serverAddress := fmt.Sprintf("%s:%s", hostAdress, strconv.Itoa(port))

	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}

	return &TcpHandler{
		port:            port,
		serverAddress:   serverAddress,
		protocolVersion: "00", // TODO
		conn:            conn,
		reader:          bufio.NewReader(conn),
		writer:          bufio.NewWriter(conn),
	}, nil
}

func (h *TcpHandler) sendMessage(message string) (string, error) {
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
