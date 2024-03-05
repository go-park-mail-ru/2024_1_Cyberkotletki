package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	"testing"
)

func TestCompilationData_GetCompilation(t *testing.T) {
	// Инициализация базы данных
	content.FilmsDatabase.InitDB()

	tests := []struct {
		name    string
		genre   string
		wantErr bool
	}{
		{
			name:    "Successful get compilation",
			genre:   "drama",
			wantErr: false,
		},
		{
			name:    "Successful get compilation",
			genre:   "action",
			wantErr: false,
		},
		{
			name:    "Successful get compilation",
			genre:   "comedian",
			wantErr: false,
		},
		{
			name:    "Unsuccessful get compilation - genre does not exist",
			genre:   "nonexistent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCompilation(tt.genre)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCompilation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
