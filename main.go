package main

import (
	"log"
	"spacex-tracker/clients"
	"spacex-tracker/configs"
	"spacex-tracker/handlers"
	"spacex-tracker/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configs
	cfg, err := configs.Load()
	if (err != nil) {
		log.Fatal(err)
	}

	client := clients.NewSpaceXClient(cfg)
	service := services.NewLaunchService(client)
	handler := handlers.NewLaunchHandler(service)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		launches := v1.Group("/launches")
		{
			launches.GET("/next", handler.GetNext)
            launches.GET("/latest", handler.GetLatest)
            launches.GET("/upcoming", handler.GetUpcoming)
            launches.GET("/past", handler.GetPast)
		}
	}

	r.Run(":8080")
}