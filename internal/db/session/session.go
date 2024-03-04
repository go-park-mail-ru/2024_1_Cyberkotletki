package session

import (
	"github.com/google/uuid"
	"sync"
)

type sessionsDB struct {
	sync.RWMutex
	sessions map[string]int
}

func (SDB *sessionsDB) NewSession(id int) string {
	sessionId := uuid.New().String()
	SDB.Lock()
	SDB.sessions[sessionId] = id
	SDB.Unlock()
	return sessionId
}

var SessionsDB = sessionsDB{sessions: make(map[string]int)}
