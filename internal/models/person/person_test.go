package person

import (
	"testing"
	"time"
)

// ДОДЕЛАТЬ ТЕСТЫ!!!!!!!

/*func TestCreatePerson(t *testing.T) {
	birthDate := time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC)
	startCareer := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	birthPlace := PlaceOfBirth{City: "Concord", Region: "California", Country: "USA"}
	genres := []string{"Drama", "Comedy"}
	career := []string{"Actor", "Director"}

	person := CreatePerson(1, "Tom", "Hanks", birthDate, 67, nil, startCareer, nil, "https://example.com/photo.jpg", birthPlace, genres, career)

	if person.GetID() != 1 || person.GetFirstName() != "Tom" || person.GetLastName() != "Hanks" {
		t.Errorf("CreatePerson() failed, expected %v, got %v", 1, person.GetID())
	}

	if !person.BirthDate.Equal(birthDate) {
		t.Errorf("CreatePerson() failed, expected %v, got %v", birthDate, person.BirthDate)
	}

	if person.Photo != "https://example.com/photo.jpg" {
		t.Errorf("CreatePerson() failed, expected %v, got %v", "https://example.com/photo.jpg", person.Photo)
	}
}*/

func TestNewPersonFromJSON(t *testing.T) {
	// формат времени “1956-07-09T00:00:00Z” : “T” - разделитель, который отделяет дату от времени,“Z” - обозначение временной зоны
	jsonStr := `{
		"id": 1,
		"first_name": "Tom",
		"last_name": "Hanks",
		"birth_date": "1956-07-09T00:00:00Z",
		"age": 67,
		"start_career": "1979-01-01T00:00:00Z",
		"photo": "https://example.com/photo.jpg",
		"birth_place": {
			"city": "Concord",
			"region": "California",
			"country": "USA"
		},
		"genres": ["Drama", "Comedy"],
		"career": ["Actor", "Director"]
	}`

	person, err := NewPersonFromJSON(jsonStr)
	// проверка на то, была ли ошибка при декодировании
	if err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	// проверки на прочтение данных
	if person.ID != 1 || person.FirstName != "Tom" || person.LastName != "Hanks" {
		t.Errorf("NewPersonFromJSON() failed, expected %v, got %v", 1, person.ID)
	}

	if !person.BirthDate.Equal(time.Date(1956, 7, 9, 0, 0, 0, 0, time.UTC)) {
		t.Errorf("NewPersonFromJSON() failed, expected %v, got %v", time.Date(1956, 7, 9, 0, 0, 0, 0, time.UTC), person.BirthDate)
	}

	if person.Photo != "https://example.com/photo.jpg" {
		t.Errorf("NewPersonFromJSON() failed, expected %v, got %v", "https://example.com/photo.jpg", person.Photo)
	}
}
