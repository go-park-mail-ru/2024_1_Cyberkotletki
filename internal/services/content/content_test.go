package content

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/db/content"
	exc "github.com/go-park-mail-ru/2024_1_Cyberkotletki/pkg/exceptions"

	"testing"
)

func TestPreviewInfoData_GetContentPreviewInfo(t *testing.T) {
	content.FilmsDatabase.InitDB()
	tests := []struct {
		name      string
		contentId int
		wantErr   error
	}{
		{
			name:      "Successful get content preview info",
			contentId: 1,
			wantErr:   nil,
		},
		{
			name:      "Unsuccessful get content preview info - content does not exist",
			contentId: 1000,
			wantErr:   exc.NotFoundErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContentPreviewInfo(tt.contentId)
			if (err != nil) != (tt.wantErr != nil) {
				t.Errorf("GetContentPreviewInfo() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && !exc.Is(err, tt.wantErr) {
				t.Errorf("GetContentPreviewInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
