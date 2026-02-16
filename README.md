# Entorno35 - NOM-035 Compliance Platform

A comprehensive SaaS solution for Mexican organizations to automate NOM-035-STPS-2018 workplace psychosocial risk assessments. Streamlines compliance reporting, staff evaluations, and risk analysis with professional dashboards and secure assessment interfaces.

## Overview

NOM-035-STPS-2018 is a Mexican federal standard that requires employers to identify, analyze, and prevent psychosocial risk factors in the workplace. Entorno35 automates this compliance process through:

- **Staff Assessment Management**: Secure, anonymous evaluation links sent to employees
- **Automated Risk Scoring**: Real-time calculation of psychosocial risk levels
- **Compliance Reporting**: Individual and company-wide NOM-035 reports with PDF export
- **Professional Dashboards**: Analytics and monitoring for HR administrators
- **Multi-tenant Architecture**: Secure isolation between organizations

### Key Features

- Complete NOM-035 question catalog (138 questions across 3 guides)
- Automated risk level calculation (Nulo, Bajo, Medio, Alto, Muy Alto)
- High-fidelity PDF export with professional NOM-035 compliant formatting
- Typeform-like assessment experience with smooth animations and transitions
- Complete keyboard accessibility (arrow keys, number selection, Enter confirmation)
- Real-time auto-save indicators and progress tracking
- CSV staff import with CURP validation
- Mobile-responsive assessment interface with focus mode layout
- Professional analytics dashboard with interactive visualizations
- Unified NOM-35 visual design system (shared tokens for auth, portal, assessment, and marketing)
- Centralized NOM-35 risk style registry (single source for labels, colors, and badge classes)
- Frontend visual regression snapshots for critical surfaces (login, dashboard shell, assessment shell)
- Secure token-based assessments with expiration handling
- Multi-company support with tenant isolation

## Project Architecture

### Technology Stack

- **Backend**: Go 1.24 with Gin web framework
- **Frontend**: Next.js 16 with React 19 and TypeScript
- **Database**: PostgreSQL 15 with GORM ORM
- **Cache**: Redis 7
- **Authentication**: JWT with bcrypt password hashing
- **UI Framework**: Tailwind CSS with Radix UI components
- **State Management**: Zustand with TanStack Query

### Directory Structure

```
Entorno35/
├── cmd/                  # Application entry points
│   ├── api/              # Main API server
│   └── seeder/           # Database seeding utility
├── internal/             # Private application code
│   ├── adapters/         # HTTP handlers and database adapters
│   ├── auth/             # Authentication context
│   ├── config/           # Configuration management
│   ├── core/             # Business logic and domain services
│   ├── database/         # Database connection utilities
│   ├── domain/           # Domain models and entities
│   └── middleware/       # HTTP middleware
├── pkg/                  # Public reusable packages
│   └── errors/           # Error handling utilities
├── migrations/           # Database schema migrations
├── tests/                # Backend test suites
│   ├── integration/      # End-to-end tests
│   └── unit/             # Unit tests
├── web/frontend/         # Next.js frontend application
│   ├── app/              # Next.js App Router pages
│   ├── components/       # React components
│   ├── lib/              # Utilities and configurations (includes NOM-35 risk style registry)
│   ├── services/         # API service layer
│   └── types/            # TypeScript type definitions
├── docs/                 # Technical and operational documentation (see docs/README.md)
├── nom035_questions.json # NOM-035 question catalog
├── risk_strategy.json    # Risk calculation configuration
└── docker-compose.prod.yml # Production stack (Portainer)
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+ (for local frontend development)
- Docker & Docker Compose
- Git

### Deployment Architecture

- **Frontend**: Deployed on Vercel (production). Pushes to enabled branches trigger automatic deploys.
- **Backend**: Self-hosted (Docker on Portainer). Pushes to `develop` or `main` that change backend code trigger a GitHub Actions build, Docker Hub push, and optional Portainer stack update.
- **Local Development**: Both frontend and backend can run locally.

For CI/CD setup (Vercel, Portainer, Cloudflare, required secrets), see [docs/CI_CD.md](docs/CI_CD.md).

### Development Setup

#### Quick Start (Recommended)

For local development, start PostgreSQL and Redis (e.g. via Docker or an existing instance), then run the backend and frontend:

```bash
# Clone and enter repository
git clone <repository-url>
cd Entorno-35

# Install dependencies
make setup
cd web/frontend && npm install && cd ../..

