package services

import (
	"context"
	"spacex-tracker/clients"
	"spacex-tracker/models"
)

type LaunchService interface {
	GetNext(ctx context.Context) (*models.Launch, error)
	GetLatest(ctx context.Context) (*models.Launch, error)
	GetUpcoming(ctx context.Context) ([]models.Launch, error)
	GetPast(ctx context.Context) ([]models.Launch, error)
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

func (service *concreteLaunchService) GetPast(ctx context.Context) ([]models.Launch, error) {
	return service.client.GetPast(ctx)
}