package main

import (
	handler "learning_api/internal/gateway/http"
	"learning_api/internal/repository"
	"learning_api/internal/service"
	todo "learning_api/pkg"
	"log"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(8000, handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}
}
