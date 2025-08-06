package main

import (
	"cimrique-redis/internal/client"
	"cimrique-redis/internal/config"
	"cimrique-redis/internal/handlers"
	"cimrique-redis/internal/service"

	"github.com/gofiber/fiber/v2"
)

func main() {
	//set up fiber, in memory storage for limiter
	app := fiber.New()

	//get configs
	AppConfig := config.LoadConfig()

	//create dependencies
	//create redis instance, and pass to handler
	client := client.NewWorkerServiceClient(AppConfig.DBParams.Address, AppConfig.DBParams.Password, AppConfig.DBParams.Number, AppConfig.DBParams.Protocol)
	service := service.NewServicesFuncs(client)
	handler := handlers.NewHandler(service)

	// create post request route that calls the handler function
	app.Post("/enqueue", handler.HandleEnqueue)

	app.Listen(AppConfig.ServerParams.ListenPort)
}
