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

#### Quick Start (Recommended)

Use the automated startup script for the easiest setup:

```bash
# Clone and enter repository
git clone <repository-url>
cd Entorno-35

# Install dependencies
make setup
cd web/frontend && npm install && cd ../..

# Start everything (one command)
./start-dev.sh
```

This will automatically:
- Start PostgreSQL and Redis containers
- Start backend API on http://localhost:8080
- Start frontend on http://localhost:3000
- Configure all environment variables

Press `Ctrl+C` to stop all services.

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

Comprehensive API documentation available in `docs/`:
- `COMPLETE_API_REFERENCE.md` - Complete API endpoint reference
- `Scoring.md` - NOM-035 scoring algorithm and thresholds
- `CSV_IMPORT_VALIDATION.md` - Staff CSV import specifications
- `BUGFIXES_JAN2026.md` - Recent bug fixes and known issues

## Data Files

- **`nom035_questions.json`** - Complete NOM-035 question catalog (138 questions)
- **`risk_strategy.json`** - Risk calculation configuration and thresholds
- **`docs/Scoring.md`** - Scoring algorithm documentation

## Project Status

This project implements a complete NOM-035 compliance solution with:

- **Backend API**: RESTful Go service with PostgreSQL and professional PDF generation
- **Frontend Application**: Modern React/Next.js interface with advanced UX
- **Assessment Engine**: Automated scoring and risk calculation
- **Reporting System**: Individual and company-wide analytics with high-fidelity PDF export
- **User Experience**: Typeform-like assessment interface with full accessibility
- **Security**: JWT authentication with multi-tenant isolation

All core features are implemented and tested. The platform features enterprise-grade UX with professional PDF reporting and is production-ready for NOM-035-STPS-2018 compliance automation.

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
