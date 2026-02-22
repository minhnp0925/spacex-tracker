package services

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"
)

type mockCache struct {
	getData   []byte
	getErr    error
	setCalled bool
	setValue  []byte
	setKey    string
}

func (m *mockCache) Get(ctx context.Context, key string) ([]byte, error) {
	return m.getData, m.getErr
}

func (m *mockCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	m.setCalled = true
	m.setValue = value
	m.setKey = key
	return nil
}

func TestGetOrSet_CacheHit(t *testing.T) {
	type testData struct {
		Name string
	}

	expected := testData{Name: "Falcon 9"}
	bytes, _ := json.Marshal(expected)

	cache := &mockCache{
		getData: bytes,
		getErr:  nil,
	}

	fetchCalled := false

	result, err := getOrSet(
		context.Background(),
		cache,
		"key",
		time.Minute,
		func(ctx context.Context) (testData, error) {
			fetchCalled = true
			return testData{}, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if fetchCalled {
		t.Fatal("fetch should not be called on cache hit")
	}

	if result.Name != "Falcon 9" {
		t.Fatal("unexpected result: ", result.Name)
	}
}

func TestGetOrSet_CacheMiss(t *testing.T) {
	type testData struct {
		Name string
	}

	cache := &mockCache{
		getErr: errors.New("miss"),
	}

	fetchCalled := false

	result, err := getOrSet(
		context.Background(),
		cache,
		"key",
		time.Minute,
		func(ctx context.Context) (testData, error) {
			fetchCalled = true
			return testData{Name: "Starship"}, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fetchCalled {
		t.Fatal("fetch should be called on cache miss")
	}

	if !cache.setCalled {
		t.Fatal("cache.Set should be called on miss")
	}

	if result.Name != "Starship" {
		t.Fatal("unexpected result")
	}
}

func TestGetOrSet_FetchError(t *testing.T) {
	cache := &mockCache{
		getErr: errors.New("miss"),
	}

	result, err := getOrSet(
		context.Background(),
		cache,
		"key",
		time.Minute,
		func(ctx context.Context) (string, error) {
			return "", errors.New("fetch failed")
		},
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if result != "" {
		t.Fatal("expected zero value")
	}

	if cache.setCalled {
		t.Fatal("cache.Set should NOT be called on fetch error")
	}
}

func TestGetOrSet_InvalidCacheData(t *testing.T) {
	type testData struct {
		Name string
	}

	cache := &mockCache{
		getData: []byte("invalid-json"),
		getErr:  nil,
	}

	fetchCalled := false

	result, err := getOrSet(
		context.Background(),
		cache,
		"key",
		time.Minute,
		func(ctx context.Context) (testData, error) {
			fetchCalled = true
			return testData{Name: "Recovered"}, nil
		},
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !fetchCalled {
		t.Fatal("fetch should be called if unmarshal fails")
	}

	if result.Name != "Recovered" {
		t.Fatal("unexpected result")
	}
}