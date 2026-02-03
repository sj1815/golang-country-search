# Golang Country Search API

A REST API service that provides country information by leveraging the [REST Countries API](https://restcountries.com/). The service includes custom in-memory caching and comprehensive unit testing.

## Features

- Search countries by name
- In-memory caching (thread-safe)
- Configurable timeouts
- Graceful shutdown
- 90%+ test coverage
- Race condition safe

## Project Structure

```
golang-country-search/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── cache/
│   │   ├── cache.go             # Thread-safe cache implementation
│   │   └── cache_test.go
│   ├── client/
│   │   ├── client.go            # HTTP client for REST Countries API
│   │   └── client_test.go
│   ├── config/
│   │   ├── config.go            # Configuration and dependency injection
│   │   └── config_test.go
│   ├── handler/
│   │   ├── countries.go         # HTTP handlers
│   │   └── countries_test.go
│   ├── model/
│   │   └── country.go           # Data models
│   ├── router/
│   │   ├── router.go            # Route definitions
│   │   └── router_test.go
│   └── service/
│       ├── countries.go         # Business logic
│       └── countries_test.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Requirements

- Go 1.21 or higher

## Installation

1. Clone the repository:
```bash
git clone https://github.com/saurabhj/golang-country-search.git
cd golang-country-search
```

2. Install dependencies:
```bash
make deps
# or
go mod tidy
```

## Running the Server

### Using Make
```bash
make run
```

### Using Go directly
```bash
go run cmd/server/main.go
```

### Build and run
```bash
go build -o server cmd/server/main.go
./server
```

The server will start on `http://localhost:8000`

## API Documentation

### Search Country

Search for a country by name.

**Endpoint:** `GET /api/countries/search`

**Query Parameters:**
| Parameter | Type   | Required | Description          |
|-----------|--------|----------|----------------------|
| name      | string | Yes      | Country name to search |

**Success Response (200 OK):**
```json
{
  "name": "India",
  "capital": "New Delhi",
  "currency": "₹",
  "population": 1417492000
}
```

**Error Responses:**

- `400 Bad Request` - Missing name parameter
```json
{
  "error": "Bad Request",
  "message": "name query parameter is required"
}
```

- `404 Not Found` - Country not found
```json
{
  "error": "Not Found",
  "message": "country not found: InvalidCountry"
}
```

### Examples

```bash
# Search for India
curl "http://localhost:8000/api/countries/search?name=India"

# Search for United States
curl "http://localhost:8000/api/countries/search?name=United%20States"

# Search for Germany
curl "http://localhost:8000/api/countries/search?name=Germany"
```

## Testing

### Run all tests
```bash
make test
```

### Run tests with race detector
```bash
make test-race
```

### Run tests with coverage
```bash
make test-coverage
```

### Generate HTML coverage report
```bash
make test-coverage-report
# Open coverage.html in browser
```

## Test Coverage

| Package   | Coverage |
|-----------|----------|
| cache     | 100%     |
| client    | 94.7%    |
| config    | 100%     |
| handler   | 94.7%    |
| router    | 100%     |
| service   | 100%     |

## Architecture

### Cache Implementation
- Custom in-memory cache built from scratch (no external libraries)
- Thread-safe using `sync.RWMutex`
- Supports concurrent reads with exclusive writes

### HTTP Client
- Configurable timeout
- Context support for cancellation
- Proper error handling

### Service Layer
- Business logic separation
- Cache interaction
- Data transformation

### Graceful Shutdown
- Handles `SIGINT` and `SIGTERM` signals
- Waits for ongoing requests to complete
- Configurable shutdown timeout

## Configuration

Default configuration values (can be modified in `internal/config/config.go`):

| Setting            | Default Value |
|--------------------|---------------|
| Server Port        | :8000         |
| HTTP Client Timeout| 10 seconds    |
| Server Read Timeout| 15 seconds    |
| Server Write Timeout| 15 seconds   |
| Shutdown Timeout   | 10 seconds    |

## License

MIT License
