package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
	handler "learning_api/internal/gateway/http"
	"learning_api/internal/repository"
	"learning_api/internal/service"
	todo "learning_api/pkg"
	"os"
	"strconv"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := gotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables %s", err.Error())
	}

	port, err := strconv.Atoi(os.Getenv("DB_PG_PORT"))
	if err != nil {
		logrus.Fatalf("postgres database port must be integer value: %s", err.Error())
	}

	db, err := repository.NewPostgresDb(repository.Config{
		Host:     os.Getenv("DB_PG_HOST"),
		Port:     port,
		Username: os.Getenv("DB_PG_USER"),
		Password: os.Getenv("DB_PG_PASSWORD"),
		DBName:   os.Getenv("DB_PG_DATABASE"),
		SSLMode:  os.Getenv("DB_PG_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)

	port, err = strconv.Atoi(os.Getenv("APP_PORT"))
	if err != nil {
		logrus.Fatalf("application port must be integer value: %s", err.Error())
	}

	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occured while running http server: %s", err.Error())
	}
}
