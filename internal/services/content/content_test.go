package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"

	"testing"
)

func TestPreviewInfoData_GetContentPreviewInfo(t *testing.T) {
	// Инициализация базы данных
	content.FilmsDatabase.InitDB()
	tests := []struct {
		name      string
		contentId int
		wantErr   bool
	}{
		{
			name:      "Successful get content preview info",
			contentId: 1,
			wantErr:   false,
		},
		{
			name:      "Unsuccessful get content preview info - content does not exist",
			contentId: 1000,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContentPreviewInfo(tt.contentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContentPreviewInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
