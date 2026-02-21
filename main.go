package main

import (
	"log"
	"spacex-tracker/services/cache"
	"spacex-tracker/clients"
	"spacex-tracker/configs"
	"spacex-tracker/handlers"
	"spacex-tracker/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configs
	cfg, err := configs.Load()
	if (err != nil) {
		log.Fatal(err)
	}

	log.Println("Using Redis address:", cfg.RedisAddress)
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddress,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	// Ensure that the connection is properly closed gracefully
	defer rdb.Close()

	redisCache := cache.NewRedisCache(rdb)

	client := clients.NewSpaceXClient(cfg)
	base := services.NewBaseLaunchService(client)
	service := services.NewCachedLaunchService(base, redisCache, cfg.CacheTTL)
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