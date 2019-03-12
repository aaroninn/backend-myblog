package sessions

import (
	"log"
	"sync"
	"time"
)

//Session struct
type Session struct {
	sessionID  string
	expireTime time.Time
	createAt   time.Time
}

//SessionsStorageInMemory is the stuct of sessions storage in memory
type SessionsStorageInMemory struct {
	sessions map[string]*Session
	rw       sync.RWMutex
	age      int
	refresh  bool
}

const defaultAge = 24 * 60 * 60

func NewSessionsStorage() *SessionsStorageInMemory {
	var session SessionsStorageInMemory
	go session.checkSessionInStorage(60 * 60 * time.Second)
	return &session
}

//New create a new session struct
//default age 1 day
func NewSession() *Session {
	var s Session
	s.expireTime = time.Now().Add(defaultAge * time.Second)
	return &s
}

func (s *SessionsStorageInMemory) Get(sessionid string) (*Session, bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	session := s.sessions[sessionid]
	if time.Now().After(session.expireTime) {
		go s.Delete(sessionid)
		return nil, false
	}

	if s.refresh == true {
		go s.RefeshSession(sessionid)
	}
	return session, true
}
func (s *SessionsStorageInMemory) Delete(sessionid string) {
	s.rw.Lock()
	defer s.rw.Unlock()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic !", err)
		}
	}()
	delete(s.sessions, sessionid)
}

func (s *SessionsStorageInMemory) checkSessionInStorage(t time.Duration) {
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic !", err)
			go s.checkSessionInStorage(t)
		}
	}()
	for {
		time.Sleep(t)
		s.checkSessions()
	}
}

func (s *SessionsStorageInMemory) checkSessions() {
	s.rw.RLock()
	defer s.rw.RUnlock()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic !", err)
		}
	}()

	for k, v := range s.sessions {
		if time.Now().After(v.expireTime) {
			go s.Delete(k)
		}
	}
}

//RefeshSession refesh expiretime of id's session
func (s *SessionsStorageInMemory) RefeshSession(sessionid string) {
	s.rw.Lock()
	defer s.rw.Unlock()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic !", err)
		}
	}()
	s.sessions[sessionid].expireTime = time.Now().Add(defaultAge * time.Second)
}
