# Entorno35 Frontend

Next.js 16 frontend for the Entorno35 NOM-035 Compliance Platform. Designed to run as a self-hosted application: you build and serve the static assets and point the API URL at your own backend.

For API reference, deployment (Vercel, Portainer), and platform documentation, see the root [docs/README.md](../../docs/README.md).

## Project Structure

```
web/frontend/
├── app/                        # Next.js App Router
│   ├── layout.tsx              # Root layout (fonts, QueryProvider, Toaster)
│   ├── page.tsx                # Marketing home
│   ├── globals.css             # Global styles and theme
│   ├── (auth)/                 # Auth routes
│   │   └── login/page.tsx      # Company login (RFC)
│   │   └── register/           # Registration + Stripe checkout confirmation
│   ├── (marketing)/layout.tsx  # Marketing layout wrapper
│   ├── (public)/assessment/    # Public assessment by token
│   │   └── [token]/            # Token-based assessment flow
│   ├── dashboard/              # Authenticated dashboard
│   │   ├── layout.tsx          # Dashboard layout and nav
│   │   ├── page.tsx            # Dashboard home
│   │   ├── assessments/        # Assessment list and report
│   │   └── staff/page.tsx      # Staff management
│   └── debug/connection/       # API connection check
├── components/
│   ├── assessments/           # Assessment wizard, table, delete dialog
│   ├── dashboard/              # Charts, metric cards, heatmaps
│   ├── layout/auth-guard.tsx   # Route protection
│   ├── marketing/              # Hero, features, compliance, pricing, nav, footer
│   ├── providers/query-provider.tsx
│   ├── staff/                  # Staff table, CSV upload, CRUD dialogs
│   └── ui/                     # Shadcn UI primitives (button, card, dialog, etc.)
├── lib/
│   ├── axios.ts                # Axios client and interceptors
│   ├── jwt.ts                  # JWT decode/session helpers (library-based decode + exp checks)
│   ├── nom35-risk.ts           # Shared NOM-35 risk labels/colors/badge classes
│   ├── query-client.ts         # TanStack Query config
│   ├── store/auth-store.ts     # Zustand auth state
│   ├── utils.ts                # cn() and helpers
│   └── translations.ts
├── services/                   # API service layer
│   ├── auth.service.ts
│   ├── assessment.service.ts
│   ├── staff.service.ts
│   └── report.service.ts
├── types/
│   └── backend.d.ts            # Backend API types (aligned with Go)
├── public/                     # Static assets
├── tests/                      # Vitest setup and tests
└── package.json
```

## Architecture

### Design Principles

1. **Service layer**: Components do not call Axios directly. All API calls go through services (e.g. `authService.login()`, `staffService.getAll()`).
2. **Type safety**: TypeScript interfaces in `types/backend.d.ts` match the Go backend domain models.
3. **Data fetching**: TanStack Query handles caching, cancellation, and refetching.
4. **Auth**: Axios interceptors attach the JWT, clear auth state on 401, and trigger client-side navigation to `/login` (no full page reload).
5. **Design system consistency**: `app/globals.css` provides shared visual tokens and reusable UI utilities for auth, portal, assessment, and marketing surfaces.
6. **Single source of truth for risk UI**: `lib/nom35-risk.ts` centralizes NOM-35 risk labels, colors, thresholds, and badge classes used by charts and tables.

### Technologies

- Next.js 16 (App Router)
- React 19, TypeScript
- TanStack Query, Axios
- Shadcn UI (Radix), Tailwind CSS, Lucide React
- Zustand (auth store), Framer Motion (assessment UX)

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend API running (e.g. `http://localhost:8080` in development)

### Installation

```bash
npm install
```

### Environment

Create `.env.local` in `web/frontend/`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

If unset, the client defaults to `http://localhost:8080`. In production, set this to your self-hosted API base URL.

### Development

```bash
npm run dev
```

Open http://localhost:3000. The dev server proxies nothing; all API calls use `NEXT_PUBLIC_API_URL`.

### Build and production

