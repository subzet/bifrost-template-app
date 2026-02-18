package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"myapp/model"
	"myapp/services"
	"myapp/storage"
	"myapp/testutil"
)

func newTestUserHandler(t *testing.T) (*UserHandler, *services.AuthService) {
	t.Helper()
	repo := model.NewUserRepository(testutil.NewTestDB(t, &model.User{}))
	authSvc := services.NewAuthService(repo)
	userSvc := services.NewUserService(repo)
	return NewUserHandler(userSvc, authSvc, storage.Noop()), authSvc
}

func mockPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestHandlerProfile(t *testing.T) {
	ctx := context.Background()
	h, authSvc := newTestUserHandler(t)
	_, _ = authSvc.Signup(ctx, "test@example.com", "password123", "testuser")

	t.Run("200 for known handle", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/testuser", nil)
		req.SetPathValue("handle", "testuser")
		w := httptest.NewRecorder()
		h.ServeProfile(mockPage())(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})

	t.Run("404 for unknown handle", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/user/unknown", nil)
		req.SetPathValue("handle", "unknown")
		w := httptest.NewRecorder()
		h.ServeProfile(mockPage())(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestHandlerEditProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("redirect to login if not authenticated", func(t *testing.T) {
		h, _ := newTestUserHandler(t)
		req := httptest.NewRequest(http.MethodGet, "/user/testuser/edit", nil)
		req.SetPathValue("handle", "testuser")
		w := httptest.NewRecorder()
		h.ServeEditProfile(mockPage())(w, req)
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "/login" {
			t.Errorf("expected redirect to /login, got %s", loc)
		}
	})

	t.Run("redirect to view if not owner", func(t *testing.T) {
		h, authSvc := newTestUserHandler(t)
		_, _ = authSvc.Signup(ctx, "user1@example.com", "password123", "user1hnd")
		token2, _ := authSvc.Signup(ctx, "user2@example.com", "password123", "user2hnd")

		req := httptest.NewRequest(http.MethodGet, "/user/user1hnd/edit", nil)
		req.SetPathValue("handle", "user1hnd")
		req.AddCookie(&http.Cookie{Name: "session", Value: token2})
		w := httptest.NewRecorder()
		h.ServeEditProfile(mockPage())(w, req)
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "/user/user1hnd" {
			t.Errorf("expected redirect to /user/user1hnd, got %s", loc)
		}
	})

	t.Run("200 if owner", func(t *testing.T) {
		h, authSvc := newTestUserHandler(t)
		token, _ := authSvc.Signup(ctx, "user@example.com", "password123", "testuser")

		req := httptest.NewRequest(http.MethodGet, "/user/testuser/edit", nil)
		req.SetPathValue("handle", "testuser")
		req.AddCookie(&http.Cookie{Name: "session", Value: token})
		w := httptest.NewRecorder()
		h.ServeEditProfile(mockPage())(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}
	})
}

func TestHandlerUpdateProfile(t *testing.T) {
	ctx := context.Background()

	t.Run("redirect to login if not authenticated", func(t *testing.T) {
		h, _ := newTestUserHandler(t)
		w := postForm(h.UpdateProfile(), "/api/user/update", url.Values{"handle": {"testuser"}})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "/login" {
			t.Errorf("expected /login, got %s", loc)
		}
	})

	t.Run("invalid handle redirects with error", func(t *testing.T) {
		h, authSvc := newTestUserHandler(t)
		token, _ := authSvc.Signup(ctx, "user@example.com", "password123", "testuser")

		req := httptest.NewRequest(http.MethodPost, "/api/user/update", strings.NewReader(url.Values{"handle": {"ab"}}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{Name: "session", Value: token})
		w := httptest.NewRecorder()
		h.UpdateProfile()(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/user/testuser/edit?error=") {
			t.Errorf("expected /user/testuser/edit?error=..., got %s", loc)
		}
	})

	t.Run("handle taken redirects with error", func(t *testing.T) {
		h, authSvc := newTestUserHandler(t)
		_, _ = authSvc.Signup(ctx, "user1@example.com", "password123", "user1hnd")
		token2, _ := authSvc.Signup(ctx, "user2@example.com", "password123", "user2hnd")

		req := httptest.NewRequest(http.MethodPost, "/api/user/update", strings.NewReader(url.Values{"handle": {"user1hnd"}}.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{Name: "session", Value: token2})
		w := httptest.NewRecorder()
		h.UpdateProfile()(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/user/user2hnd/edit?error=") {
			t.Errorf("expected /user/user2hnd/edit?error=..., got %s", loc)
		}
	})

	t.Run("success redirects to edit with success flag", func(t *testing.T) {
		h, authSvc := newTestUserHandler(t)
		token, _ := authSvc.Signup(ctx, "user@example.com", "password123", "testuser")

		form := url.Values{
			"handle":       {"newhandle"},
			"display_name": {"Test User"},
			"bio":          {"My bio"},
			"country":      {"US"},
		}
		req := httptest.NewRequest(http.MethodPost, "/api/user/update", strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.AddCookie(&http.Cookie{Name: "session", Value: token})
		w := httptest.NewRecorder()
		h.UpdateProfile()(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); loc != "/user/newhandle/edit?success=1" {
			t.Errorf("expected /user/newhandle/edit?success=1, got %s", loc)
		}
	})
}
