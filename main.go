package main

import (
	"embed"
	"log"
	"net/http"

	"myapp/handlers"
	"myapp/i18n"
	"myapp/model"
	"myapp/services"
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

	userRepo := model.NewUserRepository(database)
	authSvc := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authSvc)

	home := r.NewPage("./pages/home.tsx",
		bifrost.WithPropsLoader(func(req *http.Request) (map[string]any, error) {
			locale := i18n.DetectLocale(req)
			props := map[string]any{
				"locale": locale,
				"t":      i18n.Translations(locale),
			}
			if u := authSvc.GetUserFromRequest(req); u != nil {
				props["user"] = map[string]any{"email": u.Email}
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
			if u := authSvc.GetUserFromRequest(req); u != nil {
				props["user"] = map[string]any{"email": u.Email}
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
			if u := authSvc.GetUserFromRequest(req); u != nil {
				props["user"] = map[string]any{"email": u.Email}
			}
			return props, nil
		}),
	)

	router := http.NewServeMux()
	router.Handle("GET /{$}", home)
	router.Handle("GET /login", login)
	router.Handle("GET /signup", signup)

	router.HandleFunc("POST /api/signup", authHandler.Signup())
	router.HandleFunc("POST /api/login", authHandler.Login())
	router.HandleFunc("POST /api/logout", authHandler.Logout)
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
