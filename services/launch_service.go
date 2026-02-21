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

type baseLaunchService struct {
	client clients.SpaceXClient
}

func NewBaseLaunchService(client clients.SpaceXClient) LaunchService {
	return &baseLaunchService{
		client: client,
	}
}

func (s *baseLaunchService) GetNext(ctx context.Context) (*models.Launch, error) {
	return s.client.GetNext(ctx)
}

func (s *baseLaunchService) GetLatest(ctx context.Context) (*models.Launch, error) {
	return s.client.GetLatest(ctx)
}

func (s *baseLaunchService) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return s.client.GetUpcoming(ctx)
}

func (s *baseLaunchService) GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error) {
	launches, err := s.client.GetPast(ctx)
	if err != nil {
		return nil, err
	}

	sortOrder = strings.ToLower(sortOrder)
	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	slices.SortFunc(launches, func(a, b models.Launch) int {
		if sortOrder == "asc" {
			return a.DateUTC.Compare(b.DateUTC)
		}
		return -a.DateUTC.Compare(b.DateUTC)
	})

	return launches, nil
}