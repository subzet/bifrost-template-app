package services

import (
	"context"
	"errors"
	"net/http"

	"myapp/config"
	"myapp/model"
	"myapp/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	repo *model.UserRepository
}

func NewAuthService(repo *model.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, email, password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &model.User{Email: email, PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, user); err != nil {
		return "", ErrEmailTaken
	}

	return signToken(user.ID)
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return signToken(user.ID)
}

func (s *AuthService) GetUserFromRequest(r *http.Request) *model.User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	claims, appErr := util.ParseJwt(config.Env.JWT_SECRET, cookie.Value)
	if appErr != nil {
		return nil
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return nil
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil
	}

	user, err := s.repo.GetByID(r.Context(), userID.String())
	if err != nil {
		return nil
	}
	return user
}

func signToken(userID uuid.UUID) (string, error) {
	token, appErr := util.SignJwt(config.Env.JWT_SECRET, map[string]any{
		"sub": userID.String(),
	})
	if appErr != nil {
		return "", appErr.Error
	}
	return *token, nil
}
