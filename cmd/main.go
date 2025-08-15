package main

import (
	"cimrique-redis/internal/client"
	"cimrique-redis/internal/config"
	"cimrique-redis/internal/handlers"
	metric "cimrique-redis/internal/metrics"
	"cimrique-redis/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	//set up fiber, in memory storage for limiter
	app := fiber.New()

	//get configs
	AppConfig := config.LoadConfig()

	//metrics
	reg := prometheus.NewRegistry()
	metric := metric.NewMetric(reg)
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	//create dependencies
	//create redis instance, and pass to handler
	client := client.NewWorkerServiceClient(AppConfig.DBParams.Address, AppConfig.DBParams.Password, AppConfig.DBParams.Number, AppConfig.DBParams.Protocol)
	service := service.NewServicesFuncs(client, metric)
	handler := handlers.NewHandler(service)

	// create post request route that calls the handler function
	app.Post("/enqueue", handler.HandleEnqueue)

	// expose /metrics on the same Fiber app/port
	app.Get("/metrics", adaptor.HTTPHandler(promHandler))

	app.Listen(AppConfig.ServerParams.ListenPort)
}
