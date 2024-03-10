package collections

import (
	"reflect"
	"testing"
)

func TestGenresData_GetGenres(t *testing.T) {
	expectedGenres := []string{"action", "drama", "comedian"}
	genresData, err := GetGenres()

	if err != nil {
		t.Errorf("GetGenres() error = %v", err)
	}

	if !reflect.DeepEqual(genresData.Genres, expectedGenres) {
		t.Errorf("GetGenres() = %v, want %v", genresData.Genres, expectedGenres)
	}
}
