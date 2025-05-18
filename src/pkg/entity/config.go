package entity

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	SAMPLE_USER        string
	BUILD_DATE         string
	HOSTNAME           string
	POSTGRES_HOST      string
	POSTGRES_DB        string
	POSTGRES_USER      string
	POSTGRES_PASSWORD  string
	LINE_CHANNEL_TOKEN string
	LINE_GROUP_ID      string
}

var once sync.Once

func LoadConfig() *Config {

	once.Do(func() {
		err := godotenv.Load(".env")
		if err != nil {
			log.Printf("env file not loaded: %v", err)
		}
	})

	return &Config{
		SAMPLE_USER:        os.Getenv("SAMPLE_USER"),
		BUILD_DATE:         os.Getenv("BUILD_DATE"),
		HOSTNAME:           os.Getenv("HOSTNAME"),
		POSTGRES_HOST:      os.Getenv("POSTGRES_HOST"),
		POSTGRES_DB:        os.Getenv("POSTGRES_DB"),
		POSTGRES_USER:      os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD:  os.Getenv("POSTGRES_PASSWORD"),
		LINE_CHANNEL_TOKEN: os.Getenv("LINE_CHANNEL_TOKEN"),
		LINE_GROUP_ID:      os.Getenv("LINE_GROUP_ID"),
	}
}
