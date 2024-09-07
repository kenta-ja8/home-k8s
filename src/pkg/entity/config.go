package entity

import "os"

type Config struct {
	SAMPLE_USER       string
	BUILD_DATE        string
	HOSTNAME          string
	POSTGRES_HOST     string
	POSTGRES_DB       string
	POSTGRES_USER     string
	POSTGRES_PASSWORD string
}

func LoadConfig() *Config {
	return &Config{
		SAMPLE_USER:       os.Getenv("SAMPLE_USER"),
		BUILD_DATE:        os.Getenv("BUILD_DATE"),
		HOSTNAME:          os.Getenv("HOSTNAME"),
		POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
		POSTGRES_DB:       os.Getenv("POSTGRES_DB"),
		POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
		POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
	}
}
