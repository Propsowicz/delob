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
	defer c.Close()
	for {
		traceId := utils.GenerateKey()
		reader := bufio.NewReader(c)
		writer := bufio.NewWriter(c)

		request, errReqParse := s.getRequest(*reader, c.RemoteAddr().String())
		if errReqParse != nil {
			logger.Error(traceId, errReqParse)
			return
		}

		logger.Info(traceId, fmt.Sprintf("Request from user: %s and ip: %s", request.user, request.ip))

		if s.authManager.IsUserAuthenticated(request.user, request.ip) {
			result, err := requestHandler(traceId, request.msg)
			if err != nil {
				s.writeString(*writer, s.newResponse(fail, err.Error()), traceId)
				logger.Error(traceId, err)
			} else {
				s.writeString(*writer, s.newResponse(success, result), traceId)
				logger.Info(traceId, result)
			}
		} else {
			user, auth, errClientFirst := s.getClientFirstMessage(c, *writer, *reader, traceId)
			if errClientFirst != nil {
				logger.Error(traceId, errClientFirst)
				return
			}

			auth, errServerFirstAuth := s.prepareServerFirstMessage(auth, user)
			if errServerFirstAuth != nil {
				s.writeString(*writer, s.newResponse(fail, errServerFirstAuth.Error()), traceId)
				return
			}

			proofRequest, errProofRequest := s.getProofMessage(c, *writer, *reader, traceId, auth)
			if errProofRequest != nil {
				logger.Error(traceId, errProofRequest)
				return
			}

			verifyProofResult := s.verifyProof(proofRequest.msg, user, proofRequest.ip, auth)

			s.sendVerifierResult(*writer, traceId, verifyProofResult)
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

func (s *TcpServer) verifyProof(proof, user, ip, auth string) bool {
	return s.authManager.Verify(proof, user, ip, auth)
}

func (s *TcpServer) sendVerifierResult(writer bufio.Writer, traceId string, proofIsCorrect bool) {
	if proofIsCorrect {
		s.writeString(writer, s.newResponse(user_verified, ""), traceId)
	} else {
		s.writeString(writer, s.newResponse(user_not_verified, ""), traceId)
	}
}

func (s *TcpServer) getClientFirstMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId string) (string, string, error) {
	s.writeString(writer, s.newResponse(authentication_challenge, ""), traceId)

	clientFirstMessage, errClientFirst := s.getRequest(reader, c.RemoteAddr().String())
	if errClientFirst != nil {
		return "", "", errClientFirst
	}

	user, _, auth, errParseClientFirst := s.authManager.ParseClientFirstMessageToAuthString(clientFirstMessage.msg)
	if errParseClientFirst != nil {
		return "", "", errParseClientFirst
	}

	if user != clientFirstMessage.user {
		return "", "", fmt.Errorf("user from request and auth string does not match")
	}

	return clientFirstMessage.user, auth, nil
}

func (s *TcpServer) prepareServerFirstMessage(auth, user string) (string, error) {
	auth, errServerFirst := s.authManager.PrepareServerFirstMessage(auth, user)
	if errServerFirst != nil {
		return "", errServerFirst
	}
	return auth, nil
}

func (s *TcpServer) getProofMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId, auth string) (request, error) {
	s.writeString(writer, s.newResponse(authentication_challenge, auth), traceId)

	return s.getRequest(reader, c.RemoteAddr().String())
}
func (s *TcpServer) getRequest(reader bufio.Reader, ip string) (request, error) {
	readString, err := reader.ReadString('\n')
	if err != nil {
		return request{}, err
	}

	rawRequest := strings.TrimSpace(strings.TrimSuffix(readString, "\r\n"))
	return parseRequest(rawRequest, ip)
}
