package auth

import (
	"delob/internal/utils"
	"delob/internal/utils/logger"
	"fmt"
)

const sessionLengthInMinuntes int8 = 15

type sessionManager struct {
	sessions []session
}

type session struct {
	id             uint32
	user           string
	ip             string
	expirationTime int64
}

func newSessionManager() sessionManager {
	return sessionManager{
		sessions: []session{},
	}
}

func (s *sessionManager) AddToSession(user, ip string) error {
	_, idx, err := s.getSession(user, ip)
	if err == nil {
		s.sessions[idx].expirationTime = utils.TimestampMinutesOffset(sessionLengthInMinuntes)
		logger.Info("", fmt.Sprintf("updated session for user: %s and ip: %s", user, ip))
		return nil
	}

	s.sessions = append(s.sessions, s.createNewSession(user, ip))
	return nil
}

func (s *sessionManager) IsSessionValid(user, ip string) bool {
	session, _, err := s.getSession(user, ip)
	if err != nil {
		return false
	}

	if utils.Timestamp() < session.expirationTime {
		return true
	}
	return false
}

func (s *sessionManager) getSession(user, ip string) (session, int, error) {
	sessionId := calculateSessionId(user, ip)

	for i := range s.sessions {
		if s.sessions[i].id == sessionId {
			return s.sessions[i], i, nil
		}
	}

	return session{}, 0, fmt.Errorf("no active session exists for the user %s and ip %s.", user, ip)
}

func (s *sessionManager) createNewSession(user, ip string) session {
	logger.Info("", fmt.Sprintf("created a new session for user: %s and ip: %s", user, ip))
	return session{
		id:             calculateSessionId(user, ip),
		user:           user,
		ip:             ip,
		expirationTime: utils.TimestampMinutesOffset(sessionLengthInMinuntes),
	}
}

func calculateSessionId(user, ip string) uint32 {
	result, err := utils.Calculate(user + ip)
	if err != nil {
		panic(err.Error())
	}
	return result
}
