package session

import (
	"github.com/google/uuid"
	"sync"
)

type sessionsDB struct {
	sync.RWMutex
	Sessions map[string]int
}

func (SDB *sessionsDB) NewSession(id int) string {
	sessionId := uuid.New().String()
	SDB.Lock()
	SDB.Sessions[sessionId] = id
	SDB.Unlock()
	return sessionId
}

func (SDB *sessionsDB) CheckSession(session string) (int, bool) {
	SDB.Lock()
	defer SDB.Unlock()
	if s, ok := SDB.Sessions[session]; ok {
		return s, true
	}
	return -1, false
}

func (SDB *sessionsDB) DeleteSession(session string) bool {
	SDB.Lock()
	defer SDB.Unlock()
	if _, ok := SDB.Sessions[session]; ok {
		delete(SDB.Sessions, session)
		return true
	}
	return false
}

var SessionsDB = sessionsDB{Sessions: make(map[string]int)}
