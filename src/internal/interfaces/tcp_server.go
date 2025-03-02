package interfaces

import (
	"bufio"
	"delob/internal/utils"
	"delob/internal/utils/logger"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

// protocol:
// version, response status, msg

type TcpServer struct {
	port            int
	serverAddress   string
	protocolVersion string
}

type requestHandler func(string, string) (string, error)

func NewTcpServer(port int) TcpServer {
	buildEnv := os.Getenv("BUILD_ENV")
	hostAdress := "127.0.0.1"
	if buildEnv == "docker" {
		hostAdress = "0.0.0.0"
	}
	serverAddress := fmt.Sprintf("%s:%s", hostAdress, strconv.Itoa(port))

	return TcpServer{
		port:            port,
		serverAddress:   serverAddress,
		protocolVersion: "00", // TODO
	}
}

func (s *TcpServer) Start(requestHandler requestHandler) {
	l, err := net.Listen("tcp4", s.serverAddress)
	if err != nil {
		logger.Error("", err)
		return
	}

	logger.Info("", "Started listening for tcp connections on: "+s.serverAddress)

	defer l.Close()
	for {
		c, err := l.Accept()
		if err != nil {
			logger.Error("", err)
			return
		}
		go s.handleConnection(c, requestHandler)
	}
}

func (s *TcpServer) handleConnection(c net.Conn, requestHandler requestHandler) {
	defer c.Close()
	logger.Info("", fmt.Sprintf("Serving %s", c.RemoteAddr().String()))

	for {
		traceId := utils.GenerateKey()
		reader := bufio.NewReader(c)
		writer := bufio.NewWriter(c)
		bufferData, err := reader.ReadString('\n')
		if err != nil {
			logger.Error(traceId, err)
			return
		}

		rawExpression := strings.TrimSpace(strings.TrimSuffix(bufferData, "\r\n"))

		result, err := requestHandler(traceId, rawExpression)

		if err != nil {
			s.writeString(*writer, false, traceId, err.Error())
			logger.Error(traceId, err)
		} else {
			s.writeString(*writer, true, traceId, result)
			logger.Info(traceId, result)
		}
	}
}

func (s *TcpServer) writeString(writer bufio.Writer, isResponseSuccessfull bool, traceId, response string) {
	if _, err := writer.WriteString(s.formatResponse(isResponseSuccessfull, response)); err != nil {
		logger.Error(traceId, err)
	}
	if err := writer.Flush(); err != nil {
		logger.Error(traceId, err)
	}
}

func (s *TcpServer) formatResponse(isResponseSuccessfull bool, response string) string {
	responseStatus := "0"
	if isResponseSuccessfull {
		responseStatus = "1"
	}

	return fmt.Sprintf("%s%s%s\n", s.protocolVersion, responseStatus, response)
}
