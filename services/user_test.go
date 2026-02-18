package services

import (
	"context"
	"errors"
	"testing"

	"myapp/model"
	"myapp/testutil"
)

func newTestUserService(t *testing.T) (*UserService, *AuthService) {
	t.Helper()
	repo := model.NewUserRepository(testutil.NewTestDB(t, &model.User{}))
	return NewUserService(repo), NewAuthService(repo)
}

func TestUpdateProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userSvc, authSvc := newTestUserService(t)
		_, _ = authSvc.Signup(ctx, "user@example.com", "password123", "oldhandle")

		user, err := userSvc.GetByHandle(ctx, "oldhandle")
		if err != nil {
			t.Fatalf("GetByHandle failed: %v", err)
		}

		input := UpdateProfileInput{
			Handle:      "newhandle",
			DisplayName: "Test User",
			Bio:         "My bio",
			Country:     "US",
			SocialLinks: model.SocialLinks{Instagram: "https://instagram.com/test"},
		}
		if err := userSvc.UpdateProfile(ctx, user.ID.String(), input); err != nil {
			t.Fatalf("UpdateProfile failed: %v", err)
		}

		updated, err := userSvc.GetByHandle(ctx, "newhandle")
		if err != nil {
			t.Fatalf("GetByHandle after update failed: %v", err)
		}
		if updated.DisplayName != "Test User" {
			t.Errorf("got display_name %q, want %q", updated.DisplayName, "Test User")
		}
		if updated.SocialLinks.Instagram != "https://instagram.com/test" {
			t.Errorf("got instagram %q, want %q", updated.SocialLinks.Instagram, "https://instagram.com/test")
		}
	})

	t.Run("invalid handle", func(t *testing.T) {
		userSvc, authSvc := newTestUserService(t)
		_, _ = authSvc.Signup(ctx, "user@example.com", "password123", "testuser")
		user, _ := userSvc.GetByHandle(ctx, "testuser")

		err := userSvc.UpdateProfile(ctx, user.ID.String(), UpdateProfileInput{Handle: "ab"})
		if !errors.Is(err, ErrHandleInvalid) {
			t.Errorf("expected ErrHandleInvalid, got %v", err)
		}
	})

	t.Run("handle taken", func(t *testing.T) {
		userSvc, authSvc := newTestUserService(t)
		_, _ = authSvc.Signup(ctx, "user1@example.com", "password123", "user1hnd")
		_, _ = authSvc.Signup(ctx, "user2@example.com", "password123", "user2hnd")
		user1, _ := userSvc.GetByHandle(ctx, "user1hnd")

		err := userSvc.UpdateProfile(ctx, user1.ID.String(), UpdateProfileInput{Handle: "user2hnd"})
		if !errors.Is(err, ErrHandleTaken) {
			t.Errorf("expected ErrHandleTaken, got %v", err)
		}
	})
}

func TestUserServiceGetByHandle(t *testing.T) {
	ctx := context.Background()
	userSvc, authSvc := newTestUserService(t)
	_, _ = authSvc.Signup(ctx, "user@example.com", "password123", "testuser")

	t.Run("found", func(t *testing.T) {
		user, err := userSvc.GetByHandle(ctx, "testuser")
		if err != nil {
			t.Fatalf("GetByHandle failed: %v", err)
		}
		if user.Name != "testuser" {
			t.Errorf("got handle %q, want %q", user.Name, "testuser")
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := userSvc.GetByHandle(ctx, "nobody")
		if err == nil {
			t.Error("expected error for missing handle, got nil")
		}
	})
}
