package auth

import (
	"testing"
	"time"
)

func Test_IfTryAuthenticateReturnsFalseWhenEmpty(t *testing.T) {
	authM := NewAuthenticationManager()
	username := "testUser"
	ip := "3.4.56.78"

	result := authM.TryAuthenticate(username, ip)

	if result {
		t.Errorf("Expected false result.")
	}
}

func Test_IfTryAuthenticateReturnsTrueWhenFindNotExpiredSession(t *testing.T) {
	authM := NewAuthenticationManager()
	username := "testUser"
	ip := "3.4.56.78"
	authM.sessions = append(authM.sessions, session{
		user:           username,
		ip:             ip,
		expirationTime: time.Now().UnixMilli() + int64(time.Minute),
		hash:           "s53ef23",
	})

	result := authM.TryAuthenticate(username, ip)

	if !result {
		t.Errorf("Expected true result.")
	}
}

func Test_IfTryAuthenticateReturnsFalseWhenCannotFindNotExpiredSession(t *testing.T) {
	authM := NewAuthenticationManager()
	username := "testUser1"
	ip := "3.4.56.78232"
	authM.sessions = append(authM.sessions, session{
		user:           "a",
		ip:             "0.0.0.0",
		expirationTime: time.Now().UnixMilli(),
		hash:           "s53ef23",
	})

	result := authM.TryAuthenticate(username, ip)

	if result {
		t.Errorf("Expected false result.")
	}
}
