package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	R2AccountID       string
	R2AccessKeyID     string
	R2SecretAccessKey string
	R2BucketName      string
	Port              string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		R2AccountID:       getEnv("R2_ACCOUNT_ID", "default_account_id"),
		R2AccessKeyID:     getEnv("R2_ACCESS_KEY_ID", "default_access_key"),
		R2SecretAccessKey: getEnv("R2_SECRET_ACCESS_KEY", "default_secret_key"),
		R2BucketName:      getEnv("R2_BUCKET_NAME", "qlove-bucket"),
		Port:              getEnv("PORT", "3000"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
