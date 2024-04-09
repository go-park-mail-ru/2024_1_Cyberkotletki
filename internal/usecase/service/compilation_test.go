package service

/*
import (
	"fmt"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	mockrepo "github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCompilationService_GetCompilationTypes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                         string
		ExpectedErr                  error
		ExpectedOutput               *dto.CompilationTypeResponseList
		SetupCompilationTypeRepoMock func(repo *mockrepo.MockCompilationType)
	}{
		{
			Name:        "Успешное получение",
			ExpectedErr: nil,
			ExpectedOutput: &dto.CompilationTypeResponseList{
				CompilationTypes: []dto.CompilationTypeResponse{
					{
						CompilationType: dto.CompilationType{
							ID:   1,
							Type: "TestType1",
						},
					},
					{
						CompilationType: dto.CompilationType{
							ID:   2,
							Type: "TestType2",
						},
					},
				},
			},
			SetupCompilationTypeRepoMock: func(repo *mockrepo.MockCompilationType) {
				repo.EXPECT().GetAllCompilationTypes().Return([]*entity.CompilationType{
					{
						ID:   1,
						Type: "TestType1",
					},
					{
						ID:   2,
						Type: "TestType2",
					},
				}, nil)
			},
		},
		{
			Name:           "Ошибка получения",
			ExpectedErr:    fmt.Errorf("ошибка получения типов подборок"),
			ExpectedOutput: nil,
			SetupCompilationTypeRepoMock: func(repo *mockrepo.MockCompilationType) {
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
			mockCompilationTypeRepo := mockrepo.NewMockCompilationType(ctrl)
			mockCompilationRepo := mockrepo.NewMockCompilation(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			compilationService := NewCompilationService(mockCompilationRepo, mockCompilationTypeRepo, mockStaticRepo, mockContentRepo)
			tc.SetupCompilationTypeRepoMock(mockCompilationTypeRepo)
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
		InputCompTypeID          int
		InputCount               int
		InputPage                int
		ExpectedErr              error
		ExpectedOutput           *dto.CompilationResponseList
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
	}{
		{
			Name:            "Успешное получение",
			InputCompTypeID: 1,
			InputCount:      10,
			InputPage:       1,
			ExpectedErr:     nil,
			ExpectedOutput: &dto.CompilationResponseList{
				Compilations: []dto.CompilationResponse{
					{
						Compilation: dto.Compilation{
							ID:                1,
							Title:             "TestTitle1",
							CompilationTypeID: 1,
							PosterUploadID:    1,
						},
						ContentLength: 120,
					},
					{
						Compilation: dto.Compilation{
							ID:                2,
							Title:             "TestTitle2",
							CompilationTypeID: 1,
							PosterUploadID:    2,
						},
						ContentLength: 150,
					},
				},
			},
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationsByCompilationTypeID(1, 10, 1).Return([]*entity.Compilation{
					{
						ID:                1,
						Title:             "TestTitle1",
						CompilationTypeID: 1,
						PosterUploadID:    1,
					},
					{
						ID:                2,
						Title:             "TestTitle2",
						CompilationTypeID: 1,
						PosterUploadID:    2,
					},
				}, nil)
			},
		},
		{
			Name:            "Ошибка получения",
			InputCompTypeID: 1,
			InputCount:      10,
			InputPage:       1,
			ExpectedErr:     fmt.Errorf("ошибка получения подборок"),
			ExpectedOutput:  nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationsByCompilationTypeID(1, 10, 1).Return(nil, fmt.Errorf("ошибка получения подборок"))
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
			mockCompilationTypeRepo := mockrepo.NewMockCompilationType(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			compilationService := NewCompilationService(mockCompilationRepo, mockCompilationTypeRepo, mockStaticRepo, mockContentRepo)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			output, err := compilationService.GetCompilationsByCompilationType(tc.InputCompTypeID, tc.InputCount, tc.InputPage)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}

func TestCompilationService_GetCompilationContent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name                    string
		InputCompID             int
		ExpectedErr             error
		ExpectedOutput          []*dto.PreviewContentCardResponse
		SetupCompilationRepoMock func(repo *mockrepo.MockCompilation)
		SetupContentRepoMock     func(repo *mockrepo.MockContent)
	}{
		{
			Name:        "Успешное получение",
			InputCompID: 1,
			ExpectedErr: nil,
			ExpectedOutput: []*dto.PreviewContentCardResponse{
				{
					PreviewContentCard: dto.PreviewContentCard{
						ID:            1,
						Title:         "TestTitle1",
						OriginalTitle: "TestOriginalTitle1",
						Country:       "TestCountry1",
						Genre:         "TestGenre1",
						Director:      "TestDirector1",
						Actors:        []string{"TestActor1", "TestActor2"},
						Poster:        "TestPoster1",
						Rating:        10.0,
					},
					Type: entity.ContentTypeMovie,
				},
			},
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationContent(1, 0, 2).Return([]int{1}, nil)
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {
				repo.EXPECT().GetPreviewContent(1).Return(&entity.Content{
					ID:            1,
					Title:         "TestTitle1",
					OriginalTitle: "TestOriginalTitle1",
					Country:       []entity.Country{{Name: "TestCountry1"}},
					Genres:        []entity.Genre{{Name: "TestGenre1"}},
					Directors:     []entity.Person{{FirstName: "TestDirector", LastName: "1"}},
					Actors:        []entity.Person{{FirstName: "TestActor", LastName: "1"}, {FirstName: "TestActor", LastName: "2"}},
					PosterStaticID: 1,
					Type:          entity.ContentTypeMovie,
				}, nil)
			},
		},
		{
			Name:           "Ошибка получения",
			InputCompID:    1,
			ExpectedErr:    fmt.Errorf("ошибка получения контента подборки"),
			ExpectedOutput: nil,
			SetupCompilationRepoMock: func(repo *mockrepo.MockCompilation) {
				repo.EXPECT().GetCompilationContent(1, 0, 2).Return(nil, fmt.Errorf("ошибка получения контента подборки"))
			},
			SetupContentRepoMock: func(repo *mockrepo.MockContent) {},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockCompilationRepo := mockrepo.NewMockCompilation(ctrl)
			mockCompilationTypeRepo := mockrepo.NewMockCompilationType(ctrl)
			mockStaticRepo := mockrepo.NewMockStatic(ctrl)
			mockContentRepo := mockrepo.NewMockContent(ctrl)
			compilationService := NewCompilationService(mockCompilationRepo, mockCompilationTypeRepo, mockStaticRepo, mockContentRepo)
			tc.SetupCompilationRepoMock(mockCompilationRepo)
			tc.SetupContentRepoMock(mockContentRepo)
			output, err := compilationService.GetCompilationContent(tc.InputCompID)
			require.Equal(t, tc.ExpectedErr, err)
			require.Equal(t, tc.ExpectedOutput, output)
		})
	}
}*/
