package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"spacex-tracker/configs"
	"spacex-tracker/models"
)

type SpaceXClient interface {
	GetNext(ctx context.Context) (*models.Launch, error)
	GetLatest(ctx context.Context) (*models.Launch, error)
	GetUpcoming(ctx context.Context) ([]models.Launch, error)
	GetPast(ctx context.Context) ([]models.Launch, error)
}

type concreteSpaceXClient struct {
	base_url string
	client   *http.Client
}

func NewSpaceXClient(cfg *configs.Config) SpaceXClient {
	return &concreteSpaceXClient{
		base_url: cfg.ClientBaseURL,
		client: &http.Client{
			Timeout: cfg.ClientTimeout,
		},
	}
}

func (c *concreteSpaceXClient) getOne(ctx context.Context, url string) (*models.Launch, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status: %d", response.StatusCode)
	}

	var launch models.Launch
	if err := json.NewDecoder(response.Body).Decode(&launch); err != nil {
		return nil, err
	}

	return &launch, nil
}

func (c *concreteSpaceXClient) getList(ctx context.Context, url string) ([]models.Launch, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status: %d", response.StatusCode)
	}

	var launches []models.Launch
	if err := json.NewDecoder(response.Body).Decode(&launches); err != nil {
		return nil, err
	}

	return launches, nil
}

func (c *concreteSpaceXClient) GetNext(ctx context.Context) (*models.Launch, error) {
	url := fmt.Sprintf("%s/launches/next", c.base_url)
	return c.getOne(ctx, url)
}

func (c *concreteSpaceXClient) GetLatest(ctx context.Context) (*models.Launch, error) {
	url := fmt.Sprintf("%s/launches/latest", c.base_url)
	return c.getOne(ctx, url)
}

func (c *concreteSpaceXClient) GetUpcoming(ctx context.Context) ([]models.Launch, error) {
	url := fmt.Sprintf("%s/launches/upcoming", c.base_url)
	return c.getList(ctx, url)
}

func (c *concreteSpaceXClient) GetPast(ctx context.Context) ([]models.Launch, error) {
	url := fmt.Sprintf("%s/launches/past", c.base_url)
	return c.getList(ctx, url)
}
