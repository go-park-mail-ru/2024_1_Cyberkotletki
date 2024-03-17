package tmpDB

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"sync"
	"sync/atomic"
)

type UsersDB struct {
	sync.RWMutex
	DB          map[int]entity.User
	usersLastId atomic.Int64
}

func NewUserRepository() repository.User {
	return &UsersDB{
		DB: make(map[int]entity.User),
	}
}

func (u *UsersDB) HasUser(user entity.User) bool {
	for _, c := range u.DB {
		if user.Email == c.Email {
			return true
		}
	}
	return false
}

func (u *UsersDB) AddUser(user entity.User) (*entity.User, error) {
	u.Lock()
	defer u.Unlock()
	if u.HasUser(user) {
		return nil, entity.NewClientError("пользователь с таким email уже существует", entity.ErrBadRequest)
	}
	u.usersLastId.CompareAndSwap(u.usersLastId.Load(), u.usersLastId.Load()+1)
	user.Id = int(u.usersLastId.Load())
	u.DB[user.Id] = user
	return &user, nil
}

func (u *UsersDB) GetUserByEmail(email string) (*entity.User, error) {
	u.Lock()
	defer u.Unlock()

	for _, us := range u.DB {
		if us.Email == email {
			return &us, nil
		}
	}
	return nil, entity.NewClientError("пользователь с таким email не найден", entity.ErrNotFound)
}
