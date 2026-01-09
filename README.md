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
- PDF export for individual assessment reports
- CSV staff import with CURP validation
- Mobile-responsive assessment interface
- Professional analytics dashboard
- Secure token-based assessments
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
│   ├── lib/              # Utilities and configurations
│   ├── services/         # API service layer
│   └── types/            # TypeScript type definitions
├── docs/                 # API and technical documentation
├── nom035_questions.json # NOM-035 question catalog
├── risk_strategy.json    # Risk calculation configuration
└── docker-compose.yml    # Development environment
```

## Quick Start

### Prerequisites

- Go 1.24+
- Node.js 18+
- Docker & Docker Compose
- Git

### Development Setup

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

   # Frontend dependencies (in separate terminal)
   cd web/frontend && npm install
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the applications**
   ```bash
   # Backend API (terminal 1)
   make run

   # Frontend (terminal 2)
   cd web/frontend && npm run dev
   ```

6. **Verify installation**
   ```bash
   # Backend health check
   curl http://localhost:8080/health

   # Frontend at http://localhost:3000
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

Create a `.env` file based on `.env.example`:

```bash
cp .env.example .env
```

Key configuration options:
- `DB_URL`: PostgreSQL connection string
- `JWT_SECRET`: JWT signing key
- `CORS_ORIGIN`: Frontend URL for CORS
- `REDIS_URL`: Redis connection string

### Testing

```bash
# Run all backend tests
make test

# Run frontend tests
cd web/frontend && npm test

# Run integration tests
go test ./tests/integration/...
```

### API Documentation

Detailed endpoint documentation is available in `docs/`:
- `ASSESSMENT_ENDPOINTS.md` - Assessment management and public endpoints
- `AUTH_ENDPOINTS.md` - Authentication
- `STAFF_ENDPOINTS.md` - Staff management
- `REPORT_ENDPOINTS.md` - Reporting, analytics, and PDF export

## Data Files

- **`nom035_questions.json`** - Complete NOM-035 question catalog (138 questions)
- **`risk_strategy.json`** - Risk calculation configuration and thresholds
- **`docs/Scoring.md`** - Scoring algorithm documentation

## Project Status

This project implements a complete NOM-035 compliance solution with:

- **Backend API**: RESTful Go service with PostgreSQL
- **Frontend Application**: Modern React/Next.js interface
- **Assessment Engine**: Automated scoring and risk calculation
- **Reporting System**: Individual and company-wide analytics with PDF export
- **Security**: JWT authentication with multi-tenant isolation

All core features are implemented and tested. The platform is production-ready for NOM-035-STPS-2018 compliance automation.

For detailed development progress and roadmap, see `PROJECT_STATUS.md`.

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
