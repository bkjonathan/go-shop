package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	AWS      AWSConfig
	Upload   UploadConfig
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret              string
	ExpiresIn           time.Duration
	RefreshTokenExpires time.Duration
}

type AWSConfig struct {
	Region          string
	AccessKeyId     string
	SecretAccessKey string
	S3Bucket        string
	S3Endpoint      string
	// S3PublicEndpoint is the address browsers fetch objects from - the CDN
	// container rather than the endpoint the SDK uploads through. Empty falls
	// back to S3Endpoint.
	S3PublicEndpoint string
}

type UploadConfig struct {
	Path           string
	MaxFileSize    int64
	UploadProvider string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	jwtExpiresIn, _ := time.ParseDuration(getEnv("JWT_EXPIRES_IN", "24h"))
	refreshTokenExpires, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRES_IN", "720h"))
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)

	return &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8090"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "ecommerce_shop"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:              getEnv("JWT_SECRET", "you_jwt_secret_key"),
			ExpiresIn:           jwtExpiresIn,
			RefreshTokenExpires: refreshTokenExpires,
		},
		AWS: AWSConfig{
			Region:           getEnv("AWS_REGION", "us-east-1"),
			AccessKeyId:      getEnv("AWS_ACCESS_KEY", "test"),
			SecretAccessKey:  getEnv("AWS_SECRET_ACCESS_KEY", "test"),
			S3Bucket:         getEnv("AWS_S3_BUCKET", "ecommerce-uploads"),
			S3Endpoint:       getEnv("AWS_S3_ENDPOINT", "http://localhost:4566"),
			S3PublicEndpoint: getEnv("AWS_S3_PUBLIC_ENDPOINT", ""),
		},
		Upload: UploadConfig{
			Path:           getEnv("UPLOAD_PATH", "./uploads"),
			MaxFileSize:    maxUploadSize,
			UploadProvider: getEnv("UPLOAD_PROVIDER", "local"),
		},
	}, nil

}
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
