package auth

import (
	"delob/internal/utils"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	user       string
	salt       string
	hashed_pwd string
	iterations int
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

func (a *AuthenticationManager) AddServerFirstMessage(auth, user string) (string, error) {
	userData := loadUserData(user)

	auth = a.addToAuth(auth, fmt.Sprintf("s_nonce=%d,", generateNonce()))
	auth = a.addToAuth(auth, fmt.Sprintf("salt=%s,", userData.salt))
	auth = a.addToAuth(auth, fmt.Sprintf("iterations=%d", userData.iterations))
	return auth, nil
}

func (a *AuthenticationManager) AddClientFirstAuthString(message string) (string, error) {
	user, c_nonce, err := a.parseClientFirst(message)
	if err != nil {
		return "", err
	}

	return a.addClientFirstAuthString(user, c_nonce), nil
}

func (a *AuthenticationManager) addClientFirstAuthString(user string, c_nonce int) string {
	var auth string
	auth = a.addToAuth(auth, fmt.Sprintf("user=%s,", user))
	auth = a.addToAuth(auth, fmt.Sprintf("c_nonce=%d,", c_nonce))
	return auth
}

func (a *AuthenticationManager) addToAuth(auth string, s interface{}) string {
	var toAdd string
	switch v := s.(type) {
	case string:
		toAdd = v
	case int, int8:
		toAdd = fmt.Sprintf("%d", v)
	}

	return auth + toAdd
}

func (a *AuthenticationManager) parseClientFirst(s string) (string, int, error) {
	parts := strings.Split(s, ",")
	fmt.Println(s)
	const userPrefix string = "user="
	const cnoncePrefix string = "c_nonce="
	if parts[0][0:len(userPrefix)] == userPrefix && parts[1][0:len(cnoncePrefix)] == cnoncePrefix {
		if nonce, err := strconv.Atoi(parts[1][len(cnoncePrefix):]); err == nil {
			return parts[0][len(userPrefix):], nonce, nil
		}
	}

	return "", 1, fmt.Errorf("cannot parse client first message")
}

func (a *AuthenticationManager) Ch(user, ip string, clientNonce int) string {
	// load user's hash and password

	userData := loadUserData(user)

	serverNonce := rand.Int()
	nonce := clientNonce + serverNonce

	serverResponse := fmt.Sprint("nonce=%d,salt=%s,iterations=%d", nonce, userData.salt, userData.iterations)
	a.sessions = append(a.sessions, candidateSession(user, ip, string(clientNonce), serverResponse))

	return serverResponse
}

func loadUserData(user string) userData {
	return userData{
		user:       "test",
		salt:       "qjhklufhvkduyihetr",
		hashed_pwd: "qweqwdsfcscvrg",
		iterations: 4,
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

func generateNonce() int {
	return rand.Intn(256)
}
