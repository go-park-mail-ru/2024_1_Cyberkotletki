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

func (SDB *sessionsDB) CheckSession(session string) (int, bool) {
	SDB.Lock()
	defer SDB.Unlock()
	if s, ok := SDB.sessions[session]; ok {
		return s, true
	}
	return -1, false
}

func (SDB *sessionsDB) DeleteSession(session string) bool {
	SDB.Lock()
	defer SDB.Unlock()
	if _, ok := SDB.sessions[session]; ok {
		delete(SDB.sessions, session)
		return true
	}
	return false
}

var SessionsDB = sessionsDB{sessions: make(map[string]int)}
