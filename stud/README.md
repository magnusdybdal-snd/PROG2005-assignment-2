# API Stub Service

Mock HTTP service for third-party API development

## Quick Start

1. Run the service:
```
go run cmd/RESTstub/main.go
```

2. Default endpoints:
```
GET /countries/no  # Returns mock country data
GET /weather/no   # Returns mock weather data
GET /currency/no # Returns mock currency data
```

## Configuration

- Port: 8081 (override with PORT env var)
- Test files: /testdata/*.json

## Development

1. Add new JSON files to /testdata/
2. Create handlers in internal/stubs/
3. Register routes in cmd/RESTstub/main.go