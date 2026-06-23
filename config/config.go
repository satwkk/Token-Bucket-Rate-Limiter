package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Addr       string
	BackendUrl string
	MaxToken   float64
	RefillRate float64
}

func LoadEnv() bool {
	if err := godotenv.Load(); err != nil {
		return false
	}
	return true
}

func LoadProxyConfig() *Config {
	maxToken, err := strconv.ParseFloat(os.Getenv("MAX_TOKEN"), 64)
	if err != nil {
		log.Fatal("Error parsing MAX_TOKEN")
		return nil
	}
	refillRate, err := strconv.ParseFloat(os.Getenv("REFILL_RATE"), 64)
	if err != nil {
		log.Fatal("Error parsing REFILL_RATE")
		return nil
	}
	return &Config{
		Addr:       os.Getenv("PROXY_ADDR"),
		BackendUrl: os.Getenv("BACKEND_URL"),
		MaxToken:   maxToken,
		RefillRate: refillRate,
	}
}
