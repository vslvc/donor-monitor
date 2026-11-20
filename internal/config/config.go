package config

import (
	"os"
	"fmt"

	"github.com/joho/godotenv"
)

type Config struct{
	DatabaseURL string
}

func Load() (Config, error){
	err := godotenv.Load()

	if err != nil && !os.IsNotExist(err){
		return Config{}, fmt.Errorf("load .env: %w", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == ""{
		return Config{}, fmt.Errorf("DATABASE_URL is not ser")
	}

	return Config{
		DatabaseURL: databaseURL,
	}, nil
	
}
