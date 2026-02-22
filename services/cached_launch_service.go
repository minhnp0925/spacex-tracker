package services

import (
	"context"
	"encoding/json"
	"time"

	"spacex-tracker/services/cache"
	"spacex-tracker/models"
)

type cachedLaunchService struct {
	inner  LaunchService
	cache cache.Cache
	ttl   time.Duration
}

func NewCachedLaunchService(
	inner LaunchService,
	cache cache.Cache,
	ttl time.Duration,
) LaunchService {
	if cache == nil {
        panic("cache cannot be nil")
    }
	return &cachedLaunchService{
		inner: inner,
		cache: cache,
		ttl:   ttl,
	}
}

func getOrSet[T any](
    ctx context.Context, 
    c cache.Cache, 
    key string, 
    ttl time.Duration, 
    fetch func(context.Context) (T, error),
) (T, error) {
    // Attempt to retrieve from cache
    if data, err := c.Get(ctx, key); err == nil {
        var result T
        if err := json.Unmarshal(data, &result); err == nil {
            return result, nil
        }
    }

    // Cache miss: Fetch data via the provided function
    result, err := fetch(ctx)
    if err != nil {
        var zero T
        return zero, err
    }

    // Store in cache for future use
    if bytes, err := json.Marshal(result); err == nil {
        _ = c.Set(ctx, key, bytes, ttl)
    }

    return result, nil
}

func (c *cachedLaunchService) GetNext(ctx context.Context) (*models.Launch, error) {
	return getOrSet(ctx, c.cache, "launch:next", c.ttl, c.inner.GetNext)
}

func (c *cachedLaunchService) GetLatest(ctx context.Context) (*models.Launch, error) {
	return getOrSet(ctx, c.cache, "launch:latest", c.ttl, c.inner.GetLatest)
}

func (c *cachedLaunchService) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return getOrSet(ctx, c.cache, "launch:upcoming", c.ttl, c.inner.GetUpcoming)
}

func (c *cachedLaunchService) GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error) {
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	
	key := "launch:past"+sortOrder
	return getOrSet(ctx, c.cache, key, c.ttl, c.inner.GetUpcoming)
}