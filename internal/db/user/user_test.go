package user

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/content"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/user"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
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
		wantErr error
	}{
		{
			name: "Successful add user",
			newUser: user.User{
				Id:               2,
				Name:             "Новое имя",
				Email:            "new_email@example.com", // этого в бд нет
				PasswordHash:     "new_hashed_password",
				BirthDate:        time.Now(),
				SavedFilms:       []content.Film{},
				SavedSeries:      []content.Series{},
				SavedPersons:     []person.Person{},
				Friends:          []user.User{},
				ExpectedFilms:    []content.Film{},
				RegistrationDate: time.Now(),
			},
			wantErr: nil,
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
			wantErr: exc.AlreadyExistsErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usersDB.AddUser(tt.newUser)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("AddUser() error = %v, wantErr %v", err, tt.wantErr)
			}
			// проверка, что полученный тип ошибки = ожидаемому
			if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("AddUser() error = %v, wantErr %v", err, exc.AlreadyExistsErr)
			}
		})
	}
}

func TestUsersDB_GetUserByEmail(t *testing.T) {
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
		wantErr error
	}{
		{
			name:    "Successful get user by email",
			email:   "kristina@example.com",
			wantErr: nil,
		},
		{
			name:    "Unsuccessful get user by email - email does not exist",
			email:   "nonexistent@example.com",
			wantErr: exc.NotFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usersDB.GetUserByEmail(tt.email)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("GetUserByEmail() error = %v, wantErr %v", err, exc.NotFoundErr)
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

	tests := []struct {
		name string
		user user.User
		want bool
	}{
		{
			name: "User exists",
			user: notNewUser,
			want: true,
		},
		{
			name: "User does not exist",
			user: newUser,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := usersDB.HasUser(tt.user); got != tt.want {
				t.Errorf("HasUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
