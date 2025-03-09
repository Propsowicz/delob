package interfaces

import "fmt"

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
