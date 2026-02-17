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

type MockSpaceXClient struct{}
func (m *MockSpaceXClient) GetNext(ctx context.Context) (*models.Launch, error) {
	return nil, nil
}
func (m *MockSpaceXClient) GetLatest(ctx context.Context) (*models.Launch, error) {
	return nil, nil
}
func (m *MockSpaceXClient) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return nil, nil
}
func (m *MockSpaceXClient) GetPast(ctx context.Context) ([]models.Launch, error) {
	return nil, nil
}

var client = &MockSpaceXClient{}