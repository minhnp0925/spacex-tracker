package services

import (
	"context"
	"slices"
	"strings"

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
}

func NewLaunchService(client clients.SpaceXClient) LaunchService {
	return &concreteLaunchService{
		client: client,
	}
}

func (service *concreteLaunchService) GetNext(ctx context.Context) (*models.Launch, error) {
	return service.client.GetNext(ctx)
}

func (service *concreteLaunchService) GetLatest(ctx context.Context) (*models.Launch, error) {
	return service.client.GetLatest(ctx)
}

func (service *concreteLaunchService) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return service.client.GetUpcoming(ctx)
}

func (service *concreteLaunchService) GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error) {
	launches, err := service.client.GetPast(ctx)
	if err != nil {
		return nil, err
	}

	sortOrder = strings.ToLower(sortOrder)

	slices.SortFunc(launches, func(a, b models.Launch) int {
		if sortOrder == "asc" {
			return a.DateUTC.Compare(b.DateUTC)
		} else {
			// other sort params defaults to desc
			return -a.DateUTC.Compare(b.DateUTC)
		}
	})
	
	return launches, nil
}