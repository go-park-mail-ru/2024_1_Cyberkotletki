package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"testing"

	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	mock_usecase "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase/mocks"
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
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			mockStaticUC := mock_usecase.NewMockStatic(ctrl)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			compService := NewCompilationService(mockCompilationRepo, mockStaticUC, mockContentUC)
			result, err := compService.GetCompilationTypes()
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, result)
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
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			mockStaticUC := mock_usecase.NewMockStatic(ctrl)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			compService := NewCompilationService(mockCompilationRepo, mockStaticUC, mockContentUC)
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
		Name           string
		CompID         int
		Page           int
		ExpectedErr    error
		ExpectedOutput *dto.CompilationResponse
		SetupMock      func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic)
	}{
		{
			Name:           "Подборка не найдена",
			CompID:         1,
			Page:           1,
			ExpectedErr:    usecase.ErrCompilationNotFound,
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(nil, repository.ErrCompilationNotFound)
			},
		},
		{
			Name:           "Неожиданная ошибка при поиске подборки",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(nil, errors.New("unexpected error"))
			},
		},
		{
			Name:           "Контент в подборке не найден",
			CompID:         1,
			Page:           1,
			ExpectedErr:    usecase.ErrCompilationNotFound,
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return(nil, repository.ErrCompilationNotFound)
			},
		},
		{
			Name:           "Неожиданная ошибка при поиске контента в подборке",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении контента подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return(nil, errors.New("unexpected error"))
			},
		},
		{
			Name:           "Контент не найден",
			CompID:         1,
			Page:           1,
			ExpectedErr:    usecase.ErrContentNotFound,
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return([]int{1}, nil)
				contentUC.EXPECT().GetPreviewContentByID(1).Return(nil, usecase.ErrContentNotFound)
			},
		},
		{
			Name:           "Неожиданная ошибка при поиске контента",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении контента"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return([]int{1}, nil)
				contentUC.EXPECT().GetPreviewContentByID(1).Return(nil, errors.New("unexpected error"))
			},
		},
		{
			Name:           "Неожиданная ошибка при поиске количества контента в подборке",
			CompID:         1,
			Page:           1,
			ExpectedErr:    entity.UsecaseWrap(errors.New("ошибка при получении количества контента подборки"), errors.New("unexpected error")),
			ExpectedOutput: nil,
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return([]int{1}, nil)
				contentUC.EXPECT().GetPreviewContentByID(1).Return(&dto.PreviewContent{}, nil)
				compRepo.EXPECT().GetCompilationContentLength(1).Return(1, errors.New("unexpected error"))
			},
		},
		{
			Name:        "Успех",
			CompID:      1,
			Page:        1,
			ExpectedErr: nil,
			ExpectedOutput: &dto.CompilationResponse{
				Compilation: dto.Compilation{
					ID:                1,
					Title:             "",
					CompilationTypeID: 1,
					PosterURL:         "",
				},
				Content: []*dto.PreviewContent{
					{ID: 1},
				},
				ContentLength: 1,
				Page:          1,
				PerPage:       compilationContentLimit,
				TotalPages:    1,
			},
			SetupMock: func(compRepo *mockrepo.MockCompilation, contentUC *mock_usecase.MockContent, staticUC *mock_usecase.MockStatic) {
				compRepo.EXPECT().GetCompilation(1).Return(&entity.Compilation{ID: 1, CompilationTypeID: 1, PosterUploadID: 1}, nil)
				compRepo.EXPECT().GetCompilationContent(1, 1, compilationContentLimit).Return([]int{1}, nil)
				contentUC.EXPECT().GetPreviewContentByID(1).Return(&dto.PreviewContent{ID: 1}, nil)
				compRepo.EXPECT().GetCompilationContentLength(1).Return(1, nil)
				staticUC.EXPECT().GetStatic(1).Return("", nil)
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
			mockContentUC := mock_usecase.NewMockContent(ctrl)
			mockStaticUC := mock_usecase.NewMockStatic(ctrl)
			tc.SetupMock(mockCompilationRepo, mockContentUC, mockStaticUC)
			compService := NewCompilationService(mockCompilationRepo, mockStaticUC, mockContentUC)
			result, err := compService.GetCompilationContent(tc.CompID, tc.Page)
			require.Equal(t, tc.ExpectedErr, err)
			if tc.ExpectedErr == nil {
				require.Equal(t, tc.ExpectedOutput, result)
			}
		})
	}
}
