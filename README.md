# Entorno35 - NOM-035 Compliance Platform

Cloud-native SaaS platform for automating NOM-035-STPS-2018 compliance for Mexican organizations.

## Project Structure

```
Entorno35/
├── cmd/                  # Backend application entry points
│   ├── api/              # Main API server (Go)
│   └── seeder/           # Database seeder
├── internal/             # Backend private application code (Go)
│   ├── adapters/         # HTTP and database adapters
│   ├── auth/             # Authentication context
│   ├── config/           # Configuration management
│   ├── core/             # Core business logic
│   ├── database/         # Database connection utilities
│   ├── domain/           # Domain models and entities
│   └── middleware/       # HTTP middleware
├── pkg/                  # Backend public reusable packages (Go)
│   └── errors/           # Error handling
├── migrations/           # Database migrations
├── tests/                # Backend tests
│   ├── integration/      # Integration tests
│   └── unit/             # Unit tests
├── web/                  # Frontend application
│   └── frontend/         # Next.js 14 frontend (TypeScript)
│       ├── app/          # Next.js App Router pages
│       ├── components/   # React components
│       ├── lib/          # Utilities (axios, query client)
│       ├── services/     # API service layer
│       └── types/        # TypeScript type definitions
├── docs/                 # Documentation
│   ├── *_ENDPOINTS.md    # API endpoint documentation
│   └── ...
├── nom035_questions.json # Question catalog data
├── risk_strategy.json    # Risk calculation strategy
└── docker-compose.yml    # Docker services configuration

```

### Architecture Separation

- **Backend** (Go): Root directory structure (`cmd/`, `internal/`, `pkg/`, `migrations/`, `tests/`)
- **Frontend** (Next.js/TypeScript): `web/frontend/` directory
- **Documentation**: `docs/` directory
- **Shared Data**: Root level JSON files (question catalog, risk strategy)

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

4. **Run the backend API**
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

5. **Run the frontend (optional, in a separate terminal)**
   ```bash
   cd web/frontend
   npm install  # First time only
   npm run dev
   ```

6. **Verify backend health**
   ```bash
   curl http://localhost:8080/health
   ```

7. **Verify frontend (if running)**
   - Open http://localhost:3000 in your browser

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

**Note:** The `make run` command automatically loads environment variables from the `.env` file. No need to manually export variables.

## Architecture

- **Backend**: Go (Golang) with Gin framework
  - CORS middleware configured for frontend communication
  - Automatic .env file loading via Makefile
- **Frontend**: Next.js 14 with React
  - TanStack Query for data fetching
  - Axios client with automatic token injection
- **Database**: PostgreSQL
- **Cache/Queue**: Redis
- **Infrastructure**: Docker Compose (dev) / Kubernetes (prod)

## Essential Data Files

- **`nom035_questions.json`** - Complete question catalog (138 questions)
- **`docs/Scoring.md`** - Scoring implementation guide with thresholds

Note: Risk thresholds are stored as configuration (see `docs/Scoring.md`) and will be implemented in Phase 2's scoring engine service.

## Git Workflow

- **`main`** - Production-ready code
- **`develop`** - Integration branch (Phase 4 merged)
- **`phase-{n}/{feature}`** - Feature branches for each phase

## Current Phase

**Phase 5: Frontend Implementation** - IN PROGRESS

### Completed
- Phase 1: Core Infrastructure (project setup, database, authentication)
- Phase 2: Scoring Logic & Engine (strategy pattern, risk calculation)
- Phase 3: API Implementation (assessments, responses, staff management)
- Phase 4: Reporting & Compliance Engine (individual reports, general reports, recommendations)
- Phase 5.1: Frontend Architecture & Service Layer (Next.js setup, TypeScript types, Axios client, service layer)

### Completed (Continued)
- Phase 5.2: Authentication UI & Gatekeeper - COMPLETED ✅
  - Login screen with Company/Staff type support
  - Zustand auth store with global state management
  - Route protection (AuthGuard component)
  - Dashboard layout with navigation
  - Integration tests (12 tests passing)

- Phase 5.3: Staff Management UI - COMPLETED ✅
  - Staff data table with pagination, sorting, and filtering
  - CSV uploader with drag-and-drop, progress indication, and error handling
  - Complete staff management page with table and upload integration
  - UI components: table, pagination, dialog, alert, progress
  - Template download functionality
  - All tests passing (20/20 frontend tests)

### Next Steps
1. Phase 5.4: Dashboard (analytics, assessment creation wizard)
2. Phase 5.5: Public Assessment View (staff test-taking interface)
3. Production preparation and deployment
4. Database migration to production (password authentication)

## License

Proprietary
