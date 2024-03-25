package tmpdb

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"sync"
	"sync/atomic"
)

type UsersDB struct {
	sync.RWMutex
	DB          map[int]entity.User
	usersLastID atomic.Int64
}

func NewUserRepository() repository.User {
	return &UsersDB{
		DB: make(map[int]entity.User),
	}
}

func (u *UsersDB) HasUser(user *entity.User) bool {
	for _, c := range u.DB {
		if user.Email == c.Email {
			return true
		}
	}
	return false
}

func (u *UsersDB) AddUser(user *entity.User) (*entity.User, error) {
	u.Lock()
	defer u.Unlock()
	if u.HasUser(user) {
		return nil, entity.NewClientError("пользователь с таким email уже существует", entity.ErrAlreadyExists)
	}
	u.usersLastID.CompareAndSwap(u.usersLastID.Load(), u.usersLastID.Load()+1)
	user.ID = int(u.usersLastID.Load())
	u.DB[user.ID] = *user
	return user, nil
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
