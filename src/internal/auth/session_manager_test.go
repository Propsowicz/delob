package auth

import (
	"delob/internal/utils"
	"testing"
	"time"
)

const user string = "testUser"
const ip string = "3.4.56.78"

func Test_CanAddNewSession(t *testing.T) {
	sessionManager := newSessionManager()

	result := sessionManager.AddToSession(user, ip)

	if result != nil {
		t.Errorf("Expected no error.")
	}

	session, _, err := sessionManager.getSession(user, ip)

	if err != nil {
		t.Errorf("Expected no error.")
	}
	if session.user != user {
		t.Errorf("Wrong user.")
	}
	if session.ip != ip {
		t.Errorf("Wrong ip.")
	}
}

func Test_AddingSessionWhenSessionAlreadyExistsShouldUpdateExpirationTime(t *testing.T) {
	sessionManager := newSessionManager()

	sessionManager.AddToSession(user, ip)
	session_firstLoad, _, _ := sessionManager.getSession(user, ip)

	time.Sleep(300 * time.Millisecond)

	sessionManager.AddToSession(user, ip)
	session_secondLoad, _, _ := sessionManager.getSession(user, ip)

	if session_firstLoad.expirationTime >= session_secondLoad.expirationTime {
		t.Errorf("ExpirationTime should has been updated.")
	}
}

func Test_IfSessionIsValidForExistingUser(t *testing.T) {
	sessionManager := newSessionManager()

	sessionManager.AddToSession(user, ip)

	isSessionValid := sessionManager.IsSessionValid(user, ip)

	if !isSessionValid {
		t.Errorf("Session should be valid.")
	}
}

func Test_IfSessionIsNotValidForExistingUser(t *testing.T) {
	sessionManager := newSessionManager()

	isSessionValid := sessionManager.IsSessionValid(user, ip)

	if isSessionValid {
		t.Errorf("Session should not be valid.")
	}
}

func Test_IfSessionIsNotValidForExpirationTime(t *testing.T) {
	sessionManager := newSessionManager()

	sessionManager.AddToSession(user, ip)
	sessionManager.sessions[0].expirationTime = utils.Timestamp()

	isSessionValid := sessionManager.IsSessionValid(user, ip)

	if isSessionValid {
		t.Errorf("Session should not be valid.")
	}
}