```bash
npm run build
npm start
```

For self-hosted deployment, build on your server or in CI and serve the output (e.g. `npm start` for Node, or a static export if you switch to static generation). Ensure `NEXT_PUBLIC_API_URL` points to your backend.

## Service Layer

### Auth (`services/auth.service.ts`)

- `authService.login({ identifier, type: 'COMPANY', password })` — returns JWT token
- `authService.isAuthenticated()` — whether a non-expired valid token is present
- `authService.logout()` — clears token and state

### Billing (`services/billing.service.ts`)

- `billingService.createCheckoutSession(payload)` — creates Stripe checkout session for paid registration
- `billingService.verifyCheckoutSession(sessionId)` — confirms payment and account activation

### Staff (`services/staff.service.ts`)

- `staffService.getAll({ limit, offset })` — paginated list
- `staffService.getById(id)` — single member
- `staffService.uploadCSV(file)` — CSV import with validation and progress

### Axios client (`lib/axios.ts`)

- Injects `Authorization: Bearer <token>` from localStorage
- On 401, clears auth state and redirects client-side to `/login` (no full page reload)
- Base URL from `NEXT_PUBLIC_API_URL`

### TanStack Query

- `staleTime: 5 minutes`, `refetchOnWindowFocus: false`, `retry: 1`
- Request cancellation on unmount

## Type Definitions

`types/backend.d.ts` defines domain and API types aligned with the Go backend: `Company`, `Staff`, `Assessment`, `AuthResponse`, `PaginatedResponse<T>`, enums such as `RiskLevel`, `AssessmentStatus`, etc.

## Scripts

| Script              | Description                |
|---------------------|----------------------------|
| `npm run dev`       | Start development server   |
| `npm run build`     | Production build          |
| `npm start`         | Run production server     |
| `npm run lint`      | Run ESLint                |
| `npm test`          | Run Vitest                |
| `npm run test:coverage` | Vitest with coverage  |

## Testing and Quality Gates

- **Unit/component tests**: Vitest + Testing Library
- **Visual regression snapshots**: `tests/components/visual-regression.test.tsx`
  - Login page shell
  - Dashboard layout shell
  - Public assessment shell
- **Linting**: ESLint (`npm run lint`)
- **Production build validation**: Next build (`npm run build`)

## Development Phases (Completed)

- **Architecture and service layer**: Next.js setup, types, Axios interceptors, TanStack Query, service pattern.
- **Authentication UI**: Login (company RFC + password), Stripe registration flow, Zustand auth store, AuthGuard, dashboard layout.
- **Staff management**: Staff table (pagination, sort, filter), CSV upload (drag-and-drop, progress, errors), CRUD dialogs, template download.
- **Dashboard and analytics**: Metrics, risk distribution chart, department heatmap, assessment wizard, empty states, responsive layout.
- **Public assessment**: Token-based assessment UI, form navigation, progress, mobile-friendly, completion screen.
- **PDF export**: Report PDF download, NOM-035 layout, JWT-protected endpoint.
- **Assessment UX**: Framer Motion transitions, category badges, keyboard navigation (arrows, 1–5, Enter), auto-save indicators, accessibility.
- **Layout and marketing**: Focus-mode assessment layout, marketing landing (hero, features, compliance, pricing), subscription plans (monthly/yearly MXN), navigation and footer.
- **Design system consolidation (2026-02-16)**: Shared visual tokenized language across login, admin portal, public assessment, and marketing, preserving centralized NOM-35 risk semantics.
- **Risk style consolidation (2026-02-16)**: Centralized risk-level registry for color/label/badge mapping via `lib/nom35-risk.ts`.
- **Visual regression baseline (2026-02-16)**: Snapshot coverage added for key shells to prevent style drift.

## Current Status (February 16, 2026)

- Frontend compiles for production and passes all configured tests.
- Verified commands:
  - `npm run lint`
  - `npm run test -- --run`
  - `npm run build`
- Current Vitest scope includes behavior tests plus snapshot-based visual regression for critical routes.

## License

Proprietary.
