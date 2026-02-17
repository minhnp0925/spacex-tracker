package main

import (
	"github.com/gin-gonic/gin"
	"spacex-tracker/clients"
	"spacex-tracker/handlers"
	"spacex-tracker/services"
)

func main() {
	client := clients.NewSpaceXClient()
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