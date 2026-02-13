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
│   ├── jwt.ts                  # JWT decode/session helpers
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
4. **Auth**: Axios interceptors attach the JWT and redirect to `/login` on 401.

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

- `authService.login({ identifier, type: 'COMPANY' })` — returns JWT and user info
- `authService.isAuthenticated()` — whether a valid token is present
- `authService.logout()` — clears token and state

### Staff (`services/staff.service.ts`)

- `staffService.getAll({ limit, offset })` — paginated list
- `staffService.getById(id)` — single member
- `staffService.uploadCSV(file)` — CSV import with validation and progress

### Axios client (`lib/axios.ts`)

- Injects `Authorization: Bearer <token>` from localStorage
- On 401, redirects to `/login`
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

## Development Phases (Completed)

- **Architecture and service layer**: Next.js setup, types, Axios interceptors, TanStack Query, service pattern.
- **Authentication UI**: Login (company RFC), Zustand auth store, AuthGuard, dashboard layout.
- **Staff management**: Staff table (pagination, sort, filter), CSV upload (drag-and-drop, progress, errors), CRUD dialogs, template download.
- **Dashboard and analytics**: Metrics, risk distribution chart, department heatmap, assessment wizard, empty states, responsive layout.
- **Public assessment**: Token-based assessment UI, form navigation, progress, mobile-friendly, completion screen.
- **PDF export**: Report PDF download, NOM-035 layout, JWT-protected endpoint.
- **Assessment UX**: Framer Motion transitions, category badges, keyboard navigation (arrows, 1–5, Enter), auto-save indicators, accessibility.
- **Layout and marketing**: Focus-mode assessment layout, marketing landing (hero, features, compliance, pricing), subscription plans (monthly/yearly MXN), navigation and footer.
- **Marketing visual refresh (2026-02)**: Glassmorphism accents, gradient badges, stat highlight cards, soft grid backgrounds, and animated hero metrics to improve clarity and conversion.

## License

Proprietary.
