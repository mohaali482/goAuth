package main

import (
	"log"

	"github.com/mohaali/goAuth/auth"
	authGorm "github.com/mohaali/goAuth/auth/gorm"
	"github.com/mohaali/goAuth/config"
	"github.com/mohaali/goAuth/internal/api"
	"github.com/mohaali/goAuth/internal/http/gin"
)

func main() {
	appConfig, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	r, err := authGorm.NewGormRepository(appConfig.DB)
	if err != nil {
		panic(err)
	}
	s := auth.NewUserService(r, appConfig)
	h := gin.Handlers(*s)

	err = api.Start(appConfig.Port, h)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
