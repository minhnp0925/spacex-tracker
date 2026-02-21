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
func TestGetNext_Success(t *testing.T) {
	expected := &models.Launch{
		Id:   "1",
		Name: "Falcon 9",
	}

	mock := &MockSpaceXClient{
		GetNextFunc: func(ctx context.Context) (*models.Launch, error) {
			return expected, nil
		},
	}

	service := NewBaseLaunchService(mock)

	result, err := service.GetNext(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetNext_Error(t *testing.T) {
	mock := &MockSpaceXClient{
		GetNextFunc: func(ctx context.Context) (*models.Launch, error) {
			return nil, errors.New("client error")
		},
	}

	service := NewBaseLaunchService(mock)

	_, err := service.GetNext(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLatest_Success(t *testing.T) {
	expected := &models.Launch{
		Id:   "3",
		Name: "Crew-5",
	}

	mock := &MockSpaceXClient{
		GetLatestFunc: func(ctx context.Context) (*models.Launch, error) {
			return expected, nil
		},
	}

	service := NewBaseLaunchService(mock)

	result, err := service.GetLatest(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetLatest_Error(t *testing.T) {
	mock := &MockSpaceXClient{
		GetLatestFunc: func(ctx context.Context) (*models.Launch, error) {
			return nil, errors.New("client error")
		},
	}

	service := NewBaseLaunchService(mock)

	_, err := service.GetLatest(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUpcoming_Success(t *testing.T) {
	expected := []models.Launch{
		{Id: "1"},
		{Id: "2"},
	}

	mock := &MockSpaceXClient{
		GetUpcomingFunc: func(ctx context.Context) ([]models.Launch, error) {
			return expected, nil
		},
	}

	service := NewBaseLaunchService(mock)

	result, err := service.GetUpcoming(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %+v, got %+v", expected, result)
	}
}

func TestGetUpcoming_Error(t *testing.T) {
	mock := &MockSpaceXClient{
		GetUpcomingFunc: func(ctx context.Context) ([]models.Launch, error) {
			return nil, errors.New("client error")
		},
	}

	service := NewBaseLaunchService(mock)

	_, err := service.GetUpcoming(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPast_SortDesc(t *testing.T) {
	now := time.Now()

	launches := []models.Launch{
		{Id: "1", DateUTC: now.Add(-48 * time.Hour)},
		{Id: "2", DateUTC: now},
		{Id: "3", DateUTC: now.Add(-24 * time.Hour)},
	}

	mock := &MockSpaceXClient{
		GetPastFunc: func(ctx context.Context) ([]models.Launch, error) {
			return launches, nil
		},
	}

	service := NewBaseLaunchService(mock)

	result, err := service.GetPast(context.Background(), "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOrder := []string{"2", "3", "1"}

	for i, launch := range result {
		if launch.Id != expectedOrder[i] {
			t.Errorf("expected %s at index %d, got %s",
				expectedOrder[i], i, launch.Id)
		}
	}
}

func TestGetPast_SortAsc(t *testing.T) {
	now := time.Now()

	launches := []models.Launch{
		{Id: "1", DateUTC: now.Add(-48 * time.Hour)},
		{Id: "2", DateUTC: now},
		{Id: "3", DateUTC: now.Add(-24 * time.Hour)},
	}

	mock := &MockSpaceXClient{
		GetPastFunc: func(ctx context.Context) ([]models.Launch, error) {
			return launches, nil
		},
	}

	service := NewBaseLaunchService(mock)

	result, err := service.GetPast(context.Background(), "asc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedOrder := []string{"1", "3", "2"}

	for i, launch := range result {
		if launch.Id != expectedOrder[i] {
			t.Errorf("expected %s at index %d, got %s",
				expectedOrder[i], i, launch.Id)
		}
	}
}

func TestGetPast_Error(t *testing.T) {
	mock := &MockSpaceXClient{
		GetPastFunc: func(ctx context.Context) ([]models.Launch, error) {
			return nil, errors.New("client error")
		},
	}

	service := NewBaseLaunchService(mock)

	_, err := service.GetPast(context.Background(), "desc")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}