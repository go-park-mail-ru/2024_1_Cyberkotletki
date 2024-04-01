package redis

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
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

func NewSessionRepository(config config.Config) (repository.Session, error) {
	db := &sessionsDB{
		rdb: redis.NewClient(&redis.Options{
			Addr:     config.Auth.Redis.Addr,
			Password: config.Auth.Redis.Password,
			DB:       config.Auth.Redis.DB,
		}),
		sessionAliveTime: config.Auth.SessionAliveTime,
		ctx:              context.Background(),
	}
	if err := db.rdb.Ping(db.ctx).Err(); err != nil {
		return nil, err
	}
	return db, nil
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
		return "", entity.NewClientError("не удалось создать сессию", err, entity.ErrRedis)
	}
	// Добавляем сессию в список сессий пользователя
	err = SDB.rdb.SAdd(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(id), sessionID).Err()
	if err != nil {
		return "", entity.NewClientError("не удалось создать сессию", err, entity.ErrRedis)
	}
	return sessionID, nil
}

// CheckSession проверяет сессию и возвращает id пользователя
func (SDB *sessionsDB) CheckSession(session string) (int, error) {
	userID, err := SDB.rdb.Get(SDB.ctx, session).Result()
	if errors.Is(err, redis.Nil) {
		return 0, entity.NewClientError("пользователь с такой сессией не найден", err, entity.ErrNotFound)
	}
	if err != nil {
		return 0, entity.NewClientError("не удалось проверить сессию", err, entity.ErrRedis)
	}
	return strconv.Atoi(userID)
}

func (SDB *sessionsDB) DeleteAllSessions(userID int) error {
	// Получаем список всех сессий пользователя
	sessionIDs, err := SDB.rdb.SMembers(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(userID)).Result()
	if err != nil {
		return entity.NewClientError("не удалось получить список сессий пользователя", err, entity.ErrRedis)
	}
	// Удаляем каждую сессию
	for _, sessionID := range sessionIDs {
		err = SDB.rdb.Del(SDB.ctx, sessionID).Err()
		if err != nil {
			return entity.NewClientError("не удалось удалить сессию", err, entity.ErrRedis)
		}
	}
	// Удаляем список сессий пользователя
	err = SDB.rdb.Del(SDB.ctx, userSessionsPlaceholder+strconv.Itoa(userID)).Err()
	if err != nil {
		return entity.NewClientError("не удалось удалить список сессий пользователя", err, entity.ErrRedis)
	}
	return nil
}

func (SDB *sessionsDB) DeleteSession(session string) error {
	id, err := SDB.rdb.Get(SDB.ctx, session).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	if err != nil {
		return entity.NewClientError("не удалось удалить сессию", err, entity.ErrRedis)
	}
	err = SDB.rdb.Del(SDB.ctx, session).Err()
	if err != nil {
		return entity.NewClientError("не удалось удалить сессию", err, entity.ErrRedis)
	}
	// Удаляем сессию из списка сессий пользователя
	err = SDB.rdb.SRem(SDB.ctx, userSessionsPlaceholder+id, session).Err()
	if err != nil {
		return entity.NewClientError("не удалось удалить сессию", err, entity.ErrRedis)
	}
	return nil
}
