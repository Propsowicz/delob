package interfaces

import "strings"

type request struct {
	user string
	msg  string
	ip   string
}

func parseRequest(s, ip string) (request, error) {
	const uniqueDelimiter string = "\x1E\x1F"
	parts := strings.Split(s, uniqueDelimiter)

	r := request{}
	r.user = parts[0]
	r.msg = parts[1]
	r.ip = ip
	return r, nil
}
