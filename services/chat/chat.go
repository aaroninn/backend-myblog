package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID    string
	rw    sync.RWMutex
	conns map[string]string
}

func NewRoom() *room {
	return new(room)
}

func (r *room) addConn(userid string, conn *websocket.Conn) {
	r.rw.Lock()
	defer r.rw.Unlock()

	r.conns[userid] = conn
}
