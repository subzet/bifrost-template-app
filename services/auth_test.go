package services

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"myapp/model"
	"myapp/testutil"
)

func newTestService(t *testing.T) *AuthService {
	t.Helper()
	return NewAuthService(model.NewUserRepository(testutil.NewTestDB(t, &model.User{})))
}

func TestSignup(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns token", func(t *testing.T) {
		svc := newTestService(t)
		token, err := svc.Signup(ctx, "user@example.com", "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected a token, got empty string")
		}
	})

	t.Run("email taken", func(t *testing.T) {
		svc := newTestService(t)
		_, _ = svc.Signup(ctx, "user@example.com", "password123")
		_, err := svc.Signup(ctx, "user@example.com", "different123")
		if !errors.Is(err, ErrEmailTaken) {
			t.Errorf("expected ErrEmailTaken, got %v", err)
		}
	})
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	t.Run("success returns token", func(t *testing.T) {
		svc := newTestService(t)
		_, _ = svc.Signup(ctx, "user@example.com", "password123")

		token, err := svc.Login(ctx, "user@example.com", "password123")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if token == "" {
			t.Error("expected a token, got empty string")
		}
	})

	t.Run("wrong password", func(t *testing.T) {
		svc := newTestService(t)
		_, _ = svc.Signup(ctx, "user@example.com", "password123")

		_, err := svc.Login(ctx, "user@example.com", "wrongpassword")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		svc := newTestService(t)
		_, err := svc.Login(ctx, "nobody@example.com", "password123")
		if !errors.Is(err, ErrInvalidCredentials) {
			t.Errorf("expected ErrInvalidCredentials, got %v", err)
		}
	})
}

func TestGetUserFromRequest(t *testing.T) {
	ctx := context.Background()
	svc := newTestService(t)
	token, _ := svc.Signup(ctx, "user@example.com", "password123")

	t.Run("valid token returns user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: token})

		user := svc.GetUserFromRequest(req)
		if user == nil {
			t.Fatal("expected user, got nil")
		}
		if user.Email != "user@example.com" {
			t.Errorf("expected email user@example.com, got %s", user.Email)
		}
	})

	t.Run("invalid token returns nil", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.AddCookie(&http.Cookie{Name: "session", Value: "not.a.valid.token"})

		if svc.GetUserFromRequest(req) != nil {
			t.Error("expected nil for invalid token")
		}
	})

	t.Run("missing cookie returns nil", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if svc.GetUserFromRequest(req) != nil {
			t.Error("expected nil for missing cookie")
		}
	})
}
