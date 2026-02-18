package services

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"spacex-tracker/models"
)

// type SpaceXClient interface {
// 	GetNext(ctx context.Context) (*models.Launch, error)
// 	GetLatest(ctx context.Context) (*models.Launch, error)
// 	GetUpcoming(ctx context.Context) ([]models.Launch, error)
// 	GetPast(ctx context.Context) ([]models.Launch, error)
// }
type MockSpaceXClient struct {
	GetNextFunc     func(ctx context.Context) (*models.Launch, error)
	GetLatestFunc   func(ctx context.Context) (*models.Launch, error)
	GetUpcomingFunc func(ctx context.Context) ([]models.Launch, error)
	GetPastFunc     func(ctx context.Context) ([]models.Launch, error)
}

func (m *MockSpaceXClient) GetNext(ctx context.Context) (*models.Launch, error) {
	return m.GetNextFunc(ctx)
}

func (m *MockSpaceXClient) GetLatest(ctx context.Context) (*models.Launch, error) {
	return m.GetLatestFunc(ctx)
}

func (m *MockSpaceXClient) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return m.GetUpcomingFunc(ctx)
}

func (m *MockSpaceXClient) GetPast(ctx context.Context) ([]models.Launch, error) {
	return m.GetPastFunc(ctx)
}

