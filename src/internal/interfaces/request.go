package interfaces

import "strings"

type request struct {
	user string
	msg  string
	ip   string
}

func parseRequest(s, ip string) (request, error) {
	parts := strings.Split(s, "|||")

	r := request{}
	r.user = parts[0]
	r.msg = parts[1]
	r.ip = ip
	return r, nil
}