# Start backend (terminal 1) and frontend (terminal 2) - see Manual Setup below for env vars
./start-dev.sh
```

**Note**: `start-dev.sh` expects PostgreSQL and Redis to be available (e.g. `make docker-up` if you use a local compose, or point env vars at an existing instance). For production, the frontend is deployed on Vercel and the backend uses `docker-compose.prod.yml` on Portainer.

#### Manual Setup (Advanced)

For more control over individual services:

1. **Clone and enter the repository**
   ```bash
   git clone <repository-url>
   cd Entorno-35
   ```

2. **Start infrastructure services**
   ```bash
   make docker-up
   # Starts PostgreSQL and Redis containers
   ```

3. **Install dependencies**
   ```bash
   # Backend dependencies
   make setup

   # Frontend dependencies
   cd web/frontend && npm install && cd ../..
   ```

4. **Configure environment variables**
   ```bash
   # Required for backend
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=postgres
   export DB_NAME=entorno35
   export DB_SSLMODE=disable
   export JWT_SECRET=your-secret-key-min-32-chars-long
   export CORS_ORIGIN=http://localhost:3000
   ```

5. **Start the backend (terminal 1)**
   ```bash
   cd cmd/api && go run main.go
   ```

6. **Start the frontend (terminal 2)**
   ```bash
   cd web/frontend && npm run dev
   ```

7. **Verify installation**
   ```bash
   # Backend health check
   curl http://localhost:8080/health

   # Frontend: Open browser to http://localhost:3000
   ```

### For Administrators

After setup, create your first company account through the web interface, then:

1. Import staff data via CSV upload
2. Create assessment cycles
3. Generate secure assessment links
4. Monitor completion and view reports

### For Staff

Staff receive assessment links via email and complete evaluations through the secure, mobile-responsive interface.

## Development

### Available Commands

```bash
# Project setup
make setup          # Install Go dependencies
make docker-up      # Start PostgreSQL and Redis
make docker-down    # Stop containers
make migrate-up     # Run database migrations
make migrate-down   # Rollback migrations

# Development
make run            # Start backend API server
make build          # Build backend binary
make test           # Run all tests with coverage
make clean          # Clean build artifacts

# Frontend (from web/frontend/)
npm install         # Install dependencies
npm run dev         # Start development server
npm run build       # Build for production
npm test            # Run test suite
```

### Environment Configuration

Production uses **Portainer stack environment variables** only; no `.env` files on the server. Required backend vars are listed in `docker-compose.prod.yml` (e.g. `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `JWT_SECRET`, `CORS_ORIGIN`, `REDIS_HOST`, `SMTP_*`). Set them in Portainer when creating or editing the stack. Never commit `.env` or `.env.stripe` (they are in `.gitignore`).

### Testing

```bash
# Run all backend tests
make test

# Run frontend tests
cd web/frontend && npm test

# Run frontend linting
cd web/frontend && npm run lint

# Run frontend production build
cd web/frontend && npm run build

# Run integration tests
go test ./tests/integration/...
```

### Documentation

All substantive documentation (API, deployment, scoring, CSV import, SMTP) lives in the **docs/** directory. The only non-README markdown at repo root is **SECURITY.md** (credentials and pre-commit check), kept at root for GitHub’s Security policy link.

- **[docs/README.md](docs/README.md)** – Documentation index and list of essential docs.

Key documents:

- **API and integration**: [docs/API_REFERENCE.md](docs/API_REFERENCE.md) – REST endpoints, request/response formats, auth.
- **Scoring**: [docs/SCORING.md](docs/SCORING.md) – NOM-035 polarity rules, risk thresholds, domain grouping.
- **Deployment**: [docs/CI_CD.md](docs/CI_CD.md) – Vercel, Portainer, Cloudflare, required secrets.
- **Email**: [docs/SMTP_CONFIGURATION.md](docs/SMTP_CONFIGURATION.md) – SMTP and provider setup.

## Data Files

- **`nom035_questions.json`** – NOM-035 question catalog (138 questions, Guide I/II/III).
- **`risk_strategy.json`** – Risk calculation configuration and thresholds (used by scoring service).

## Project Status

Status snapshot as of **February 16, 2026**:

- Frontend style architecture is consolidated with shared NOM-35 tokens and reusable visual utilities.
- Risk-level presentation now uses a single shared module (`web/frontend/lib/nom35-risk.ts`) across charts, tables, and reports.
- Marketing, login, dashboard/admin, and public assessment flows follow a consistent visual language.
- Visual regression snapshot coverage is in place for key surfaces (`tests/components/visual-regression.test.tsx`).
- Quality gates verified:
  - `cd web/frontend && npm run lint`
  - `cd web/frontend && npm run test -- --run`
  - `cd web/frontend && npm run build`

The project implements a complete NOM-035 compliance solution with:

- **Backend API**: RESTful Go service with PostgreSQL and professional PDF generation
- **Frontend Application**: Modern React/Next.js interface with advanced UX
- **Assessment Engine**: Automated scoring and risk calculation
- **Reporting System**: Individual and company-wide analytics with high-fidelity PDF export
- **User Experience**: Typeform-like assessment interface with full accessibility
- **Security**: JWT authentication with multi-tenant isolation

All core features are implemented and tested. The platform is production-ready for NOM-035-STPS-2018 compliance automation.

## Contributing

1. Fork the repository
2. Create a feature branch from `develop`
3. Make your changes with tests
4. Ensure all tests pass
5. Submit a pull request

### Commit Convention

```
<type>(<scope>): <description>

<body>

<footer>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`

## License

Proprietary software. All rights reserved.
