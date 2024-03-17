package tmpDB

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/google/uuid"
	"sync"
)

type sessionsDB struct {
	sync.RWMutex
	Sessions map[string]int
}

func NewSessionRepository() repository.Session {
	return &sessionsDB{
		Sessions: make(map[string]int),
	}
}

func (SDB *sessionsDB) NewSession(id int) string {
	sessionId := uuid.New().String()
	SDB.Lock()
	SDB.Sessions[sessionId] = id
	SDB.Unlock()
	return sessionId
}

func (SDB *sessionsDB) CheckSession(session string) bool {
	SDB.Lock()
	defer SDB.Unlock()
	if _, ok := SDB.Sessions[session]; ok {
		return true
	}
	return false
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
