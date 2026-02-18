package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"myapp/i18n"
	"myapp/services"
)

type AuthHandler struct {
	svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := i18n.DetectLocale(r)
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/signup?error="+url.QueryEscape(i18n.T(locale, "error.emailPasswordRequired")), http.StatusSeeOther)
			return
		}
		if password != confirmPassword {
			http.Redirect(w, r, "/signup?error="+url.QueryEscape(i18n.T(locale, "error.passwordsMismatch")), http.StatusSeeOther)
			return
		}
		if len(password) < 8 {
			http.Redirect(w, r, "/signup?error="+url.QueryEscape(i18n.T(locale, "error.passwordTooShort")), http.StatusSeeOther)
			return
		}

		token, err := h.svc.Signup(r.Context(), email, password)
		if err != nil {
			errKey := "error.somethingWrong"
			if errors.Is(err, services.ErrEmailTaken) {
				errKey = "error.emailTaken"
			}
			http.Redirect(w, r, "/signup?error="+url.QueryEscape(i18n.T(locale, errKey)), http.StatusSeeOther)
			return
		}

		setSessionCookie(w, token)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := i18n.DetectLocale(r)
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if email == "" || password == "" {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(i18n.T(locale, "error.emailPasswordRequired")), http.StatusSeeOther)
			return
		}

		token, err := h.svc.Login(r.Context(), email, password)
		if err != nil {
			http.Redirect(w, r, "/login?error="+url.QueryEscape(i18n.T(locale, "error.invalidCredentials")), http.StatusSeeOther)
			return
		}

		setSessionCookie(w, token)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	clearSessionCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
