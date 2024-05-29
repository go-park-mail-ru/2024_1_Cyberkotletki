package redis

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

const userSessionsPlaceholder = "user_sessions:"

type sessionsDB struct {
	rdb              *redis.Client
	sessionAliveTime int
	ctx              context.Context
}

func NewSessionRepository(rdb *redis.Client, sessionAliveTime int) repository.Session {
	db := &sessionsDB{
		rdb:              rdb,
		sessionAliveTime: sessionAliveTime,
		ctx:              context.Background(),
	}
	return db
}

func (SDB *sessionsDB) NewSession(id int) (string, error) {
	sessionID := uuid.NewString()
	// сначала нужно убедиться, что сессии с таким ключом нет
	for {
		_, err := SDB.rdb.Get(SDB.ctx, strconv.Itoa(id)).Result()
		if errors.Is(err, redis.Nil) {
			break
		}
		sessionID = uuid.NewString()
	}
	err := SDB.rdb.Set(SDB.ctx, sessionID, id, time.Duration(SDB.sessionAliveTime)*time.Second).Err()
	if err != nil {
		return "", entity.RedisWrap(errors.New("не удалось создать сессию"), err)
	}
	// Добавляем сессию в список сессий пользователя
	err = SDB.rdb.SAdd(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(id), sessionID).Err()
	if err != nil {
		return "", entity.RedisWrap(errors.New("не удалось создать сессию"), err)
	}
	return sessionID, nil
}

// CheckSession проверяет сессию и возвращает id пользователя
func (SDB *sessionsDB) CheckSession(session string) (int, error) {
	userID, err := SDB.rdb.Get(SDB.ctx, session).Result()
	if errors.Is(err, redis.Nil) {
		return 0, repository.ErrSessionNotFound
	}
	if err != nil {
		return 0, entity.RedisWrap(errors.New("не удалось проверить сессию"), err)
	}
	return strconv.Atoi(userID)
}

func (SDB *sessionsDB) DeleteAllSessions(userID int) error {
	// Получаем список всех сессий пользователя
	sessionIDs, err := SDB.rdb.SMembers(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(userID)).Result()
	if err != nil {
		return entity.RedisWrap(errors.New("не удалось получить список сессий пользователя"), err)
	}
	// Удаляем каждую сессию
	for _, sessionID := range sessionIDs {
		err = SDB.rdb.Del(SDB.ctx, sessionID).Err()
		if err != nil {
			return entity.RedisWrap(errors.New("не удалось удалить сессию"), err)
		}
	}
	// Удаляем список сессий пользователя
	err = SDB.rdb.Del(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(userID)).Err()
	if err != nil {
		return entity.RedisWrap(errors.New("не удалось удалить список сессий пользователя"), err)
	}
	return nil
}

func (SDB *sessionsDB) DeleteSession(session string) error {
	id, err := SDB.rdb.Get(SDB.ctx, session).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return entity.RedisWrap(errors.New("не удалось удалить сессию"), err)
	}
	err = SDB.rdb.Del(SDB.ctx, session).Err()
	if err != nil {
		return entity.RedisWrap(errors.New("не удалось удалить сессию"), err)
	}
	// Удаляем сессию из списка сессий пользователя
	err = SDB.rdb.SRem(SDB.ctx, userSessionsPlaceholder+id, session).Err()
	if err != nil {
		return entity.RedisWrap(errors.New("не удалось удалить сессию"), err)
	}
	return nil
}
