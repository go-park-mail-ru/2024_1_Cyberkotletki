package user

import (
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/exceptions"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"sync"
	"sync/atomic"
	"time"
)

type UsersDB struct {
	DB          map[int]user.User
	dbMutex     sync.RWMutex
	usersLastId int64
}

var UsersDatabase = new(UsersDB)

// InitDB Инициализирует небольшую таблицу пользователей
func (u *UsersDB) InitDB() {
	u.dbMutex.Lock()
	defer u.dbMutex.Unlock()

	atomic.AddInt64(&u.usersLastId, 10)
	u.DB = make(map[int]user.User)

	// Заполнение базы данных DB
	u.DB[1] = user.User{
		Id:               1,
		Name:             "Egor",
		Email:            "egor@example.com",
		PasswordHash:     "hashed_password1",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}
	u.DB[2] = user.User{
		Id:               2,
		Name:             "Sasha",
		Email:            "sasha@example.com",
		PasswordHash:     "hashed_password2",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}

	u.DB[3] = user.User{
		Id:               3,
		Name:             "Kristina",
		Email:            "kristina@example.com",
		PasswordHash:     "hashed_password3",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}
}

func (u *UsersDB) HasUser(user user.User) bool {
	for _, c := range u.DB {
		if user.Email == c.Email {
			return true
		}
	}
	return false
}

func (u *UsersDB) AddUser(userObj user.User) (*user.User, *exc.Exception) {
	u.dbMutex.Lock()
	defer u.dbMutex.Unlock()
	if u.HasUser(userObj) {
		return nil, &exc.Exception{
			When:  time.Now(),
			What:  "Пользователь с таким Email уже существует",
			Layer: exc.Database,
			Type:  exc.AlreadyExists,
		}
	}
	atomic.AddInt64(&u.usersLastId, u.usersLastId+1)
	userObj.Id = int(u.usersLastId)
	u.DB[userObj.Id] = userObj
	return &userObj, nil
}

func (u *UsersDB) GetUserByEmail(email string) (*user.User, *exc.Exception) {
	u.dbMutex.Lock()
	defer u.dbMutex.Unlock()

	for _, us := range u.DB {
		if us.Email == email {
			return &us, nil
		}
	}
	return nil, &exc.Exception{
		When:  time.Now(),
		What:  "Пользователь с таким email не найден",
		Layer: exc.Database,
		Type:  exc.NotFound,
	}
}
