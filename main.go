package main

import (
	"embed"
	"log"
	"net/http"

	"myapp/config"
	"myapp/handlers"
	"myapp/i18n"
	"myapp/model"
	"myapp/services"
	"myapp/storage"
	"myapp/util"

	"github.com/3-lines-studio/bifrost"
)

//go:embed all:.bifrost
var bifrostFS embed.FS

func main() {
	database := util.Db

	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer sqlDB.Close()

	if err := i18n.Load(); err != nil {
		log.Fatalf("Failed to load translations: %v", err)
	}

	var store storage.Storage
	if config.Env.STORAGE_TYPE == "s3" {
		s, err := storage.NewS3Storage(
			config.Env.S3_ENDPOINT,
			config.Env.S3_REGION,
			config.Env.S3_BUCKET,
			config.Env.S3_KEY_ID,
			config.Env.S3_APPLICATION_KEY,
			config.Env.S3_BASE_URL,
		)
		if err != nil {
			log.Fatalf("Failed to create S3 storage: %v", err)
		}
		store = s
		log.Print("Using S3 as storage")
	} else {
		s, err := storage.NewLocalStorage("./uploads", config.Env.APP_URL+"/uploads")
		if err != nil {
			log.Fatalf("Failed to create local storage: %v", err)
		}
		store = s
	}

	userRepo := model.NewUserRepository(database)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, authService, store)

	userProps := func(req *http.Request) map[string]any {
		if u := authService.GetUserFromRequest(req); u != nil {
			return map[string]any{"email": u.Email, "handle": u.Name}
		}
		return nil
	}

	profileProps := func(p *model.User) map[string]any {
		return map[string]any{
			"handle":      p.Name,
			"displayName": p.DisplayName,
			"bio":         p.Bio,
			"country":     p.Country,
			"email":       p.Email,
			"avatarURL":   p.AvatarURL,
			"socialLinks": map[string]any{
				"instagram": p.SocialLinks.Instagram,
				"facebook":  p.SocialLinks.Facebook,
				"linkedin":  p.SocialLinks.Linkedin,
				"x":         p.SocialLinks.X,
			},
		}
	}

	app := bifrost.New(
		bifrostFS,
		bifrost.Page("/", "./pages/home.tsx", bifrost.WithLoader(
			func(req *http.Request) (map[string]any, error) {
				locale := i18n.DetectLocale(req)
				props := map[string]any{
					"locale": locale,
					"t":      i18n.Translations(locale),
				}
				if u := userProps(req); u != nil {
					props["user"] = u
				}
				return props, nil
			},
		)),
		bifrost.Page("/login", "./pages/login.tsx", bifrost.WithLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			props := map[string]any{
				"locale": locale,
				"t":      i18n.Translations(locale),
			}
			if e := req.URL.Query().Get("error"); e != "" {
				props["error"] = e
			}
			if u := userProps(req); u != nil {
				props["user"] = u
			}
			return props, nil
		})),
		bifrost.Page("/signup", "./pages/signup.tsx", bifrost.WithLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			props := map[string]any{
				"locale": locale,
				"t":      i18n.Translations(locale),
			}
			if e := req.URL.Query().Get("error"); e != "" {
				props["error"] = e
			}
			if u := userProps(req); u != nil {
				props["user"] = u
			}
			return props, nil
		})),
		bifrost.Page("/user/{handle}", "./pages/profile.tsx", bifrost.WithLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			handle := req.PathValue("handle")
			profile, err := userService.GetByHandle(req.Context(), handle)
			if err != nil {
				return nil, err
			}
			currentUser := authService.GetUserFromRequest(req)
			isOwner := currentUser != nil && currentUser.Name == handle
			props := map[string]any{
				"locale":  locale,
				"t":       i18n.Translations(locale),
				"profile": profileProps(profile),
				"isOwner": isOwner,
			}
			if currentUser != nil {
				props["user"] = map[string]any{"email": currentUser.Email, "handle": currentUser.Name}
			}
			return props, nil
		})),
		bifrost.Page("/user/{handle}/edit", "./pages/profile-edit.tsx", bifrost.WithLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			handle := req.PathValue("handle")
			profile, err := userService.GetByHandle(req.Context(), handle)
			if err != nil {
				return nil, err
			}
			props := map[string]any{
				"locale":  locale,
				"t":       i18n.Translations(locale),
				"profile": profileProps(profile),
			}
			if e := req.URL.Query().Get("error"); e != "" {
				props["error"] = e
			}
			if req.URL.Query().Get("success") == "1" {
				props["success"] = true
			}
			if u := userProps(req); u != nil {
				props["user"] = u
			}
			return props, nil
		})),
	)

	defer app.Stop()

	api := http.NewServeMux()

	api.HandleFunc("POST /api/signup", authHandler.Signup())
	api.HandleFunc("POST /api/login", authHandler.Login())
	api.HandleFunc("POST /api/logout", authHandler.Logout)
	api.HandleFunc("POST /api/user/update", userHandler.UpdateProfile())
	api.HandleFunc("POST /api/set-lang", handleSetLang)

	log.Fatal(http.ListenAndServe(":8080", app.Wrap(api)))
}

func handleSetLang(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("lang")
	if lang != "en" && lang != "es" {
		lang = "en"
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "lang",
		Value:    lang,
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "/"
	}
	http.Redirect(w, r, referer, http.StatusSeeOther)
}
