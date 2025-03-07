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

func (s *TcpServer) getRequest(reader bufio.Reader, ip string) (request, error) {
	readString, err := reader.ReadString('\n')
	if err != nil {
		return request{}, err
	}

	rawRequest := strings.TrimSpace(strings.TrimSuffix(readString, "\r\n"))
	return parseRequest(rawRequest, ip)
}

func (s *TcpServer) handleConnection(c net.Conn, requestHandler requestHandler) {
	defer c.Close()
	logger.Info("", fmt.Sprintf("Serving %s", c.RemoteAddr().String()))

	for {
		traceId := utils.GenerateKey()
		reader := bufio.NewReader(c)
		writer := bufio.NewWriter(c)

		request, errReqParse := s.getRequest(*reader, c.RemoteAddr().String())
		if errReqParse != nil {
			logger.Error(traceId, errReqParse)
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
			clientFirstMessage, errClientFirst := s.getClientFirstMessage(c, *writer, *reader, traceId)
			if errClientFirst != nil {
				logger.Error(traceId, errClientFirst)
				return
			}

			auth, errAddClientFirst := s.authManager.AddClientFirstAuthString(clientFirstMessage.msg)
			if errAddClientFirst != nil {
				logger.Error(traceId, errAddClientFirst)
				return
			}

			auth, errAddClientFirstMsg := s.addServerFirstAuthString(auth, clientFirstMessage.user)
			if errAddClientFirstMsg != nil {
				logger.Error(traceId, errAddClientFirstMsg)
				return
			}

			clientFinalMessage, errClientFinal := s.getClientFinalMessage(c, *writer, *reader, traceId, auth)
			if errClientFinal != nil {
				logger.Error(traceId, errClientFinal)
				return
			}

			fmt.Printf(clientFinalMessage.user)
			fmt.Printf(clientFinalMessage.msg)

			// s.writeString(*writer, s.newResponse(authChallenge, "some data"), traceId)

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

func (s *TcpServer) getClientFirstMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId string) (request, error) {
	s.writeString(writer, s.newResponse(authChallenge, ""), traceId)

	return s.getRequest(reader, c.RemoteAddr().String())
}

func (s *TcpServer) addServerFirstAuthString(auth, user string) (string, error) {
	auth, errServerFirst := s.authManager.AddServerFirstMessage(auth, user)
	if errServerFirst != nil {
		return "", errServerFirst
	}
	return auth, nil
}

func (s *TcpServer) getClientFinalMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId, auth string) (request, error) {
	s.writeString(writer, s.newResponse(authChallenge, auth), traceId)

	return s.getRequest(reader, c.RemoteAddr().String())
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
