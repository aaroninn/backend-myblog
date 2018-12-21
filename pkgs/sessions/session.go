package sessions

import (
	"sync"
	"time"
)

//Session struct
type Session struct {
	sessionID string
	values    map[string]interface{}
	maxAge    time.Duration
}

type SessionsStorage struct {
	Sessions sync.Map
}

const defaultAge = 24 * 60 * 60

func NewSessionsStorage() *SessionsStorage {
	var session SessionsStorage
	return &session
}

//New create a new session struct
//default age 1 day
func New() *Session {
	var s Session
	s.maxAge = defaultAge * time.Second
	return &s
}

func (s *SessionsStorage) Get(sessionid string) (*Session, bool) {
	value, ok := s.Sessions.Load(sessionid)
	if !ok {
		return nil, false
	}

	return value.(*Session), ok
}

func (s *SessionsStorage) Delete(session *Session) {
	s.Sessions.Delete(session.SessionID)
}

func (s *SessionsStorage) Save(session *Session) {
	s.Sessions.Store(session.sessionID, s)
	go time.AfterFunc(session.maxAge, func() {
		s.Sessions.Delete(session.SessionID)
	})
}

func (s *Session) SetSessionID(id string) {
	s.sessionID = id
}

func (s *Session) SexMaxage(age int) {
	s.maxAge = time.Duration(age) * time.Second
}

func (s *Session) StoreValue(key string, value interface{}) {
	s.values[key] = value
}

func (s *Session) SessionID() string {
	return s.sessionID
}

func (s *Session) MaxAge() time.Duration {
	return s.maxAge
}

func (s *Session) Value(key string) interface{} {
	return s.values[key]
}
