package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DB             string
	Secret         string
	AccessExpTime  int
	RefreshExpTime int
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	accessExpTime, err := strconv.Atoi(os.Getenv("AccessExpTime"))
	if err != nil {
		panic(err)
	}

	refreshExpTime, err := strconv.Atoi(os.Getenv("RefreshExpTime"))
	if err != nil {
		panic(err)
	}

	config := &Config{
		Port:           os.Getenv("PORT"),
		DB:             os.Getenv("DB"),
		Secret:         os.Getenv("SECRET"),
		AccessExpTime:  accessExpTime,
		RefreshExpTime: refreshExpTime,
	}

	return config, nil
}
