package service

import (
	"errors"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
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
			Name:        "Successful retrieval of compilation types",
			ExpectedErr: nil,
			ExpectedOutput: &dto.CompilationTypeResponseList{
				CompilationTypes: []dto.CompilationType{
					{ID: 1, Type: "Type1"},
					{ID: 2, Type: "Type2"},
				},
			},
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetAllCompilationTypes().Return([]entity.CompilationType{
					{ID: 1, Name: "Type1"},
					{ID: 2, Name: "Type2"},
				}, nil)
			},
		},
		{
			Name:           "Unexpected error",
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении типов подборок"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetAllCompilationTypes().Return(nil, errors.New("unexpected error"))
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
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			compService := NewCompilationService(mockCompilationRepo,
				mockStaticRepo, mockContentRepo)
			result, err := compService.GetCompilationTypes()
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				require.Equal(t, tc.ExpectedOutput, result)
			}
		})
	}
}

func TestCompilationService_GetCompilationsByCompilationType(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                     string
		CompTypeID               int
		ExpectedErr              error
		ExpectedOutput           *dto.CompilationResponseList
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
	}{
		{
			Name:           "Unexpected error",
			CompTypeID:     1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении подборок по типу"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationsByTypeID(1).Return(nil, errors.New("unexpected error"))
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
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			compService := NewCompilationService(mockCompilationRepo,
				mockStaticRepo, mockContentRepo)
			result, err := compService.GetCompilationsByCompilationType(tc.CompTypeID)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				require.Equal(t, tc.ExpectedOutput, result)
			}
		})
	}
}

func TestCompilationService_GetCompilationContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                     string
		CompID                   int
		Page                     int
		ExpectedErr              error
		ExpectedOutput           *dto.CompilationResponse
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
		SetupStaticRepoMock      func(repo *mockrepo.MockStatic)
	}{

		{
			Name:           "Unexpected error",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilation(1).Return(entity.Compilation{}, errors.New("unexpected error"))
			},
		},
		{
			Name:           "Compilation not found",
			CompID:         1,
			Page:           1,
			ExpectedErr:    usecase.ErrCompilationNotFound,
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilation(1).Return(entity.Compilation{}, repository.ErrCompilationNotFound)
			},
		},
		{
			Name:           "Error getting compilation content",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении контента подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilation(1).Return(entity.Compilation{ID: 1}, nil)
				repo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return(nil, errors.New("unexpected error"))
			},
		},
		{
			Name:           "Error getting compilation content length",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении количества контента подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilation(1).Return(entity.Compilation{ID: 1}, nil)
				repo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return([]int{1, 2, 3}, nil)
				repo.EXPECT().GetCompilationContentLength(1).Return(0, errors.New("unexpected error"))
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
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			compService := NewCompilationService(mockCompilationRepo,
				mockStaticRepo, mockContentRepo)
			result, err := compService.GetCompilationContent(tc.CompID, tc.Page)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				require.Equal(t, tc.ExpectedOutput, result)
			}
		})
	}
}
