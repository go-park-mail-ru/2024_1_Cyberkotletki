package user

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	"testing"
	"time"
)

func TestUsersDB_Init(t *testing.T) {
	var usersDB UsersDB
	usersDB.InitDB()

	if len(usersDB.DB) != 3 {
		t.Errorf("Expected length of UsersDB to be 3, got %d", len(usersDB.DB))
	}

	if userObj, ok := usersDB.DB[1]; !ok || userObj.Name != "Egor" {
		t.Errorf("Expected to find user_Obj with name Egor, got %v", userObj)
	}
}
func TestUsersDB_AddUser(t *testing.T) {
	var usersDB UsersDB
	usersDB.InitDB()

	tests := []struct {
		name    string
		newUser user.User
		wantErr bool
	}{
		{
			name: "Successful add user",
			newUser: user.User{
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
			},
			wantErr: false,
		},
		{
			name: "Unsuccessful add user - email already exists",
			newUser: user.User{
				Id:               3,
				Name:             "Другое имя",
				Email:            "new_email@example.com", // это уже есть в бд
				PasswordHash:     "another_hashed_password",
				BirthDate:        time.Now(),
				SavedFilms:       []content.Film{},
				SavedSeries:      []content.Series{},
				SavedPersons:     []person.Person{},
				Friends:          []user.User{},
				ExpectedFilms:    []content.Film{},
				RegistrationDate: time.Now(),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usersDB.AddUser(tt.newUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestTestUsersDB_GetUserByEmail(t *testing.T) {
	var usersDB UsersDB
	usersDB.InitDB()

	testUser := user.User{
		Id:    2,
		Email: "kristina@example.com",
	}
	usersDB.DB[testUser.Id] = testUser

	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "Successful get user by email",
			email:   "kristina@example.com",
			wantErr: false,
		},
		{
			name:    "Unsuccessful get user by email - email does not exist",
			email:   "nonexistent@example.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usersDB.GetUserByEmail(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsersDB_HasUser(t *testing.T) {
	var usersDB UsersDB
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
