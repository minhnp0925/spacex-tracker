package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"spacex-tracker/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type mockLaunchService struct {
	nextResult *models.Launch
	nextErr    error

	latestResult *models.Launch
	latestErr    error

	upcomingResult []models.Launch
	upcomingErr    error

	pastResult []models.Launch
	pastErr    error
}

func (m *mockLaunchService) GetNext(ctx context.Context) (*models.Launch, error) {
	return m.nextResult, m.nextErr
}

func (m *mockLaunchService) GetLatest(ctx context.Context) (*models.Launch, error) {
	return m.latestResult, m.latestErr
}

func (m *mockLaunchService) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	return m.upcomingResult, m.upcomingErr
}

func (m *mockLaunchService) GetPast(ctx context.Context, sortOrder string) ([]models.Launch, error) {
	return m.pastResult, m.pastErr
}

func setupRouter(service *mockLaunchService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	handler := NewLaunchHandler(service)

	r := gin.New()
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

	return r
}

func TestGetNext_Success(t *testing.T) {
	mockSvc := &mockLaunchService{
		nextResult: &models.Launch{
			Name: "Falcon 9",
		},
	}

	router := setupRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/launches/next", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Falcon 9") {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}

func TestGetNext_Error(t *testing.T) {
	mockSvc := &mockLaunchService{
		nextErr: errors.New("service failed"),
	}

	router := setupRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/launches/next", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetUpcoming_Success(t *testing.T) {
	mockSvc := &mockLaunchService{
		upcomingResult: []models.Launch{
			{Name: "Starship"},
			{Name: "Falcon Heavy"},
		},
	}

	router := setupRouter(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/launches/upcoming", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Starship") {
		t.Fatal("missing expected launch")
	}
}

func TestGetPast_WithQuery(t *testing.T) {
	mockSvc := &mockLaunchService{
		pastResult: []models.Launch{
			{Name: "Old Falcon"},
		},
	}

	router := setupRouter(mockSvc)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/launches/past?sort=asc",
		nil,
	)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}