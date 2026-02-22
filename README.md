A RESTful backend service built with Go and Gin that tracks SpaceX launch data by integrating with the public SpaceX API.

## Instructions

### Running locally

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
### Running on docker

#### Build the container
```bash
docker-compose up --build
```

Service will start at http://localhost:8080

## Deployment
The service is currently deployed at https://spacex-tracker.onrender.com