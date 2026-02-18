package services

import (
	"context"

	"myapp/model"
)

type UserService struct {
	repo *model.UserRepository
}

func NewUserService(repo *model.UserRepository) *UserService {
	return &UserService{repo: repo}
}

type UpdateProfileInput struct {
	Handle      string
	DisplayName string
	Bio         string
	Country     string
	SocialLinks model.SocialLinks
	AvatarURL   string // empty means keep existing avatar
}

func (s *UserService) UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) error {
	if !model.HandleRegex.MatchString(input.Handle) {
		return ErrHandleInvalid
	}

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if input.Handle != user.Name {
		existing, err := s.repo.GetByHandle(ctx, input.Handle)
		if err == nil && existing != nil {
			return ErrHandleTaken
		}
	}

	user.Name = input.Handle
	user.DisplayName = input.DisplayName
	user.Bio = input.Bio
	user.Country = input.Country
	user.SocialLinks = input.SocialLinks
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}

	return s.repo.Update(ctx, user)
}

func (s *UserService) GetByHandle(ctx context.Context, handle string) (*model.User, error) {
	return s.repo.GetByHandle(ctx, handle)
}
