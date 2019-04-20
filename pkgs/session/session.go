package session

import (
	"log"
	"sync"
	"time"
)

//Session struct
type Session struct {
	sessionID  string
	data       interface{}
	expireTime time.Time
	createAt   time.Time
}

//NewSession return a ptr of Session struct with default expiretime
//Default expiretime is 1 day
//Use *Session.SetExpireTime to modify expiretime
func NewSession(sessionid string) *Session {
	session := new(Session)
	session.sessionID = sessionid
	session.expireTime = time.Now().Add(defaultAge * time.Second)
	return session
}

//SetExpireTime set the expire time of session
func (s *Session) SetExpireTime(t int) {
	s.expireTime = time.Now().Add(time.Duration(t) * time.Second)
}

//SetData add data in session
func (s *Session) SetData(data interface{}) {
	s.data = data
}

//GetData return data storage in session
func (s *Session) GetData() interface{} {
	return s.data
}

//SessionsStorageInMemory is the struct of sessions storage in memory
type SessionsStorageInMemory struct {
	sessions map[string]*Session
	rw       sync.RWMutex
	age      int
	refresh  bool
}

const defaultAge = 24 * 60 * 60

//NewSessionsStorage return a ptr of SessionsStorageInMemory struct
func NewSessionsStorage() *SessionsStorageInMemory {
	session := new(SessionsStorageInMemory)
	//check every hour to make sure expireout session is deleted
	go session.checkSessionInStorage(60 * 60 * time.Second)
	return session
}

//Add add a session in storage
//NewSession return a ptr of Session Struct
func (s *SessionsStorageInMemory) Add(session *Session) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.sessions[session.sessionID] = session
}

//Get find session in storage by sessionid
//If session not exist it will return nil, false
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

//Delete delete session in storage by sessionid
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

//Update update session's data
func (s *SessionsStorageInMemory) Update(sessionid string, data interface{}) error {
	s.rw.Lock()
	defer s.rw.Unlock()
	_, ok := s.sessions[sessionid]
	if !ok {
		return errSessionNotExist
	}

	s.sessions[sessionid].data = data
	return nil
}

//SessionAmount return the amount of sessions storage in memory
func (s *SessionsStorageInMemory) SessionAmount() int {
	s.rw.RLock()
	defer s.rw.RUnlock()
	return len(s.sessions)
}

func (s *SessionsStorageInMemory) deleteSessions(sessionids []string) {
	s.rw.Lock()
	defer s.rw.Unlock()
	defer func() {
		err := recover()
		if err != nil {
			log.Println("panic !", err)
		}
	}()
	for _, v := range sessionids {
		delete(s.sessions, v)
	}
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

//SetAutoFresh auto fresh when Get
func (s *SessionsStorageInMemory) SetAutoFresh() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.refresh = true
}

func (s *SessionsStorageInMemory) DisableAutoFresh() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.refresh = false
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

	sessionids := make([]string, 0)
	for _, v := range s.sessions {
		if time.Now().After(v.expireTime) {
			sessionids = append(sessionids, v.sessionID)
		}
	}

	s.deleteSessions(sessionids)
}

//RefeshSession refesh expiretime of id's session
//If session not exist, it will return a error
func (s *SessionsStorageInMemory) RefeshSession(sessionid string) error {
	s.rw.Lock()
	defer s.rw.Unlock()

	_, ok := s.sessions[sessionid]
	if !ok {
		return errSessionNotExist
	}
	s.sessions[sessionid].expireTime = time.Now().Add(defaultAge * time.Second)
	return nil
}
