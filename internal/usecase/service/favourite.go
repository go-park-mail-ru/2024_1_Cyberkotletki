package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
)

type FavouriteService struct {
	contentUC     usecase.Content
	favouriteRepo repository.Favourite
}

func NewFavouriteService(favouriteRepo repository.Favourite, contentUC usecase.Content) usecase.Favourite {
	return &FavouriteService{
		favouriteRepo: favouriteRepo,
		contentUC:     contentUC,
	}
}

func (f FavouriteService) CreateFavourite(userID, contentID int, category string) error {
	currentStatus, err := f.GetStatus(userID, contentID)
	switch {
	case errors.Is(err, usecase.ErrFavouriteNotFound):
		// всё ок, добавляем статус
		break
	case err != nil:
		return entity.UsecaseWrap(err, errors.New("ошибка при получении статуса контента в избранном в FavouriteService"))
	default:
		if currentStatus.Status == category {
			return nil
		}
		// сначала удаляем старый статус
		err = f.favouriteRepo.DeleteFavourite(userID, contentID)
		if err != nil {
			return entity.UsecaseWrap(err, errors.New("ошибка при удалении из избранного в FavouriteService"))
		}
	}

	err = f.favouriteRepo.CreateFavourite(userID, contentID, category)
	switch {
	case errors.Is(err, repository.ErrFavouriteContentNotFound):
		return usecase.ErrFavouriteContentNotFound
	case err != nil:
		return entity.UsecaseWrap(err, errors.New("ошибка при добавлении в избранное в FavouriteService"))
	default:
		return nil
	}
}

func (f FavouriteService) DeleteFavourite(userID, contentID int) error {
	err := f.favouriteRepo.DeleteFavourite(userID, contentID)
	switch {
	case errors.Is(err, repository.ErrFavouriteNotFound):
		return usecase.ErrFavouriteNotFound
	case err != nil:
		return entity.UsecaseWrap(err, errors.New("ошибка при удалении из избранного в FavouriteService"))
	default:
		return nil
	}
}

func (f FavouriteService) GetFavourites(userID int) (*dto.FavouritesResponse, error) {
	favourites, err := f.favouriteRepo.GetFavourites(userID)
	switch {
	case errors.Is(err, repository.ErrFavouriteUserNotFound):
		return nil, usecase.ErrFavouriteUserNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(err, errors.New("ошибка при получении избранного контента в FavouriteService"))
	}
	response := dto.FavouritesResponse{
		Favourites: make([]dto.Favourite, len(favourites)),
	}
	for index, favourite := range favourites {
		content, err := f.contentUC.GetPreviewContentByID(favourite.ContentID)
		if err != nil {
			return nil, entity.UsecaseWrap(err, errors.New("ошибка при получении контента из избранного в FavouriteService"))
		}
		response.Favourites[index] = dto.Favourite{
			Content:  *content,
			Category: favourite.Category,
		}
	}
	return &response, nil
}

func (f FavouriteService) GetStatus(userID, contentID int) (*dto.FavouriteStatusResponse, error) {
	status, err := f.favouriteRepo.GetFavourite(userID, contentID)
	switch {
	case errors.Is(err, repository.ErrFavouriteNotFound):
		return nil, usecase.ErrFavouriteNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(err,
			errors.New("ошибка при получении статуса контента в избранном в FavouriteService"))
	default:
		return &dto.FavouriteStatusResponse{
			Status: status.Category,
		}, nil
	}
}
