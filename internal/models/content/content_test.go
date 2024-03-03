package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/person"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/models/small_models"
	"testing"
	"time"
)

func TestContentMethods(t *testing.T) {
	var c Content
	content := c.NewContentFull(
		1,
		"Test Title",
		[]small_models.Country{{Id: 1, Name: "Test Country"}},
		[]small_models.Genre{{Id: 1, Name: "Test Genre"}},
		[]person.Person{{Id: 1, FirstName: "Test Director"}},
		[]person.Person{{Id: 1, FirstName: "Test Writer"}},
		[]person.Person{{Id: 1, FirstName: "Test Producer"}},
		[]person.Person{{Id: 1, FirstName: "Test Cinematographer"}},
		"Test Slogan",
		[]person.Person{{Id: 1, FirstName: "Test Composer"}},
		[]person.Person{{Id: 1, FirstName: "Test Artist"}},
		[]person.Person{{Id: 1, FirstName: "Test Editor"}},
		123,
		321,
		[]small_models.BoxOffice{{Id: 1, Revenue: 33}},
		[]small_models.Audience{{Id: 1, AudienceT: 4}},
		time.Now(),
		time.Now(),
		18,
		8.5,
		[]person.Person{{Id: 1, FirstName: "Test Actor"}},
		[]person.Person{{Id: 1, FirstName: "Test Dubbing"}},
		[]small_models.Award{{Id: 1, AwardType: "Test Award"}},
		"Test Description",
		"Test Poster",
		"Test Playback",
	)

	if content.GetID() != 1 {
		t.Errorf("Expected Id to be 1, got %d", content.GetID())
	}

	if content.GetTitle() != "Test Title" {
		t.Errorf("Expected title to be 'Test Title', got '%s'", content.GetTitle())
	}

	if len(content.GetCountry()) != 1 || content.GetCountry()[0].Name != "Test Country" {
		t.Errorf("Expected country to be 'Test Country', got '%v'", content.GetCountry())
	}

	if len(content.GetGenres()) != 1 || content.GetGenres()[0].Name != "Test Genre" {
		t.Errorf("Expected genre to be 'Test Genre', got '%v'", content.GetGenres())
	}

	if len(content.GetDirectors()) != 1 || content.GetDirectors()[0].FirstName != "Test Director" {
		t.Errorf("Expected director to be 'Test Director', got '%v'", content.GetDirectors())
	}

	newCountry := small_models.Country{Id: 2, Name: "New Country"}
	content.AddCountry(newCountry)
	if len(content.GetCountry()) != 2 || content.GetCountry()[1].Name != "New Country" {
		t.Errorf("Expected country to be 'New Country', got '%v'", content.GetCountry())
	}
}
