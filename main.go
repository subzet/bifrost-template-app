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

	r, err := bifrost.New(bifrost.WithAssetsFS(bifrostFS))
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	defer r.Stop()

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
	} else {
		s, err := storage.NewLocalStorage("./uploads", config.Env.APP_URL+"/uploads")
		if err != nil {
			log.Fatalf("Failed to create local storage: %v", err)
		}
		store = s
	}

	userRepo := model.NewUserRepository(database)
	authSvc := services.NewAuthService(userRepo)
	userSvc := services.NewUserService(userRepo)
	authHandler := handlers.NewAuthHandler(authSvc)
	userHandler := handlers.NewUserHandler(userSvc, authSvc, store)

	userProps := func(req *http.Request) map[string]any {
		if u := authSvc.GetUserFromRequest(req); u != nil {
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

	home := r.NewPage("./pages/home.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			props := map[string]any{
				"locale": locale,
				"t":      i18n.Translations(locale),
			}
			if u := userProps(req); u != nil {
				props["user"] = u
			}
			return props, nil
		}),
	)

	login := r.NewPage("./pages/login.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
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
		}),
	)

	signup := r.NewPage("./pages/signup.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
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
		}),
	)

	profilePage := r.NewPage("./pages/profile.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			handle := req.PathValue("handle")
			profile, err := userSvc.GetByHandle(req.Context(), handle)
			if err != nil {
				return nil, err
			}
			currentUser := authSvc.GetUserFromRequest(req)
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
		}),
	)

	editPage := r.NewPage("./pages/profile-edit.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			handle := req.PathValue("handle")
			profile, err := userSvc.GetByHandle(req.Context(), handle)
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
		}),
	)

	router := http.NewServeMux()
	router.Handle("GET /{$}", home)
	router.Handle("GET /login", login)
	router.Handle("GET /signup", signup)
	router.Handle("GET /user/{handle}", userHandler.ServeProfile(profilePage))
	router.Handle("GET /user/{handle}/edit", userHandler.ServeEditProfile(editPage))
	// Serve locally uploaded files
	router.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	router.HandleFunc("POST /api/signup", authHandler.Signup())
	router.HandleFunc("POST /api/login", authHandler.Login())
	router.HandleFunc("POST /api/logout", authHandler.Logout)
	router.HandleFunc("POST /api/user/update", userHandler.UpdateProfile())
	router.HandleFunc("POST /api/set-lang", handleSetLang)

	assetRouter := http.NewServeMux()
	bifrost.RegisterAssetRoutes(assetRouter, r, router)

	addr := ":8080"
	log.Printf("Server running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, assetRouter); err != nil {
		log.Fatal(err)
	}
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
