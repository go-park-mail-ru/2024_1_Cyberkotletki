package service

import (
	"errors"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
	"strings"
)

type ReviewService struct {
	reviewRepo  repository.Review
	userRepo    repository.User
	contentRepo repository.Content
	staticRepo  repository.Static
}

func NewReviewService(
	reviewRepo repository.Review,
	userRepo repository.User,
	contentRepo repository.Content,
	staticRepo repository.Static,
) usecase.Review {
	return &ReviewService{
		reviewRepo:  reviewRepo,
		userRepo:    userRepo,
		contentRepo: contentRepo,
		staticRepo:  staticRepo,
	}
}

// reviewEntityToDTO конвертирует entity.Review в dto.ReviewResponse, добавляя дополнительные поля автора и контента
func (r *ReviewService) reviewEntityToDTO(reviewEntity *entity.Review) (*dto.ReviewResponse, error) {
	author, err := r.userRepo.GetUserByID(reviewEntity.AuthorID)
	switch {
	case errors.Is(err, repository.ErrUserNotFound):
		return nil, usecase.ErrUserNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении пользователя"), err)
	}
	var authorName string
	authorName = author.Name
	if author.Name == "" {
		// если пользователь не указал имя, то показываем его мейл
		authorName = author.Email
	}
	var avatar string
	avatar, err = r.staticRepo.GetStatic(author.AvatarUploadID)
	switch {
	case errors.Is(err, repository.ErrStaticNotFound):
		// аватара может и не быть
		avatar = ""
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении аватара"), err)
	}
	content, err := r.contentRepo.GetContent(reviewEntity.ContentID)
	switch {
	case errors.Is(err, repository.ErrContentNotFound):
		return nil, usecase.ErrContentNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении контента"), err)
	}
	return &dto.ReviewResponse{
		Review: dto.Review{
			ID:        reviewEntity.ID,
			AuthorID:  reviewEntity.AuthorID,
			ContentID: reviewEntity.ContentID,
			Rating:    reviewEntity.ContentRating,
			Title:     reviewEntity.Title,
			Text:      reviewEntity.Text,
			CreatedAt: reviewEntity.CreatedAt.String(),
			Likes:     reviewEntity.Likes,
			Dislikes:  reviewEntity.Dislikes,
		},
		AuthorName:   authorName,
		AuthorAvatar: avatar,
		ContentName:  content.Title,
	}, nil
}

// reviewEntitiesToDTO конвертирует массив entity.Review в массив dto.ReviewResponse
func (r *ReviewService) reviewEntitiesToDTO(reviewEntities []*entity.Review) (*dto.ReviewResponseList, error) {
	reviews := make([]dto.ReviewResponse, 0, len(reviewEntities))
	for i, review := range reviewEntities {
		toDTO, err := r.reviewEntityToDTO(review)
		if err != nil {
			return nil, err
		}
		reviews[i] = *toDTO
	}
	return &dto.ReviewResponseList{Reviews: reviews}, nil
}

// GetLatestReviews возвращает последние count отзывов
func (r *ReviewService) GetLatestReviews(count int) (*dto.ReviewResponseList, error) {
	reviewEntities, err := r.reviewRepo.GetLatestReviews(count)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении последних отзывов"), err)
	}
	reviewDTOs, err := r.reviewEntitiesToDTO(reviewEntities)
	if err != nil {
		return nil, err
	}
	reviewDTOs.Page = 1
	reviewDTOs.Count = count
	reviewDTOs.Pages = 1
	reviewDTOs.Total = count
	return reviewDTOs, nil
}

// GetUserReviews возвращает count отзывов на странице с номером page пользователя с userID
func (r *ReviewService) GetUserReviews(userID, count, page int) (*dto.ReviewResponseList, error) {
	reviewEntities, err := r.reviewRepo.GetReviewsByAuthorID(userID, page, count)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзывов пользователя"), err)
	}
	reviewDTOs, err := r.reviewEntitiesToDTO(reviewEntities)
	if err != nil {
		return nil, err
	}
	reviewDTOs.Total, err = r.reviewRepo.GetReviewsCountByAuthorID(userID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении количества отзывов пользователя"), err)
	}
	reviewDTOs.Page = page
	reviewDTOs.Count = count
	reviewDTOs.Pages = (reviewDTOs.Total + count - 1) / count
	return reviewDTOs, nil
}

func (r *ReviewService) GetContentReviews(contentID, count, page int) (*dto.ReviewResponseList, error) {
	reviewEntities, err := r.reviewRepo.GetReviewsByContentID(contentID, page, count)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзывов контента"), err)
	}
	reviewDTOs, err := r.reviewEntitiesToDTO(reviewEntities)
	if err != nil {
		return nil, err
	}
	reviewDTOs.Total, err = r.reviewRepo.GetReviewsCountByContentID(contentID)
	if err != nil {
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении количества отзывов пользователя"), err)
	}
	reviewDTOs.Page = page
	reviewDTOs.Count = count
	reviewDTOs.Pages = (reviewDTOs.Total + count - 1) / count
	return reviewDTOs, nil
}

