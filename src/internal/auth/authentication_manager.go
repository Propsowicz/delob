package auth

import (
	"delob/internal/utils"
	"encoding/hex"
	"fmt"
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

	auth = addToAuth(auth, fmt.Sprintf("s_nonce=%d,", generateNonce()))
	auth = addToAuth(auth, fmt.Sprintf("salt=%s,", hex.EncodeToString(userData.Salt)))
	auth = addToAuth(auth, fmt.Sprintf("iterations=%d", userData.Iterations))
	return auth, nil
}

func (a *AuthenticationManager) ParseClientFirstMessageToAuthString(message string) (string, error) {
	user, c_nonce, err := parseClientFirst(message)
	if err != nil {
		return "", err
	}

	return addClientFirstAuthString(user, c_nonce), nil
}

func (a *AuthenticationManager) Verify(proof, user, auth string) bool {
	userData := loadUserData(user)
	clientSignature := computeHmacHash(userData.Stored_key, []byte(auth))

	serverSideProof := hex.EncodeToString(xorBytes(userData.Client_key, clientSignature))

	return serverSideProof == proof
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
