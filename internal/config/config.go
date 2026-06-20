package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port     string
	BotToken string
	Dsn      string
}

var GlobalConfig *Config

func LoadConfig() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	GlobalConfig = &Config{
		Port:     os.Getenv("PORT"),
		BotToken: os.Getenv("BOT_TOKEN"),
		Dsn:      os.Getenv("DSN"),
	}
}
