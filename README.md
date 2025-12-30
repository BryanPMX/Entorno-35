# Entorno35 - NOM-035 Compliance Platform

Cloud-native SaaS platform for automating NOM-035-STPS-2018 compliance for Mexican organizations.

## Project Structure

```
Entorno35/
├── cmd/
│   └── api/              # Application entry point
├── internal/             # Private application code
│   ├── adapters/        # HTTP and database adapters
│   ├── auth/            # Authentication context
│   ├── config/          # Configuration management
│   ├── core/            # Core business logic
│   ├── database/        # Database connection utilities
│   ├── domain/          # Domain models and entities
│   └── middleware/      # HTTP middleware
├── pkg/                  # Public reusable packages
│   ├── errors/          # Error handling
│   └── utils/           # Utility functions
├── migrations/           # Database migrations
├── web/                  # Frontend (Next.js)
└── docs/                 # Documentation

```

## Quick Start

### Prerequisites

- Go 1.21+
- Node.js 18+ (for frontend)
- Docker & Docker Compose
- PostgreSQL 15+ (via Docker)
- Redis 7+ (via Docker)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd Entorno-35
   ```

2. **Start Docker services**
   ```bash
   make docker-up
   # or
   docker-compose up -d
   ```

3. **Install Go dependencies**
   ```bash
   make setup
   # or
   go mod download
   ```

4. **Run the application**
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

5. **Verify health**
   ```bash
   curl http://localhost:8080/health
   ```

## Development

### Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests with coverage
make clean         # Clean build artifacts
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make migrate-up    # Run database migrations
make migrate-down  # Rollback database migrations
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

## Architecture

- **Backend**: Go (Golang) with Gin framework
- **Frontend**: Next.js 14 with React
- **Database**: PostgreSQL
- **Cache/Queue**: Redis
- **Infrastructure**: Docker Compose (dev) / Kubernetes (prod)

## Essential Data Files

- **`nom035_questions.json`** - Complete question catalog (138 questions)
- **`docs/Scoring.md`** - Scoring implementation guide with thresholds

Note: Risk thresholds are stored as configuration (see `docs/Scoring.md`) and will be implemented in Phase 2's scoring engine service.

## Git Workflow

- **`main`** - Production-ready code
- **`develop`** - Integration branch (Phase 3 merged)
- **`phase-{n}/{feature}`** - Feature branches for each phase

## Current Phase

**Phase 3: API Implementation** - COMPLETED

### Completed
- Phase 1: Core Infrastructure (project setup, database, authentication)
- Phase 2: Scoring Logic & Engine (strategy pattern, risk calculation)
- Phase 3: API Implementation (assessments, responses, staff management)

### Next Steps
1. Phase 4: Frontend Implementation
2. Production preparation and deployment

## License

Proprietary
