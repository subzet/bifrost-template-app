package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"myapp/model"
	"myapp/services"
	"myapp/testutil"
)

func newTestHandler(t *testing.T) *AuthHandler {
	t.Helper()
	repo := model.NewUserRepository(testutil.NewTestDB(t, &model.User{}))
	return NewAuthHandler(services.NewAuthService(repo))
}

func postForm(handler http.HandlerFunc, target string, values url.Values) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func sessionCookie(w *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range w.Result().Cookies() {
		if c.Name == "session" {
			return c
		}
	}
	return nil
}

func TestHandlerSignup(t *testing.T) {
	t.Run("empty fields redirect to signup with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {""}, "password": {""}, "confirm_password": {""}, "handle": {""},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("missing handle redirects to signup with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {""},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("invalid handle redirects to signup with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"ab"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("passwords mismatch redirect to signup with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"different123"}, "handle": {"testuser"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("password too short redirect to signup with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"short"}, "confirm_password": {"short"}, "handle": {"testuser"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("success sets session cookie and redirects to edit profile", func(t *testing.T) {
		w := postForm(newTestHandler(t).Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"testuser"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if w.Header().Get("Location") != "/user/testuser/edit" {
			t.Errorf("expected redirect to /user/testuser/edit, got %s", w.Header().Get("Location"))
		}
		if c := sessionCookie(w); c == nil || c.Value == "" {
			t.Error("expected session cookie to be set")
		}
	})

	t.Run("duplicate email redirect to signup with error", func(t *testing.T) {
		h := newTestHandler(t)
		form := url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"user1hnd"},
		}
		postForm(h.Signup(), "/api/signup", form)
		w := postForm(h.Signup(), "/api/signup", url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"user2hnd"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})

	t.Run("duplicate handle redirect to signup with error", func(t *testing.T) {
		h := newTestHandler(t)
		form := url.Values{
			"email": {"user1@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"testuser"},
		}
		postForm(h.Signup(), "/api/signup", form)
		w := postForm(h.Signup(), "/api/signup", url.Values{
			"email": {"user2@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"testuser"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/signup?error=") {
			t.Errorf("expected /signup?error=..., got %s", loc)
		}
	})
}

func TestHandlerLogin(t *testing.T) {
	t.Run("empty fields redirect to login with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Login(), "/api/login", url.Values{
			"email": {""}, "password": {""},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/login?error=") {
			t.Errorf("expected /login?error=..., got %s", loc)
		}
	})

	t.Run("invalid credentials redirect to login with error", func(t *testing.T) {
		w := postForm(newTestHandler(t).Login(), "/api/login", url.Values{
			"email": {"nobody@example.com"}, "password": {"password123"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if loc := w.Header().Get("Location"); !strings.HasPrefix(loc, "/login?error=") {
			t.Errorf("expected /login?error=..., got %s", loc)
		}
	})

	t.Run("success sets session cookie and redirects home", func(t *testing.T) {
		h := newTestHandler(t)
		creds := url.Values{
			"email": {"user@example.com"}, "password": {"password123"}, "confirm_password": {"password123"}, "handle": {"testuser"},
		}
		postForm(h.Signup(), "/api/signup", creds)

		w := postForm(h.Login(), "/api/login", url.Values{
			"email": {"user@example.com"}, "password": {"password123"},
		})
		if w.Code != http.StatusSeeOther {
			t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
		}
		if w.Header().Get("Location") != "/" {
			t.Errorf("expected redirect to /, got %s", w.Header().Get("Location"))
		}
		if c := sessionCookie(w); c == nil || c.Value == "" {
			t.Error("expected session cookie to be set")
		}
	})
}

func TestHandlerLogout(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/logout", nil)
	w := httptest.NewRecorder()
	newTestHandler(t).Logout(w, req)

	if w.Code != http.StatusSeeOther {
		t.Errorf("expected %d, got %d", http.StatusSeeOther, w.Code)
	}
	if w.Header().Get("Location") != "/" {
		t.Errorf("expected redirect to /, got %s", w.Header().Get("Location"))
	}
	if c := sessionCookie(w); c == nil || c.MaxAge != -1 {
		t.Error("expected session cookie to be cleared (MaxAge=-1)")
	}
}
