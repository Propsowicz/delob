package interfaces

import "fmt"

type status int8

const (
	fail                     status = 0
	success                  status = 1
	user_verified            status = 7
	user_not_verified        status = 8
	authentication_challenge status = 9
)

func (s *TcpServer) newResponse(status status, msg string) string {
	return fmt.Sprintf("%s%d%s\n", s.protocolVersion, status, msg)
}
