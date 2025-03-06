package interfaces

import (
	"bufio"
	"delob/internal/auth"
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
	authManager     auth.AuthenticationManager
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
		authManager:     auth.NewAuthenticationManager(),
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

	// auth

	// s.authManager.TryAuthenticate("qwe", "12")

	defer c.Close()
	logger.Info("", fmt.Sprintf("Serving %s", c.RemoteAddr().String()))

	for {
		traceId := utils.GenerateKey()
		reader := bufio.NewReader(c)
		writer := bufio.NewWriter(c)
		readString, err := reader.ReadString('\n')
		if err != nil {
			logger.Error(traceId, err)
			return
		}

		rawRequest := strings.TrimSpace(strings.TrimSuffix(readString, "\r\n"))
		request, errReqParse := parseRequest(rawRequest, c.RemoteAddr().String())
		if errReqParse != nil {
			return
		}

		if s.authManager.TryAuthenticate(request.user, request.ip) {

			result, err := requestHandler(traceId, request.msg)
			if err != nil {
				s.writeString(*writer, s.newResponse(fail, err.Error()), traceId)
				logger.Error(traceId, err)
			} else {
				s.writeString(*writer, s.newResponse(success, result), traceId)
				logger.Info(traceId, result)
			}
		} else {
			s.writeString(*writer, s.newResponse(authChallenge, ""), traceId)
			logger.Info(traceId, "auth challenge started")

			readString, err := reader.ReadString('\n')
			if err != nil {
				logger.Error(traceId, err)
				return
			}

			rawRequest := strings.TrimSpace(strings.TrimSuffix(readString, "\r\n"))
			request, errReqParse := parseRequest(rawRequest, c.RemoteAddr().String())
			if errReqParse != nil {
				return
			}

			fmt.Println("new request")
			fmt.Println(request)

			s.writeString(*writer, s.newResponse(authChallenge, "some data"), traceId)

		}
	}
}

func (s *TcpServer) writeString(writer bufio.Writer, response, traceId string) {
	if _, err := writer.WriteString(response); err != nil {
		logger.Error(traceId, err)
	}
	if err := writer.Flush(); err != nil {
		logger.Error(traceId, err)
	}
}

type request struct {
	user string
	msg  string
	ip   string
}

func parseRequest(s, ip string) (request, error) {
	// user=<>,msg=<>

	parts := strings.Split(s, "|||")

	// for now I use this separator |||
	r := request{}
	r.user = parts[0]
	r.msg = parts[1]
	r.ip = ip
	return r, nil
}

type status int8

const (
	fail          status = 0
	success       status = 1
	authChallenge status = 9
)

func (s *TcpServer) newResponse(status status, msg string) string {
	return fmt.Sprintf("%s%d%s\n", s.protocolVersion, status, msg)
}
