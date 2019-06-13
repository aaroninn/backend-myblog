package chat

import (
	"github.com/gorilla/websocket"
	"hypermedlab/backend-myblog/services/chat"

	"sync"
)

type Manager struct {
	rw    sync.RWMutex
	rooms map[string]*chat.Room
	conns map[string]*websocket.Conn
}

func NewManager() {
	return new(Manager)
}

func (m *Manager) AddRoom(room *chat.Room) {
	m.rw.Lock()
	defer m.rw.Unlock()
	m.rooms[room.ID] = room
}

func (m *Manager) AddConn(userid string) {

}
