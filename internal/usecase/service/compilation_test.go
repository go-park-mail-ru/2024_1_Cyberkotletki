package service

import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCompilationService_GetCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                     string
		ExpectedErr              error
		ExpectedOutput           *dto.CompilationTypeResponseList
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
	}{
		{
			Name:           "Ошибка получения",
			ExpectedErr:    fmt.Errorf("ошибка получения типов подборок"),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetAllCompilationTypes().Return(nil, fmt.Errorf("ошибка получения типов подборок"))
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCompilationRepo := mockrepo.NewMockCompilation(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			compilationService := NewCompilationService(mockCompilationRepo, mockStaticRepo, mockContentRepo, mockReviewRepo)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			output, err := compilationService.GetCompilationTypes()
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}

func TestCompilationService_GetCompilationsByCompilationType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                     string
		Input                    int
		ExpectedErr              error
		ExpectedOutput           *dto.CompilationResponseList
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
		SetupStaticRepoMock      func(repo *mockrepo.MockStatic)
		SetupContentRepoMock     func(repo *mockrepo.MockContent)
		SetupReviewRepoMock      func(repo *mockrepo.MockReview)
	}{
		{
			Name:           "Ошибка получения",
			Input:          1,
			ExpectedErr:    fmt.Errorf("ошибка получения подборок"),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationsByTypeID(1).Return(nil, fmt.Errorf("ошибка получения подборок"))
			},
			SetupStaticRepoMock:  func(repo *mockrepo.MockStatic) {},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
			SetupReviewRepoMock:  func(repo *mockrepo.MockReview) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockCompilationRepo := mockrepo.NewMockCompilation(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockReviewRepo := mockrepo.NewMockReview(ctrl)
			compilationService := NewCompilationService(mockCompilationRepo, mockStaticRepo,
				mockContentRepo, mockReviewRepo)

			tc.SetupCompilationRepoMock(mockCompilationRepo)
			tc.SetupStaticRepoMock(mockStaticRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			tc.SetupReviewRepoMock(mockReviewRepo)
			output, err := compilationService.GetCompilationsByCompilationType(tc.Input)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}
