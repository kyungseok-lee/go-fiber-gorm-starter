# Spindle API (fiber-gorm-starter)

A production-ready REST API skeleton built with Go Fiber v2 and GORM, supporting both MySQL and PostgreSQL databases.

## Features

- 🚀 **Go Fiber v2** - Fast HTTP framework
- 🗄️ **GORM** - Powerful ORM with MySQL/PostgreSQL support
- 🔄 **Database Migration** - golang-migrate integration
- 📊 **Observability** - Prometheus metrics, structured logging (zap)
- 🛡️ **Security** - CORS, security headers, API key authentication
- 🐳 **Docker** - Multi-stage builds with distroless images
- 🧪 **Testing** - Unit tests and integration tests
- 📖 **API Documentation** - Swagger/OpenAPI integration
- ⚡ **Performance** - k6 load testing scripts
- 🔧 **CI/CD** - GitHub Actions with linting, testing, security scanning

## Quick Start

### Prerequisites

- Go 1.22+
- Docker & Docker Compose
- Make (optional, for convenience commands)

### Local Development

1. **Clone the repository**
   ```bash
   git clone https://github.com/kyungseok-lee/go-fiber-gorm-starter.git
   cd fiber-gorm-starter
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Start database (MySQL)**
   ```bash
   docker-compose --profile mysql up -d
   ```

4. **Run database migrations**
   ```bash
   ./scripts/migrate.sh up
   ```

5. **Seed database (optional)**
   ```bash
   go run scripts/seed.go
   ```

6. **Start the application**
   ```bash
   go run cmd/server/main.go
   ```

7. **Access the API**
   - API: http://localhost:8080
   - Health: http://localhost:8080/health
   - Swagger Docs: http://localhost:8080/docs/index.html
   - Metrics: http://localhost:8080/metrics

### Using PostgreSQL Instead

To switch to PostgreSQL:

1. **Update environment**
   ```bash
   # In .env file
   DB_DRIVER=postgres
   DB_PORT=5432
   ```

2. **Start PostgreSQL**
   ```bash
   docker-compose --profile postgres up -d
   ```

3. **Run migrations**
   ```bash
   ./scripts/migrate.sh up
   ```

## API Endpoints

### Users
- `GET /v1/users` - List users with pagination
- `GET /v1/users/:id` - Get user by ID
- `POST /v1/users` - Create new user
- `PUT /v1/users/:id` - Update user
- `DELETE /v1/users/:id` - Delete user

### System
- `GET /health` - Health check
- `GET /ready` - Readiness check (includes dependencies)
- `GET /metrics` - Prometheus metrics
- `GET /docs/*` - Swagger documentation (dev only)

## Database Management

### Migrations

```bash
# Apply all migrations
./scripts/migrate.sh up

# Apply specific number of migrations
./scripts/migrate.sh up 1

# Rollback migrations
./scripts/migrate.sh down 1

# Check migration status
./scripts/migrate.sh status

# Create new migration
./scripts/migrate.sh create add_user_profile
```

### Switching Databases

The application supports both MySQL and PostgreSQL. Switch by setting environment variables:

**MySQL:**
```env
DB_DRIVER=mysql
DB_PORT=3306
# MySQL connection string format is automatically handled
```

**PostgreSQL:**
```env
DB_DRIVER=postgres
DB_PORT=5432
DB_SSL_MODE=disable
# PostgreSQL connection string format is automatically handled
```

## Docker Deployment

### Using Docker Compose

**With MySQL:**
```bash
docker-compose --profile mysql --profile app up -d
```

**With PostgreSQL:**
```bash
docker-compose --profile postgres --profile app up -d
```

### Building Docker Image

```bash
# Build image
docker build -t fiber-gorm-starter .

# Run container
docker run -p 8080:8080 --env-file .env fiber-gorm-starter
```

## Development Tools

### Make Commands

```bash
# Run the application
make run

# Run tests
make test

# Run linter
make lint

# Build binary
make build

# Generate Swagger docs
make swag

# Database operations
make migrate-up
make migrate-down
make seed

# Docker operations
make docker-up     # MySQL
make docker-up-pg  # PostgreSQL
make docker-down

# Format code
make fmt

# Run all checks (lint + test + build)
make check
```

### Code Generation

Generate Swagger documentation:
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/server/main.go -o ./docs
```

## Testing

### Unit Tests
```bash
go test -v ./...
```

### Integration Tests
```bash
# Start test database
docker-compose --profile mysql up -d

# Run tests with database
ENV=test DB_NAME=fiber_gorm_starter_test go test -v ./...
```

### Load Testing
```bash
# Install k6
# https://k6.io/docs/getting-started/installation/

# Run performance tests
k6 run scripts/k6/users-smoke.js
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment (local/dev/prod) | `local` |
| `PORT` | Server port | `8080` |
| `DB_DRIVER` | Database driver (mysql/postgres) | `mysql` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `3306` |
| `DB_USER` | Database user | `user` |
| `DB_PASS` | Database password | `password` |
| `DB_NAME` | Database name | `fiber_gorm_starter` |
| `DB_SSL_MODE` | SSL mode (postgres only) | `disable` |
| `DB_MAX_OPEN` | Max open connections | `25` |
| `DB_MAX_IDLE` | Max idle connections | `10` |
| `DB_MAX_LIFETIME` | Connection max lifetime | `300s` |
| `API_KEY` | API key for authentication | `` |
| `LOG_LEVEL` | Logging level | `info` |
| `METRICS_ENABLED` | Enable Prometheus metrics | `true` |
| `PPROF_ENABLED` | Enable pprof endpoints | `false` |

### Database Connection Pool

Optimize connection pool settings based on your workload:

```env
# For high-traffic applications
DB_MAX_OPEN=100
DB_MAX_IDLE=20
DB_MAX_LIFETIME=600s

# For low-traffic applications
DB_MAX_OPEN=10
DB_MAX_IDLE=5
DB_MAX_LIFETIME=300s
```

## Monitoring & Observability

### Metrics

Prometheus metrics are available at `/metrics`:
- HTTP request duration and count
- HTTP request size
- Database connection pool stats
- Go runtime metrics

### Health Checks

- `/health` - Basic application health
- `/ready` - Readiness check including database connectivity

### Logging

Structured logging with zap:
- Request/response logging with correlation IDs
- Different log levels for different environments
- PII-safe logging practices

### Profiling

Enable pprof for debugging (development only):
```env
PPROF_ENABLED=true
```

Access profiling endpoints at `/debug/pprof/`

## Security

### Security Headers
- HSTS with preload
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block

### CORS Policy
- Development: Allow all origins
- Production: Restricted to specific domains (configure in middleware)

### API Authentication
Simple API key authentication (expandable to JWT):
```bash
curl -H "Authorization: Bearer your-api-key" http://localhost:8080/v1/users
```

## Performance Optimization

### Connection Pool Tuning
- Monitor `db_connections_open` metric
- Adjust `DB_MAX_OPEN` based on concurrent load
- Set appropriate `DB_MAX_LIFETIME` for connection recycling

### Log Sampling
For high-traffic applications, implement log sampling:
```go
// Example: Sample 10% of requests in production
if env == "prod" && rand.Float64() > 0.1 {
    return // Skip logging
}
```

### Caching Strategy
Add Redis for caching (infrastructure ready):
```bash
docker-compose --profile redis --profile mysql --profile app up -d
```

## Deployment Strategies

### Blue-Green Deployment
1. Deploy new version alongside current
2. Switch traffic using load balancer
3. Monitor metrics and rollback if needed

### Canary Deployment
1. Route small percentage of traffic to new version
2. Gradually increase traffic if metrics are healthy
3. Full rollout or rollback based on results

### Rolling Updates
1. Update instances one by one
2. Wait for health checks before proceeding
3. Maintain zero-downtime deployment

## Future Enhancements

### Authentication & Authorization
- [ ] JWT token-based authentication
- [ ] Role-based access control (RBAC)
- [ ] OAuth2/OIDC integration
- [ ] Rate limiting per user

### Caching
- [ ] Redis integration
- [ ] Cache-aside pattern implementation
- [ ] Cache invalidation strategies

### Database
- [ ] Read replica support
- [ ] Database sharding
- [ ] Query optimization and indexing
- [ ] Soft delete with audit logs

### Observability
- [ ] OpenTelemetry distributed tracing
- [ ] Custom business metrics
- [ ] Log aggregation (ELK stack)
- [ ] APM integration

### Testing
- [ ] Contract testing with Pact
- [ ] Chaos engineering tests
- [ ] Performance regression tests
- [ ] Security penetration tests

### DevOps
- [ ] Kubernetes manifests
- [ ] Helm charts
- [ ] ArgoCD GitOps
- [ ] Infrastructure as Code (Terraform)

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make changes and add tests
4. Run linting and tests: `make check`
5. Commit with conventional commits format
6. Push and create a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Spindle** - A production-ready Go API starter template