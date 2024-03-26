package redis

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/config"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type sessionsDB struct {
	rdb              *redis.Client
	sessionAliveTime int
	ctx              context.Context
}

func NewSessionRepository(logger echo.Logger, config config.Config) repository.Session {
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
		logger.Fatal("Не удалось подключиться к Redis", err.Error())
	}
	logger.Info("Redis с сессиями успешно запущен")
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
		return "", entity.NewClientError("не удалось создать сессию", err, entity.ErrRedis)
	}
	return sessionID, nil
}

func (SDB *sessionsDB) CheckSession(session string) (bool, error) {
	_, err := SDB.rdb.Get(SDB.ctx, session).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		return false, entity.NewClientError("не удалось проверить сессию", err, entity.ErrRedis)
	}
	return true, nil
}

func (SDB *sessionsDB) DeleteSession(session string) (bool, error) {
	err := SDB.rdb.Del(SDB.ctx, session).Err()
	if err != nil {
		return false, entity.NewClientError("не удалось удалить сессию", err, entity.ErrRedis)
	}
	return true, nil
}
