package db

import (
	"errors"
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

// InitUsersDB Инициализирует небольшую таблицу пользователей
func (u *UsersDB) InitUsersDB() {
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

func (u *UsersDB) HasUserWithEmailAndPassword(user user.User) bool {
	for _, c := range u.DB {
		if user.Email == c.Email && user.PasswordHash == c.PasswordHash {
			return true
		}
	}
	return false
}

func (u *UsersDB) AddUser(id int, name string, email string, passwordHash string, birthDate time.Time, savedFilms []content.Film,
	savedSeries []content.Series, savedPersons []person.Person, friends []user.User, expectedFilms []content.Film,
	registrationDate time.Time) (*user.User, error) {

	user_obj := *user.NewUserFull(id, name, email, passwordHash, birthDate, savedFilms,
		savedSeries, savedPersons, friends, expectedFilms, registrationDate)
	u.dbMutex.Lock()
	defer u.dbMutex.Unlock()
	if u.HasUser(user_obj) {
		err := errors.New("user_obj with this id email already exists")
		return &user_obj, err
	}
	atomic.AddInt64(&u.usersLastId, u.usersLastId+1)
	user_obj.Id = int(u.usersLastId)
	u.DB[user_obj.Id] = user_obj
	return &user_obj, nil
}

func (u *UsersDB) GetUser(id int) (*user.User, error) {
	u.dbMutex.Lock()
	defer u.dbMutex.Unlock()

	user_obj, ok := u.DB[id]
	if ok {
		return &user_obj, nil
	}
	err := errors.New("user with this id not found")
	return nil, err
}