func (r *ReviewService) GetReview(reviewID int) (*dto.ReviewResponse, error) {
	reviewEntity, err := r.reviewRepo.GetReviewByID(reviewID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return nil, usecase.ErrReviewNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), err)
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) GetContentReviewByAuthor(authorID, contentID int) (*dto.ReviewResponse, error) {
	reviewEntity, err := r.reviewRepo.GetContentReviewByAuthor(authorID, contentID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return nil, usecase.ErrReviewNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), err)
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) CreateReview(review dto.ReviewCreate) (*dto.ReviewResponse, error) {
	if err := entity.ValidateReview(review.Rating, review.Title, review.Text); err != nil {
		return nil, usecase.ReviewErrorIncorrectData{Err: err}
	}
	reviewEntity, err := r.reviewRepo.AddReview(&entity.Review{
		AuthorID:      review.UserID,
		ContentID:     review.ContentID,
		ContentRating: review.Rating,
		Title:         strings.TrimSpace(review.Title),
		Text:          strings.TrimSpace(review.Text),
	})
	switch {
	case errors.Is(err, repository.ErrReviewViolation):
		return nil, usecase.ErrReviewContentNotFound
	case errors.Is(err, repository.ErrReviewBadRequest):
		return nil, usecase.ReviewErrorIncorrectData{Err: err}
	case errors.Is(err, repository.ErrReviewAlreadyExists):
		return nil, usecase.ErrReviewAlreadyExists
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при добавлении отзыва"), err)
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) EditReview(review dto.ReviewUpdate) (*dto.ReviewResponse, error) {
	currentReview, err := r.reviewRepo.GetReviewByID(review.ReviewID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return nil, usecase.ErrReviewNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), err)
	}
	if currentReview.AuthorID != review.UserID {
		return nil, usecase.ErrReviewForbidden
	}
	if err = entity.ValidateReview(review.Rating, review.Title, review.Title); err != nil {
		return nil, usecase.ReviewErrorIncorrectData{Err: err}
	}
	err = r.reviewRepo.UpdateReview(&entity.Review{
		ID:            review.ReviewID,
		AuthorID:      review.UserID,
		ContentRating: review.Rating,
		Title:         strings.TrimSpace(review.Title),
		Text:          strings.TrimSpace(review.Text),
	})
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return nil, usecase.ErrReviewNotFound
	case errors.Is(err, repository.ErrReviewBadRequest):
		return nil, usecase.ReviewErrorIncorrectData{Err: err}
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при обновлении отзыва"), err)
	}
	reviewEntity, err := r.reviewRepo.GetReviewByID(review.ReviewID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return nil, usecase.ErrReviewNotFound
	case err != nil:
		return nil, entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), err)
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) DeleteReview(reviewID, userID int) error {
	review, err := r.reviewRepo.GetReviewByID(reviewID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return usecase.ErrReviewNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при получении отзыва"), err)
	}
	if review.AuthorID != userID {
		return usecase.ErrReviewForbidden
	}
	err = r.reviewRepo.DeleteReviewByID(reviewID)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return usecase.ErrReviewNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при удалении отзыва"), err)
	}
	return nil
}

func (r *ReviewService) VoteReview(reviewID, userID int, vote bool) error {
	userVote, err := r.reviewRepo.IsVotedByUser(reviewID, userID)
	if err != nil {
		return entity.UsecaseWrap(errors.New("ошибка при проверке оценки"), err)
	}
	switch {
	case userVote == 1 && vote:
		// пользователь уже поставил лайк
		return nil
	case userVote == -1 && !vote:
		// пользователь уже поставил дизлайк
		return nil
	case userVote == 1:
		// пользователь поставил дизлайк, но хочет поставить лайк
		err = r.reviewRepo.UnVoteReview(reviewID, userID)
	case userVote == -1:
		// пользователь поставил лайк, но хочет поставить дизлайк
		err = r.reviewRepo.UnVoteReview(reviewID, userID)
	}
	switch {
	case errors.Is(err, repository.ErrReviewVoteNotFound):
		// пользователь еще не оценивал отзыв - так и должно быть
		break
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при удалении оценки"), err)
	}
	err = r.reviewRepo.VoteReview(reviewID, userID, vote)
	switch {
	case errors.Is(err, repository.ErrReviewNotFound):
		return usecase.ErrReviewNotFound
	case errors.Is(err, repository.ErrReviewVoteAlreadyExists):
		return usecase.ErrReviewVoteAlreadyExists
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при оценке отзыва"), err)
	}
	return nil
}

func (r *ReviewService) IsVotedByUser(userID, reviewID int) (int, error) {
	isVoted, err := r.reviewRepo.IsVotedByUser(reviewID, userID)
	if err != nil {
		return 0, entity.UsecaseWrap(errors.New("ошибка при проверке оценки"), err)
	}
	return isVoted, nil
}

func (r *ReviewService) UnVoteReview(userID, reviewID int) error {
	err := r.reviewRepo.UnVoteReview(reviewID, userID)
	switch {
	case errors.Is(err, repository.ErrReviewVoteNotFound):
		return usecase.ErrReviewVoteNotFound
	case err != nil:
		return entity.UsecaseWrap(errors.New("ошибка при удалении оценки"), err)
	}
	return nil
}
