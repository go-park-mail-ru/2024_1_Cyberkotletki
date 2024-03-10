package collections

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"
	"testing"
)

func TestCompilationData_GetCompilation(t *testing.T) {
	// Инициализация базы данных
	content.FilmsDatabase.InitDB()

	tests := []struct {
		name    string
		genre   string
		wantErr error
	}{
		{
			name:    "Successful get compilation",
			genre:   "drama",
			wantErr: nil,
		},
		{
			name:    "Successful get compilation",
			genre:   "action",
			wantErr: nil,
		},
		{
			name:    "Successful get compilation",
			genre:   "comedian",
			wantErr: nil,
		},
		{
			name:    "Unsuccessful get compilation - genre does not exist",
			genre:   "nonexistent",
			wantErr: exc.NotFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCompilation(tt.genre)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetCompilation() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("GetCompilation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
