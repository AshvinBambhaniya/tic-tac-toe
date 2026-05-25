package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// AllConfig variable of type AppConfig
var AllConfig AppConfig

// AppConfig type AppConfig
type AppConfig struct {
	IsDevelopment bool   `envconfig:"IS_DEVELOPMENT" default:"false"`
	Debug         bool   `envconfig:"DEBUG" default:"false"`
	Env           string `envconfig:"APP_ENV" default:"local"`
	Port          string `envconfig:"APP_PORT" default:":3000"`
	Secret        string `envconfig:"JWT_SECRET" default:"ThisIsKey"`
	FrontendURL   string `envconfig:"FRONTEND_URL" default:"http://localhost:3000"`
	DB            DBConfig
}

// GetConfig Collects all configs
func GetConfig() AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("warning .env file not found, scanning from OS ENV")
	}

	AllConfig = AppConfig{}

	err = envconfig.Process("", &AllConfig)
	if err != nil {
		log.Fatal(err)
	}

	return AllConfig
}

// GetConfigByName Collects all configs
func GetConfigByName(key string) string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv(key)
}

// LoadTestEnv loads environment variables from .env.testing file
func LoadTestEnv() AppConfig {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = godotenv.Load(fmt.Sprintf("%s/.env.testing", cwd))
	if err != nil {
		log.Fatal(err)
	}
	return GetConfig()
}
