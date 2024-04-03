package service

import (
	"fmt"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestAuth_Logout(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                string
		ExpectedErr          error
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:        "Существующая сессия",
			Input:       "session1",
			ExpectedErr: nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session1").Return(nil)
			},
		},
		{
			Name:        "Несуществующая сессия",
			Input:       "session2",
			ExpectedErr: fmt.Errorf("не удалось удалить сессию"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteSession("session2").Return(fmt.Errorf("не удалось удалить сессию"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := mockrepo.NewMockSession(ctrl)
			authService := AuthService{
				sessionRepo: mockSessionRepo,
			}
			tc.SetupSessionRepoMock(mockSessionRepo)
			err := authService.Logout(tc.Input)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestAuth_LogoutAll(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                int
		ExpectedErr          error
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:        "Существующий пользователь",
			Input:       1,
			ExpectedErr: nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteAllSessions(1).Return(nil)
			},
		},
		{
			Name:        "Несуществующий пользователь",
			Input:       2,
			ExpectedErr: fmt.Errorf("не удалось удалить сессии"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().DeleteAllSessions(2).Return(fmt.Errorf("не удалось удалить сессии"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := mockrepo.NewMockSession(ctrl)
			authService := AuthService{
				sessionRepo: mockSessionRepo,
			}
			tc.SetupSessionRepoMock(mockSessionRepo)
			err := authService.LogoutAll(tc.Input)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestAuth_GetUserIDBySession(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                string
		ExpectedUserID       int
		ExpectedErr          error
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:           "Существующая сессия",
			Input:          "session1",
			ExpectedUserID: 1,
			ExpectedErr:    nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session1").Return(1, nil)
			},
		},
		{
			Name:           "Несуществующая сессия",
			Input:          "session2",
			ExpectedUserID: 0,
			ExpectedErr:    fmt.Errorf("не удалось проверить сессию"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().CheckSession("session2").Return(0, fmt.Errorf("не удалось проверить сессию"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := mockrepo.NewMockSession(ctrl)
			authService := AuthService{
				sessionRepo: mockSessionRepo,
			}
			tc.SetupSessionRepoMock(mockSessionRepo)
			userID, err := authService.GetUserIDBySession(tc.Input)
			require.EqualValues(t, tc.ExpectedUserID, userID)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}

func TestAuth_CreateSession(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                 string
		Input                int
		ExpectedSession      string
		ExpectedErr          error
		SetupSessionRepoMock func(repo *mockrepo.MockSession)
	}{
		{
			Name:            "Существующий пользователь",
			Input:           1,
			ExpectedSession: "session1",
			ExpectedErr:     nil,
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(1).Return("session1", nil)
			},
		},
		{
			Name:            "Несуществующий пользователь",
			Input:           2,
			ExpectedSession: "",
			ExpectedErr:     fmt.Errorf("не удалось создать сессию"),
			SetupSessionRepoMock: func(repo *mockrepo.MockSession) {
				repo.EXPECT().NewSession(2).Return("", fmt.Errorf("не удалось создать сессию"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockSessionRepo := mockrepo.NewMockSession(ctrl)
			authService := AuthService{
				sessionRepo: mockSessionRepo,
			}
			tc.SetupSessionRepoMock(mockSessionRepo)
			session, err := authService.CreateSession(tc.Input)
			require.EqualValues(t, tc.ExpectedSession, session)
			require.EqualValues(t, tc.ExpectedErr, err)
		})
	}
}
