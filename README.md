A RESTful backend service built with Go and Gin that tracks SpaceX launch data by integrating with the public SpaceX API.

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

#### Step 2: Install dependencies
```bash
go mod download
```
#### Step 3: Run the service

```bash
go run .
```
### Run on docker

#### Build the container
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