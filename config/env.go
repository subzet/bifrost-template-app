package config

import "os"

type authEnv struct {
	JWT_SECRET string
	DB_DSN     string

	// Storage: "local" (default) or "s3"
	STORAGE_TYPE string
	// APP_URL is used to build public URLs for local storage (e.g. http://localhost:8080)
	APP_URL string

	// Backblaze B2 / S3-compatible storage
	S3_ENDPOINT        string
	S3_BUCKET          string
	S3_KEY_ID          string
	S3_APPLICATION_KEY string
	S3_REGION          string
	// S3_BASE_URL is the public base URL for uploaded files
	S3_BASE_URL string
}

var Env authEnv = authEnv{
	JWT_SECRET: getenvDefault("JWT_SECRET", "dev-secret-change-me"),
	DB_DSN:     getenvDefault("DB_DSN", "file:dev.db"),

	STORAGE_TYPE: getenvDefault("STORAGE_TYPE", "local"),
	APP_URL:      getenvDefault("APP_URL", "http://localhost:8080"),

	S3_ENDPOINT:        os.Getenv("S3_ENDPOINT"),
	S3_BUCKET:          os.Getenv("S3_BUCKET"),
	S3_KEY_ID:          os.Getenv("S3_KEY_ID"),
	S3_APPLICATION_KEY: os.Getenv("S3_APPLICATION_KEY"),
	S3_REGION:          getenvDefault("S3_REGION", "us-west-004"),
	S3_BASE_URL:        os.Getenv("S3_BASE_URL"),
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
