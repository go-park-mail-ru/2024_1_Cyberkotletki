package exceptions

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		l    Layer
		t    Type
		msgs []string
	}{
		{
			name: "0 messages",
			l:    Service,
			t:    Forbidden,
			msgs: []string{},
		},
		{
			name: "1 message",
			l:    Database,
			t:    NotFound,
			msgs: []string{"message1"},
		},
		{
			name: "2 messages",
			l:    Server,
			t:    Internal,
			msgs: []string{"message1", "message2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.l, tt.t, tt.msgs...)
			if err == nil {
				t.Errorf("New() error = nil")
			}
		})
	}
}

func TestError_Error(t *testing.T) {
	err := New(Service, Forbidden, "тестовое сообщение")
	want := "[" + err.(Error).When.String() +
		"] ошибка бизнес-слоя: нет доступа: : тестовое сообщение"

	if got := err.Error(); got != want {
		t.Errorf("Error() = %v, want %v", got, want)
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		name string
		err1 error
		err2 error
		want bool
	}{
		{
			name: "Same type",
			err1: New(Service, Forbidden),
			err2: ForbiddenErr,
			want: true,
		},
		{
			name: "Different type",
			err1: New(Service, Forbidden),
			err2: NotFoundErr,
			want: false,
		},
		{
			name: "Non-Error type",
			err1: errors.New("неизведанная ошибка"),
			err2: ForbiddenErr,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.err1, tt.err2); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
