package main

import (
	"log"
	"net/http"
	"spacex-tracker/clients"
	"spacex-tracker/configs"
	"spacex-tracker/handlers"
	"spacex-tracker/services"
	"spacex-tracker/services/cache"

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

	r.GET("/health", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

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