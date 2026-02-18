package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"myapp/i18n"
	"myapp/model"
	"myapp/services"
	"myapp/storage"
)

type UserHandler struct {
	userSvc *services.UserService
	authSvc *services.AuthService
	store   storage.Storage
}

func NewUserHandler(userSvc *services.UserService, authSvc *services.AuthService, store storage.Storage) *UserHandler {
	return &UserHandler{userSvc: userSvc, authSvc: authSvc, store: store}
}

func (h *UserHandler) ServeProfile(page http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle := r.PathValue("handle")
		if _, err := h.userSvc.GetByHandle(r.Context(), handle); err != nil {
			http.NotFound(w, r)
			return
		}
		page.ServeHTTP(w, r)
	}
}

func (h *UserHandler) ServeEditProfile(page http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handle := r.PathValue("handle")
		currentUser := h.authSvc.GetUserFromRequest(r)
		if currentUser == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if currentUser.Name != handle {
			http.Redirect(w, r, "/user/"+handle, http.StatusSeeOther)
			return
		}
		page.ServeHTTP(w, r)
	}
}

func (h *UserHandler) UpdateProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		locale := i18n.DetectLocale(r)
		currentUser := h.authSvc.GetUserFromRequest(r)
		if currentUser == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Parse multipart form (supports both file uploads and plain form submissions).
		_ = r.ParseMultipartForm(10 << 20)

		oldHandle := currentUser.Name
		input := services.UpdateProfileInput{
			Handle:      strings.TrimSpace(r.FormValue("handle")),
			DisplayName: strings.TrimSpace(r.FormValue("display_name")),
			Bio:         strings.TrimSpace(r.FormValue("bio")),
			Country:     strings.TrimSpace(r.FormValue("country")),
			SocialLinks: model.SocialLinks{
				Instagram: strings.TrimSpace(r.FormValue("instagram")),
				Facebook:  strings.TrimSpace(r.FormValue("facebook")),
				Linkedin:  strings.TrimSpace(r.FormValue("linkedin")),
				X:         strings.TrimSpace(r.FormValue("x")),
			},
		}

		if file, header, err := r.FormFile("avatar"); err == nil {
			defer file.Close()
			contentType := header.Header.Get("Content-Type")
			if ext := imageExt(contentType); ext != "" {
				key := "avatars/" + currentUser.ID.String() + ext
				if avatarURL, err := h.store.Upload(r.Context(), key, file, header.Size, contentType); err == nil {
					input.AvatarURL = avatarURL
				}
			}
		}

		if err := h.userSvc.UpdateProfile(r.Context(), currentUser.ID.String(), input); err != nil {
			errKey := "error.somethingWrong"
			if errors.Is(err, services.ErrHandleTaken) {
				errKey = "error.handleTaken"
			} else if errors.Is(err, services.ErrHandleInvalid) {
				errKey = "error.handleInvalid"
			}
			http.Redirect(w, r, "/user/"+oldHandle+"/edit?error="+url.QueryEscape(i18n.T(locale, errKey)), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/user/"+input.Handle+"/edit?success=1", http.StatusSeeOther)
	}
}

func imageExt(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ""
	}
}
