package db

import (
	UsersDB "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/user"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"time"

	//"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"testing"
)

func TestInitU(t *testing.T) {
	var usersDB UsersDB.UsersDB
	usersDB.InitDB()

	if len(usersDB.DB) != 3 {
		t.Errorf("Expected length of UsersDB to be 3, got %d", len(usersDB.DB))
	}

	if user_Obj, ok := usersDB.DB[1]; !ok || user_Obj.Name != "Egor" {
		t.Errorf("Expected to find user_Obj with name Egor, got %v", user_Obj)
	}
}

func TestAddUser(t *testing.T) {
	var usersDB UsersDB.UsersDB
	usersDB.InitDB()

	newUser := user.User{
		Id:               2,
		Name:             "Новое имя",
		Email:            "new_email@example.com",
		PasswordHash:     "new_hashed_password",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}

	userObj, err := usersDB.AddUser(newUser)

	if err != nil {
		t.Errorf("TestAddUser failed, unexpected error: %v", err)
	}

	if userObj == nil {
		t.Errorf("TestAddUser failed, userObj is nil")
	}

	if userObj.Email != newUser.Email {
		t.Errorf("111TestAddUser failed, userObj does not match newUser")
	}
	if userObj.Id == 4 {
		t.Errorf("TestAddUser failed, userObj does not match newUser")
	}

	if _, ok := usersDB.DB[newUser.Id]; !ok {
		t.Errorf("TestAddUser failed, user not added to DB")
	}
}

func TestGetUserByEmail(t *testing.T) {
	var usersDB UsersDB.UsersDB
	usersDB.InitDB()

	// Добавьте пользователя в базу данных перед тестированием функции GetUserByEmail
	testUser := user.User{
		Id:    2,
		Email: "kristina@example.com",
		// Заполните остальные поля, если это необходимо
	}
	usersDB.DB[testUser.Id] = testUser

	user_Obj, err := usersDB.GetUserByEmail("kristina@example.com")

	if err != nil {
		t.Errorf("TestGetUserByEmail failed, unexpected error: %v", err)
	}

	if user_Obj == nil {
		t.Errorf("TestGetUserByEmail failed, user_Obj is nil")
	}

	if user_Obj.Id != testUser.Id || user_Obj.Email != testUser.Email {
		t.Errorf("TestGetUserByEmail failed, user_Obj does not match testUser")
	}
}

func TestHasUser(t *testing.T) {
	var usersDB UsersDB.UsersDB
	usersDB.InitDB()

	newUser := user.User{
		Id:               2,
		Name:             "Новое имя",
		Email:            "new_email@example.com",
		PasswordHash:     "new_hashed_password",
		BirthDate:        time.Now(),
		SavedFilms:       []content.Film{},
		SavedSeries:      []content.Series{},
		SavedPersons:     []person.Person{},
		Friends:          []user.User{},
		ExpectedFilms:    []content.Film{},
		RegistrationDate: time.Now(),
	}
	notNewUser := user.User{
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

	okFalse := usersDB.HasUser(newUser)
	okTrue := usersDB.HasUser(notNewUser)

	if !okTrue && okFalse {
		t.Errorf("TestHasUser failed, expected okTrue and okFalse, got %v and %v", okTrue, okFalse)
	}
}
