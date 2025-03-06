package auth

import (
	"delob/internal/utils"
	"fmt"
	"math/rand/v2"
)

type AuthenticationManager struct {

	// handshake client <-> data base SCRAM (password)
	// tcp client <-> data base
	//
	// 1. Handshake
	// 2. userData persistance
	// 3. session storage

	sessions []session
}

type session struct {
	ip             string
	expirationTime int64
	user           string
	hash           string
	client_first   string
	server_first   string
	client_final   string
}

type userData struct {
	user            string
	salt            string
	hashed_pwd      string
	iteration_count int
}

func NewAuthenticationManager() AuthenticationManager {
	return AuthenticationManager{
		sessions: []session{},
	}
}

// 0-1 failure/success
// server responses:
// 9 - not auth -> challenge me!
// 8 - here is challenge data
// 6/7 forbidden/access
// client --> server
// server -> data (salt, random number) -> client
// client -> generate hash -> server
// server generate hash -> compare hashes
// ok -> 7 & save session
// not OK -> forbid - 6

func (a *AuthenticationManager) TryAuthenticate(user, ip string) bool {
	if len(a.sessions) == 0 {
		return false
	}

	for i := range a.sessions {
		if a.sessions[i].user == user && a.sessions[i].ip == ip && utils.Timestamp() <= a.sessions[i].expirationTime {
			return true
		}
	}

	return false
}

func (a *AuthenticationManager) Ch(user, ip string, clientNonce int) string {
	// load user's hash and password

	userData := loadUserData(user)

	serverNonce := rand.Int()
	nonce := clientNonce + serverNonce

	serverResponse := fmt.Sprintf("nonce=%d,salt=%s,iterations=%d", nonce, userData.salt, userData.iteration_count)
	a.sessions = append(a.sessions, candidateSession(user, ip, string(clientNonce), serverResponse))

	return serverResponse
}

func loadUserData(user string) userData {
	return userData{
		user:            "test",
		salt:            "qjhklufhvkduyihetr",
		hashed_pwd:      "qweqwdsfcscvrg",
		iteration_count: 4,
	}
}

func candidateSession(user, ip, client_first, server_first string) session {
	return session{
		user:           user,
		ip:             ip,
		expirationTime: 0,
		client_first:   client_first,
		server_first:   server_first,
	}
}

func remove(s []session, i int) []session {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
