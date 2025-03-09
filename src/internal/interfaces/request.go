package interfaces

import "strings"

type request struct {
	user string
	msg  string
	ip   string
}

func parseRequest(s, ip string) (request, error) {
	// CRLF (Carriage Return + Line Feed) as separator
	parts := strings.Split(s, "\r\n")

	r := request{}
	r.user = parts[0]
	r.msg = parts[1]
	r.ip = ip
	return r, nil
}
