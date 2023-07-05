package main

import (
	"log"

	"github.com/mohaali482/goAuth/auth"
	authGorm "github.com/mohaali482/goAuth/auth/gorm"
	"github.com/mohaali482/goAuth/config"
	fiberHandler "github.com/mohaali482/goAuth/internal/http/fiber"
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
	// h := gin.Handlers(*s)
	app := fiberHandler.App(*s)

	// err = api.Start(appConfig.Port, h)
	err = app.Listen(appConfig.Port)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
