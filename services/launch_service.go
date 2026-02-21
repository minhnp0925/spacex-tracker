package services

import (
	"context"
	"encoding/json"
	"log"
	"slices"
	"strings"

	"spacex-tracker/cache"
	"spacex-tracker/clients"
	"spacex-tracker/models"
)

type LaunchService interface {
	GetNext(ctx context.Context) (*models.Launch, error)
	GetLatest(ctx context.Context) (*models.Launch, error)
	GetUpcoming(ctx context.Context) ([]models.Launch, error)
	GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error)
}

type concreteLaunchService struct {
	client clients.SpaceXClient
	cache cache.Cache
}

func NewLaunchService(client clients.SpaceXClient, cache cache.Cache) LaunchService {
	return &concreteLaunchService{
		client: client,
		cache: cache,
	}
}

//TODO: Refactor the copy-and-pasted cache logic
func (service *concreteLaunchService) GetNext(ctx context.Context) (*models.Launch, error) {
	cacheKey := "launch:next"

	if service.cache != nil {
		data, err := service.cache.Get(ctx, cacheKey)
		// cache hit
		if err != nil {
			log.Print(err)
		}
		if err == nil {
			var launch models.Launch
			if err := json.Unmarshal(data, &launch); err == nil {
				return &launch, nil
			}
		}
	}

	// calls external
	launch, err := service.client.GetNext(ctx)
	if err != nil {
		return nil, err
	}

	// write into cache
	if service.cache != nil {
		bytes, _ := json.Marshal(launch)
		//TODO: Add a TTL config with ParseTime instead of "60"
		service.cache.Set(ctx, cacheKey, bytes, 60)
	}

	return launch, nil
}

func (service *concreteLaunchService) GetLatest(ctx context.Context) (*models.Launch, error) {
	cacheKey := "launch:latest"

	if service.cache != nil {
		data, err := service.cache.Get(ctx, cacheKey)
		// cache hit
		if err == nil {
			var launch models.Launch
			if err := json.Unmarshal(data, &launch); err == nil {
				return &launch, nil
			}
		}
	}

	// calls external
	launch, err := service.client.GetLatest(ctx)
	if err != nil {
		return nil, err
	}

	// write into cache
	if service.cache != nil {
		bytes, _ := json.Marshal(launch)
		//TODO: Add a TTL config with ParseTime instead of "60"
		service.cache.Set(ctx, cacheKey, bytes, 60)
	}

	return launch, nil
}

func (service *concreteLaunchService) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	cacheKey := "launch:upcoming"

	if service.cache != nil {
		data, err := service.cache.Get(ctx, cacheKey)
		// cache hit
		if err == nil {
			var launches []models.Launch
			if err := json.Unmarshal(data, &launches); err == nil {
				return launches, nil
			}
		}
	}

	// calls external
	launches, err := service.client.GetUpcoming(ctx)
	if err != nil {
		return nil, err
	}

	// write into cache
	if service.cache != nil {
		bytes, _ := json.Marshal(launches)
		//TODO: Add a TTL config with ParseTime instead of "60"
		service.cache.Set(ctx, cacheKey, bytes, 60)
	}

	return launches, nil
}

func (service *concreteLaunchService) GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error) {
	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" {
		sortOrder = "desc"
	}
	
	cacheKey := "launch:past:"+sortOrder

	if service.cache != nil {
		data, err := service.cache.Get(ctx, cacheKey)
		// cache hit
		if err == nil {
			var launches []models.Launch
			if err := json.Unmarshal(data, &launches); err == nil {
				return launches, nil
			}
		}
	}

	// calls external
	launches, err := service.client.GetPast(ctx)
	if err != nil {
		return nil, err
	}
	sortOrder = strings.ToLower(sortOrder)

	slices.SortFunc(launches, func(a, b models.Launch) int {
		if sortOrder == "asc" {
			return a.DateUTC.Compare(b.DateUTC)
		} else { // sortOrder == "desc"
			return -a.DateUTC.Compare(b.DateUTC)
		}
	})

	// write into cache
	if service.cache != nil {
		bytes, _ := json.Marshal(launches)
		//TODO: Add a TTL config with ParseTime instead of "60"
		service.cache.Set(ctx, cacheKey, bytes, 60)
	}

	return launches, nil
}