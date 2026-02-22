package main

import (
	"context"
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

	log.Println("Using Redis URL:", cfg.RedisURL)
	
	var rdb *redis.Client
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err == nil {
		candidate := redis.NewClient(opt)

		if err := candidate.Ping(context.Background()).Err(); err != nil {
			log.Printf("Redis unreachable: %v. Running cacheless.", err)
		} else {
			rdb = candidate
			defer rdb.Close()
			log.Println("Redis connected successfully")
		}
	} else {
		log.Println("Invalid Redis URL, running cacheless")
	}
	
	client := clients.NewSpaceXClient(cfg)
	base := services.NewBaseLaunchService(client)
	var service services.LaunchService
	
	if rdb != nil {
		redisCache := cache.NewRedisCache(rdb)
		service = services.NewCachedLaunchService(base, redisCache, cfg.CacheTTL)
	} else {
		service = base
	}

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