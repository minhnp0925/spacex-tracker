A RESTful backend service built with Go and Gin that tracks SpaceX launch data by integrating with the public SpaceX API. The service also implements a Redis cache using cache-aside strategy.

## Endpoints

| API | Method | Path | Description |
|-----|--------|------|-------------|
| Next launch | GET | `/api/v1/launches/next` | Returns the next launch. |
| Latest launch | GET | `/api/v1/launches/latest` | Returns the latest launch. |
| Upcoming launches | GET | `/api/v1/launches/upcoming` | Returns an array of upcoming launches. |
| Past launches | GET | `/api/v1/launches/past` | Returns an array of past launches. Optional query param `?sort=asc\|desc` to sort by time. Defaults to `desc`.|

## Response schema
```go
type Launch struct {
    Id string `json:"id"`
    Name string `json:"name"`
    DateUTC time.Time `json:"date_utc"`
    Success *bool `json:"success,omitempty"`
    Upcoming bool `json:"upcoming"`
    Details string `json:"details,omitempty"`
}
```

## Instructions

### Run locally

#### Step 1: Add environment variables
```bash
cp .env.example .env
```
The following environment variables are required:

| Variable          | Description                                                                | Example                         |
| ----------------- | -------------------------------------------------------------------------- | ------------------------------- |
| `REDIS_URL`       | Redis connection string used for caching (optional in local fallback mode) | `redis://redis:6379`            |
| `CLIENT_BASE_URL` | Base URL of the SpaceX public API                                          | `https://api.spacexdata.com/v4` |
| `CLIENT_TIMEOUT`  | HTTP client timeout (in seconds)                                           | `5`                             |
| `CACHE_TTL`       | Cache time-to-live in seconds for GET responses                            | `60`                            |

> If `REDIS_URL` is not set or invalid, the service will automatically fall back to running without caching.

#### Step 2: Install dependencies
```bash
go mod download
```
#### Step 3: Run the service

```bash
go run .
```
### Run on docker

#### Step 1: Add environment variables
```bash
cp .env.example .env
```
The following environment variables are required:

| Variable          | Description                                                                | Example                         |
| ----------------- | -------------------------------------------------------------------------- | ------------------------------- |
| `REDIS_URL`       | Redis connection string used for caching (optional in local fallback mode) | `redis://redis:6379`            |
| `CLIENT_BASE_URL` | Base URL of the SpaceX public API                                          | `https://api.spacexdata.com/v4` |
| `CLIENT_TIMEOUT`  | HTTP client timeout (in seconds)                                           | `5`                             |
| `CACHE_TTL`       | Cache time-to-live in seconds for GET responses                            | `60`                            |

> If `REDIS_URL` is not set or invalid, the service will automatically fall back to running without caching.

#### Step 2: Build the container
```bash
docker-compose up --build
```
Service will start at http://localhost:8080

### Run unit tests
This repository contains unit tests for the service layer (service logic and caching logic) and handler (httptest). To run tests:
```bash
go test ./...
```

## Deployment
The service is currently deployed at https://spacex-tracker.onrender.com