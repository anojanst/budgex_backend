# Budgex Backend

A Go-based backend API for the Budgex personal finance application, built with Fiber web framework and Clerk authentication.

## Project Structure

```
budgex_backend/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP request handlers
│   │   │   ├── budget.go        # Budget CRUD operations
│   │   │   ├── category.go      # Category CRUD operations
│   │   │   ├── health.go        # Health check endpoint
│   │   │   ├── me.go           # User info endpoint
│   │   │   └── transaction.go   # Transaction CRUD operations
│   │   ├── middleware/          # HTTP middleware
│   │   │   ├── clerk_fiber.go   # Clerk authentication middleware
│   │   │   ├── logz.go         # Structured JSON logging
│   │   │   └── trace.go        # OpenTelemetry tracing
│   │   └── router.go           # API routes and middleware setup
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── db/
│   │   └── db.go              # Database connection and migrations
│   ├── docs/                  # Generated Swagger documentation
│   ├── models/                # Database models
│   │   └── models.go
│   └── observability/         # Logging and tracing
│       ├── logger.go          # Zap structured logging
│       └── tracing.go         # OpenTelemetry tracing
├── .env                       # Environment variables
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
└── makefile                   # Build and run commands
```

## Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **PostgreSQL 13+** - [Download](https://www.postgresql.org/download/)
- **Clerk Account** - [Sign up](https://clerk.com/) for authentication

## Environment Variables

Create a `.env` file in the project root with the following variables:

```env
# Server Configuration
PORT=8080
SERVICE_NAME=budgex-backend
LOG_LEVEL=info

# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/budgex?sslmode=disable

# Clerk Authentication
CLERK_SECRET_KEY=sk_test_your_clerk_secret_key_here

# OpenTelemetry (Optional)
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
OTEL_EXPORTER_OTLP_HEADERS=api-key=your-api-key
```

### Environment Variable Descriptions

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `PORT` | Server port | No | 8080 |
| `SERVICE_NAME` | Service name for logging/tracing | No | budgex-backend |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | No | info |
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `CLERK_SECRET_KEY` | Clerk secret key for JWT verification | Yes | - |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OpenTelemetry collector endpoint | No | - |
| `OTEL_EXPORTER_OTLP_HEADERS` | Headers for OTLP exporter | No | - |

## Setup Steps

### 1. Clone and Navigate
```bash
git clone https://github.com/anojanst/budgex_backend.git
cd budgex_backend
```

### 2. Install Dependencies
```bash
go mod tidy
```

### 3. Install Development Tools
```bash
# Swagger documentation generator
go install github.com/swaggo/swag/cmd/swag@latest

# Add to PATH (if needed)
export PATH=$PATH:$(go env GOPATH)/bin
```

### 4. Database Setup
```bash
# Create PostgreSQL database
createdb budgex

# Or using psql
psql -U postgres -c "CREATE DATABASE budgex;"
```

### 5. Configure Environment
```bash
# Copy example environment file
cp .env.example .env

# Edit with your values
nano .env
```

### 6. Run Database Migrations
```bash
go run cmd/server/main.go
# Migrations run automatically on startup
```

### 7. Generate Swagger Documentation
```bash
swag init -g cmd/server/main.go -o internal/docs
```

## Common Commands

### Development Commands

```bash
# Run the server
go run cmd/server/main.go

# Run with hot reload (if using air)
air

# Run in background
go run cmd/server/main.go &

# Stop background server
pkill -f "go run cmd/server/main.go"
```

### Build Commands

```bash
# Build binary
go build -o bin/server cmd/server/main.go

# Run binary
./bin/server

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o bin/server-linux cmd/server/main.go
```

### Database Commands

```bash
# Connect to database
psql -U postgres -d budgex

# Reset database (drop and recreate)
dropdb budgex && createdb budgex

# View database schema
psql -U postgres -d budgex -c "\dt"
```

### Documentation Commands

```bash
# Generate Swagger docs
swag init -g cmd/server/main.go -o internal/docs

# Regenerate docs after changes
swag init -g cmd/server/main.go -o internal/docs --parseDependency --parseInternal

# View Swagger UI
open http://localhost:8080/swagger/index.html
```

### Testing Commands

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test ./internal/api/handlers -v
```

### Debugging Commands

```bash
# Check server status
curl http://localhost:8080/api/health

# Test authentication
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/me

# View server logs
tail -f /var/log/budgex-backend.log  # if logging to file
```

## API Endpoints

### Public Endpoints
- `GET /api/health` - Health check
- `GET /swagger/index.html` - Swagger UI

### Protected Endpoints (Require Bearer Token)

#### Authentication
- `GET /api/me` - Get current user ID

#### Transactions
- `GET /api/transactions/` - List transactions
- `POST /api/transactions/` - Create transaction

#### Categories
- `GET /api/categories/` - List categories
- `POST /api/categories/` - Create category

#### Budgets
- `GET /api/budgets/` - List budgets
- `POST /api/budgets/` - Upsert budget

## Authentication

The API uses Clerk for authentication. Include the JWT token in the Authorization header:

```bash
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8080/api/me
```

## Logging and Monitoring

### Structured JSON Logging
The application uses Zap for structured JSON logging with the following fields:
- `timestamp` - ISO8601 timestamp
- `level` - Log level (info, warn, error)
- `msg` - Log message
- `service` - Service name
- `trace_id` - OpenTelemetry trace ID
- `span_id` - OpenTelemetry span ID
- `user_id` - Authenticated user ID
- `request_id` - Unique request ID

### OpenTelemetry Tracing
Configure tracing by setting the `OTEL_EXPORTER_OTLP_ENDPOINT` environment variable.

## Troubleshooting

### Common Issues

1. **"Request Header Fields Too Large" Error**
   - Solution: The server is configured with 64KB header buffer to handle large JWT tokens

2. **"Address already in use" Error**
   - Solution: Kill existing processes: `pkill -f "go run cmd/server/main.go"`

3. **"CLERK_SECRET_KEY environment variable is not set"**
   - Solution: Ensure `.env` file exists and contains valid Clerk secret key

4. **Database connection errors**
   - Solution: Verify PostgreSQL is running and `DATABASE_URL` is correct

5. **Unauthorized errors with valid token**
   - Solution: Check if token is expired or Clerk secret key is correct

### Debug Mode

Enable debug logging by setting:
```env
LOG_LEVEL=debug
```

## Development Workflow

1. Make code changes
2. Run tests: `go test ./...`
3. Generate docs: `swag init -g cmd/server/main.go -o internal/docs`
4. Start server: `go run cmd/server/main.go`
5. Test endpoints via Swagger UI or curl
6. Commit changes

## Production Deployment

1. Build binary: `go build -o bin/server cmd/server/main.go`
2. Set production environment variables
3. Run migrations
4. Start server: `./bin/server`
5. Configure reverse proxy (nginx/Apache)
6. Set up monitoring and logging

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes
4. Add tests
5. Update documentation
6. Submit a pull request

## License

[Add your license here]
