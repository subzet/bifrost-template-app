package model

import (
	"context"
	"errors"
	"fmt"
	"myapp/util"
	"regexp"

	"gorm.io/gorm"
)

var HandleRegex = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{2,29}$`)

type SocialLinks struct {
	Instagram string `json:"instagram"`
	Facebook  string `json:"facebook"`
	Linkedin  string `json:"linkedin"`
	X         string `json:"x"`
}

type User struct {
	util.Entity
	Email        string      `json:"email"        gorm:"uniqueIndex;not null"`
	PasswordHash string      `json:"-"            gorm:"column:password_hash;not null"`
	Name         string      `json:"name"`
	DisplayName  string      `json:"display_name"`
	Bio          string      `json:"bio"`
	Country      string      `json:"country"`
	SocialLinks  SocialLinks `json:"social_links" gorm:"serializer:json"`
	AvatarURL    string      `json:"avatar_url"`
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Where("deleted_at is null").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		Where("deleted_at is null").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) GetByHandle(ctx context.Context, handle string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("name = ?", handle).
		Where("deleted_at is null").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by handle: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Model(&User{}).
		Where("id = ?", id).
		Where("deleted_at is null").
		Update("deleted_at", gorm.Expr("CURRENT_TIMESTAMP"))

	if result.Error != nil {
		return fmt.Errorf("failed to soft-delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

// ExistsByEmail – useful for registration checks
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&User{}).
		Where("email = ?", email).
		Where("deleted_at is null").
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check email existence: %w", err)
	}
	return count > 0, nil
}
