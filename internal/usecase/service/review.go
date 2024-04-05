package service

import (
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/entity/dto"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/repository"
	"github.com/go-park-mail-ru/2024_1_Cyberkotletki/internal/usecase"
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
	author, err := r.userRepo.GetUser(map[string]interface{}{"id": reviewEntity.AuthorID})
	if err != nil {
		return nil, err
	}
	var authorName string
	if author.Name == "" {
		// если пользователь не указал имя, то показываем его мейл
		authorName = author.Email
	} else {
		authorName = author.Name
	}
	var avatar string
	// аватарки может и не быть при
	if author.AvatarUploadID != 0 {
		avatar, err = r.staticRepo.GetStatic(author.AvatarUploadID)
		if err != nil {
			return nil, err
		}
	}
	content, err := r.contentRepo.GetContent(reviewEntity.ContentID)
	if err != nil {
		return nil, err
	}
	return &dto.ReviewResponse{
		Review: dto.Review{
			ID:        reviewEntity.ID,
			AuthorID:  reviewEntity.AuthorID,
			ContentID: reviewEntity.ContentID,
			Rating:    reviewEntity.Rating,
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
func (r *ReviewService) reviewEntitiesToDTO(reviewEntities []*entity.Review) (*[]dto.ReviewResponse, error) {
	reviews := make([]dto.ReviewResponse, 0)
	for _, review := range reviewEntities {
		toDTO, err := r.reviewEntityToDTO(review)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, *toDTO)
	}
	return &reviews, nil
}

// GetLatestReviews возвращает последние count отзывов
func (r *ReviewService) GetLatestReviews(count int) (*[]dto.ReviewResponse, error) {
	reviewEntities, err := r.reviewRepo.GetLatestReviews(count)
	if err != nil {
		return nil, err
	}
	return r.reviewEntitiesToDTO(reviewEntities)
}

// GetUserReviews возвращает count отзывов на странице с номером page пользователя с userID
func (r *ReviewService) GetUserReviews(userID, count, page int) (*[]dto.ReviewResponse, error) {
	reviewEntities, err := r.reviewRepo.GetReviewsByAuthorID(userID, page, count)
	if err != nil {
		return nil, err
	}
	return r.reviewEntitiesToDTO(reviewEntities)
}

func (r *ReviewService) GetContentReviews(contentID, count, page int) (*[]dto.ReviewResponse, error) {
	reviewEntities, err := r.reviewRepo.GetReviewsByContentID(contentID, page, count)
	if err != nil {
		return nil, err
	}
	return r.reviewEntitiesToDTO(reviewEntities)
}

func (r *ReviewService) GetReview(reviewID int) (*dto.ReviewResponse, error) {
	reviewEntity, err := r.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		return nil, err
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) CreateReview(review dto.ReviewCreate) (*dto.ReviewResponse, error) {
	if err := entity.ValidateReview(review.Rating, review.Title, review.Text); err != nil {
		return nil, err
	}
	reviewEntity, err := r.reviewRepo.AddReview(&entity.Review{
		AuthorID:  review.UserID,
		ContentID: review.ContentID,
		Rating:    review.Rating,
		Title:     review.Title,
		Text:      review.Text,
	})
	if err != nil {
		return nil, err
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) EditReview(review dto.ReviewUpdate) (*dto.ReviewResponse, error) {
	currentReview, err := r.reviewRepo.GetReviewByID(review.ReviewID)
	if err != nil {
		return nil, err
	}
	if currentReview.AuthorID != review.UserID {
		return nil, entity.NewClientError("Вы не можете редактировать чужой отзыв", entity.ErrForbidden)
	}
	if err = entity.ValidateReview(review.Rating, review.Title, review.Title); err != nil {
		return nil, err
	}
	reviewEntity, err := r.reviewRepo.UpdateReview(&entity.Review{
		ID:       review.ReviewID,
		AuthorID: review.UserID,
		Rating:   review.Rating,
		Title:    review.Title,
		Text:     review.Text,
	})
	if err != nil {
		return nil, err
	}
	return r.reviewEntityToDTO(reviewEntity)
}

func (r *ReviewService) DeleteReview(reviewID, userID int) error {
	review, err := r.reviewRepo.GetReviewByID(reviewID)
	if err != nil {
		return err
	}
	if review.AuthorID != userID {
		return entity.NewClientError("Вы не можете удалить чужой отзыв", entity.ErrForbidden)
	}
	return r.reviewRepo.DeleteReviewByID(reviewID)
}

func (r *ReviewService) LikeReview(userID, reviewID int) error {
	return r.reviewRepo.LikeReview(reviewID, userID, true)
}

func (r *ReviewService) DislikeReview(userID, reviewID int) error {
	return r.reviewRepo.LikeReview(reviewID, userID, false)
}

func (r *ReviewService) IsLikedByUser(userID, reviewID int) (int, error) {
	return r.reviewRepo.IsLikedByUser(reviewID, userID)
}

func (r *ReviewService) UnlikeReview(userID, reviewID int) error {
	return r.reviewRepo.UnlikeReview(reviewID, userID)
}

func (r *ReviewService) GetUserRating(userID int) (int, error) {
	return r.reviewRepo.GetAuthorRating(userID)
}

func (r *ReviewService) GetContentRating(contentID int) (int, error) {
	return r.reviewRepo.GetContentRating(contentID)
}
