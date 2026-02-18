package config

import "os"

type authEnv struct {
	JWT_SECRET string
	DB_DSN     string
}

var Env authEnv = authEnv{
	JWT_SECRET: getenvDefault("JWT_SECRET", "dev-secret-change-me"),
	DB_DSN:     getenvDefault("DB_DSN", "file:dev.db"),
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
