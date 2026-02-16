# Documentation

This directory contains all substantive technical and operational documentation for the Entorno35 NOM-035 Compliance Platform. Use this index to find the right document for your task.

**Where documentation lives**

- **In this directory (`docs/`)**: All non-README documentation (CI/CD, API reference, Scoring, SMTP). Nothing else.
- **At repo root**: [README.md](../README.md) (project overview, quick start) and [SECURITY.md](../SECURITY.md) (credentials, pre-commit). SECURITY.md is at root so GitHub can link it as the repository Security policy.
- **Next to code**: Package READMEs only (e.g. `cmd/README.md`, `internal/README.md`, `web/frontend/README.md`) for structure and usage; they are not duplicate docs.

**Essential docs (only these)**

| Location | Document | Purpose |
|----------|----------|---------|
| Root | README.md | Project overview, quick start, commands |
| Root | SECURITY.md | Credentials, pre-commit check, required env vars |
| docs/ | README.md (this file) | Documentation index |
| docs/ | CI_CD.md | Deployment (Vercel, Portainer, Cloudflare) |
| docs/ | API_REFERENCE.md | REST API reference |
| docs/ | SCORING.md | NOM-035 scoring rules and thresholds |
| docs/ | SMTP_CONFIGURATION.md | Email configuration |

Package READMEs under `cmd/`, `internal/`, `pkg/`, `migrations/`, `tests/`, `web/frontend/` are for code navigation only; they are not part of the docs set above.

---

## Current Status Snapshot (February 16, 2026)

- Project state is production-ready for the currently implemented feature set.
- Frontend architecture now includes:
  - Shared visual token system in `web/frontend/app/globals.css` for consistent brand surfaces.
  - Centralized NOM-35 risk style registry in `web/frontend/lib/nom35-risk.ts` for charts, reports, and badges.
  - Visual regression snapshots in `web/frontend/tests/components/visual-regression.test.tsx`.
- Latest verified frontend quality gates:
  - `cd web/frontend && npm run lint`
  - `cd web/frontend && npm run test -- --run`
  - `cd web/frontend && npm run build`

For the most current setup and status details:

- Root project status: [README.md](../README.md)
- Frontend implementation status: [web/frontend/README.md](../web/frontend/README.md)

---

## Overview

| Document | Description |
|----------|-------------|
| [CI_CD.md](CI_CD.md) | Continuous integration and deployment: Vercel (frontend), GitHub Actions, Portainer (backend), Cloudflare. Production branch, required secrets, and step-by-step Portainer and Cloudflare checks. |
| [API_REFERENCE.md](API_REFERENCE.md) | REST API reference: authentication, staff, assessments, scoring, and reports. Request/response formats, status codes, and query parameters. |
| [SCORING.md](SCORING.md) | NOM-035 scoring implementation: polarity rules (Guide II and III), risk level thresholds, domain grouping, and verified mapping to the official NOM-035-STPS-2018 document. |
| [SMTP_CONFIGURATION.md](SMTP_CONFIGURATION.md) | Email configuration for assessment invitations: required env vars, Gmail and SendGrid setup, production providers, and troubleshooting. |

---

## By Role

**Deploying or operating the platform**

- Start with [CI_CD.md](CI_CD.md) for Vercel, Portainer, and Cloudflare.
- Backend env vars are defined in the repo root `docker-compose.prod.yml`; set their values in Portainer.
- See root [SECURITY.md](../SECURITY.md) for credentials, pre-commit checks, and required environment variables.

**Integrating with the API**

- Use [API_REFERENCE.md](API_REFERENCE.md) for endpoints, payloads, and auth.
- Use [SCORING.md](SCORING.md) if you need risk level calculation or domain logic.

**Configuring email**

- Use [SMTP_CONFIGURATION.md](SMTP_CONFIGURATION.md) for SMTP and provider setup.

---

## Data and Configuration Files (repo root)

- `nom035_questions.json` – NOM-035 question catalog (138 questions, Guide I/II/III).
- `risk_strategy.json` – Risk calculation configuration and thresholds used by the scoring service.

---

## Package and Module Documentation

Package-level READMEs live next to the code:

- **Backend**: `cmd/README.md`, `internal/README.md`, and subpackages under `internal/` (adapters, config, core, domain, etc.).
- **Frontend**: `web/frontend/README.md` – Next.js app structure, services, and scripts.
- **Migrations**: `migrations/README.md` – Schema migrations and usage.
- **Tests**: `tests/integration/README.md` – Integration test setup and scenarios.

For project setup, quick start, and high-level architecture, see the root [README.md](../README.md).
