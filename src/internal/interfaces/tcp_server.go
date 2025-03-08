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
			user, auth, errClientFirst := s.getClientFirstMessage(c, *writer, *reader, traceId)
			if errClientFirst != nil {
				logger.Error(traceId, errClientFirst)
				return
			}

			auth, errServerFirstAuth := s.prepareServerFirstMessage(auth, user)
			if errServerFirstAuth != nil {
				return
			}

			proofRequest, errProofRequest := s.getProofMessage(c, *writer, *reader, traceId, auth)
			if errProofRequest != nil {
				logger.Error(traceId, errProofRequest)
				return
			}

			verifyProofResult := s.verifyProof(proofRequest.msg, user, auth)

			s.sendVerifierResult(c, *writer, *reader, traceId, verifyProofResult)

			fmt.Println(proofRequest.msg)
			fmt.Println(verifyProofResult)

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

func (s *TcpServer) verifyProof(proof, user, auth string) bool {
	return s.authManager.Verify(proof, user, auth)
}

func (s *TcpServer) sendVerifierResult(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId string, proofIsCorrect bool) {
	if proofIsCorrect {
		// save to session
		s.writeString(writer, s.newResponse(proofVerified, ""), traceId)
	} else {
		s.writeString(writer, s.newResponse(proofNotVerified, ""), traceId)
	}
}

func (s *TcpServer) getClientFirstMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId string) (string, string, error) {
	s.writeString(writer, s.newResponse(authChallenge, ""), traceId)

	clientFirstMessage, errClientFirst := s.getRequest(reader, c.RemoteAddr().String())
	if errClientFirst != nil {
		return "", "", errClientFirst
	}

	auth, errParseClientFirst := s.authManager.ParseClientFirstMessageToAuthString(clientFirstMessage.msg)
	if errParseClientFirst != nil {
		return "", "", errParseClientFirst
	}
	return clientFirstMessage.user, auth, nil
}

func (s *TcpServer) prepareServerFirstMessage(auth, user string) (string, error) {
	auth, errServerFirst := s.authManager.AddServerFirstMessage(auth, user)
	if errServerFirst != nil {
		return "", errServerFirst
	}
	return auth, nil
}

func (s *TcpServer) getProofMessage(c net.Conn, writer bufio.Writer, reader bufio.Reader, traceId, auth string) (request, error) {
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
	fail             status = 0
	success          status = 1
	proofVerified    status = 7
	proofNotVerified status = 8
	authChallenge    status = 9
)

func (s *TcpServer) newResponse(status status, msg string) string {
	return fmt.Sprintf("%s%d%s\n", s.protocolVersion, status, msg)
}
