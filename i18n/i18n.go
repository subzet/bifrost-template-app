package i18n

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"
)

//go:embed locales/*.json
var localesFS embed.FS

var translations map[string]map[string]string

var supportedLocales = []string{"en", "es"}

func Load() error {
	translations = make(map[string]map[string]string)
	for _, locale := range supportedLocales {
		data, err := localesFS.ReadFile("locales/" + locale + ".json")
		if err != nil {
			return err
		}
		var t map[string]string
		if err := json.Unmarshal(data, &t); err != nil {
			return err
		}
		translations[locale] = t
	}
	return nil
}

func DetectLocale(r *http.Request) string {
	if cookie, err := r.Cookie("lang"); err == nil {
		for _, l := range supportedLocales {
			if cookie.Value == l {
				return l
			}
		}
	}

	accept := r.Header.Get("Accept-Language")
	for _, part := range strings.Split(accept, ",") {
		lang := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		lang = strings.SplitN(lang, "-", 2)[0]
		for _, l := range supportedLocales {
			if lang == l {
				return l
			}
		}
	}

	return "en"
}

func T(locale, key string) string {
	if t, ok := translations[locale]; ok {
		if val, ok := t[key]; ok {
			return val
		}
	}
	if val, ok := translations["en"][key]; ok {
		return val
	}
	return key
}

func Translations(locale string) map[string]string {
	if t, ok := translations[locale]; ok {
		return t
	}
	return translations["en"]
}
