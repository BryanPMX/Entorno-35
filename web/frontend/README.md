# Entorno35 Frontend

Next.js 14+ frontend application for the Entorno35 NOM-035 Compliance Platform.

## Project Structure

```
web/frontend/
├── app/                   # Next.js App Router
│   ├── layout.tsx        # Root layout with QueryProvider
│   ├── page.tsx          # Home page
│   └── globals.css       # Global styles
├── components/            # React components
│   ├── providers/        # Context providers
│   │   └── query-provider.tsx  # TanStack Query provider
│   └── ui/               # Shadcn UI components (to be added)
├── lib/                   # Utilities and configurations
│   ├── axios.ts          # Axios client with interceptors
│   ├── query-client.ts   # TanStack Query configuration
│   └── utils.ts          # Utility functions (cn, etc.)
├── services/              # API service layer
│   ├── auth.service.ts   # Authentication service
│   ├── staff.service.ts  # Staff management service
│   └── report.service.ts # Report and PDF export service
├── types/                 # TypeScript type definitions
│   └── backend.d.ts      # Backend API types (matches Go models)
├── public/                # Static assets
└── package.json          # Dependencies and scripts
```

## Architecture

### Design Principles

1. **Service Layer Pattern**: Components never call `axios` directly. All API calls go through service methods (`authService.login()`, `staffService.getAll()`, etc.)
2. **Type Safety**: Strict TypeScript interfaces matching the Go backend domain models
3. **Race Condition Prevention**: TanStack Query handles request cancellation, caching, and background re-fetching
4. **Auth Sync**: Axios interceptors automatically inject tokens and handle 401 redirects

### Key Technologies

- **Next.js 14+** (App Router) - React framework
- **TypeScript** - Type safety
- **TanStack Query (React Query)** - Data fetching and state management
- **Axios** - HTTP client
- **Shadcn UI** - Component library
- **Tailwind CSS** - Styling
- **Lucide React** - Icons

## Getting Started

### Prerequisites

- Node.js 18+ and npm
- Backend API running on `http://localhost:8080` (or configure via environment variable)

### Installation

```bash
# Install dependencies
npm install
```

### Environment Configuration

Create a `.env.local` file in the root of `web/frontend/`:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

If not set, defaults to `http://localhost:8080`.

### Development

```bash
# Start development server
npm run dev

# Open http://localhost:3000
```

### Build

```bash
# Build for production
npm run build

# Start production server
npm start
```

## Service Layer

All API interactions go through the service layer:

### Auth Service (`services/auth.service.ts`)

```typescript
import { authService } from '@/services/auth.service';

// Login
const response = await authService.login({
  identifier: 'ABC123456789',
  type: 'COMPANY'
});

// Check authentication status
const isAuth = authService.isAuthenticated();

// Logout
authService.logout();
```

### Staff Service (`services/staff.service.ts`)

```typescript
import { staffService } from '@/services/staff.service';

// Get paginated staff list
const staff = await staffService.getAll({ limit: 50, offset: 0 });

// Get single staff member
const member = await staffService.getById('uuid-here');

// Upload CSV
const result = await staffService.uploadCSV(file);
```

**CSV Upload Features:**
- Drag-and-drop file upload
- Progress tracking and error reporting
- Template download functionality
- Validation and success metrics
- Batch processing with transaction safety

## Axios Client

The Axios client (`lib/axios.ts`) is configured with:

- **Automatic token injection**: Reads token from `localStorage` and adds `Authorization: Bearer <token>` header
- **401 handling**: Automatically redirects to `/login` on 401 Unauthorized responses
- **Base URL configuration**: Configurable via `NEXT_PUBLIC_API_URL` environment variable

## TanStack Query

TanStack Query is configured with:

- `refetchOnWindowFocus: false` - Prevents refetch on tab focus
- `staleTime: 5 minutes` - Data is fresh for 5 minutes
- `retry: 1` - Retry once on failure
- Automatic request cancellation on unmount

### Usage Example

```typescript
'use client';

import { useQuery } from '@tanstack/react-query';
import { staffService } from '@/services/staff.service';

export function StaffList() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['staff', { limit: 50, offset: 0 }],
    queryFn: () => staffService.getAll({ limit: 50, offset: 0 })
  });

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error: {error.message}</div>;

  return <div>{/* Render staff list */}</div>;
}
```

## Type Definitions

All backend types are defined in `types/backend.d.ts`:

- `Company`, `Staff`, `Assessment` - Domain models
- `AuthResponse`, `LoginRequest` - API request/response types
- `PaginatedResponse<T>`, `ImportResult` - Utility types
- Enums: `RiskLevel`, `AssessmentStatus`, `SubscriptionStatus`, etc.

These types strictly match the Go backend domain models for type safety across the stack.

## Development Phases

### Phase 5.1: Architecture & Service Layer ✅ COMPLETE

- Next.js project setup
- TypeScript type definitions
- Axios client with interceptors
- Service layer pattern
- TanStack Query configuration

### Phase 5.2: Authentication UI ✅ COMPLETE

- Login screen with Company/Staff type support
- Zustand auth store with global state management
- Route protection (AuthGuard component)
- Dashboard layout with navigation
- JWT token decoding and session management

### Phase 5.3: Staff Management UI ✅ COMPLETE

- Staff data table with pagination, sorting, and filtering
- CSV uploader with drag-and-drop, progress indication, and error handling
- Complete staff management page integration
- UI components: table, pagination, dialog, alert, progress
- Template download functionality
- All tests passing (20/20 frontend tests)

### Phase 5.4: Dashboard & Analytics ✅ COMPLETE

- Professional admin dashboard with metrics and interactive charts
- Risk distribution donut chart with hover effects and tooltips
- Department heatmap bar chart with risk-based color coding
- Assessment creation wizard with multi-step modal
- Recharts integration for data visualization
- Sophisticated animations and professional UI/UX
- Empty state handling with onboarding cards
- Responsive design and mobile-friendly interface

### Phase 5.5: Public Assessment View ✅ COMPLETE

- Minimalist staff assessment interface with distraction-free design
- Form-based question navigation with auto-advance
- Progress bar with completion tracking
- Mobile-first, touch-optimized interface
- Secure token-based authentication
- Professional completion screen

### Phase 6.0: PDF Export ✅ COMPLETE

- PDF download button in assessment report pages
- Server-side PDF generation with NOM-035 compliant formatting
- Professional report layout with all assessment data
- Automatic filename generation and browser download
- Secure JWT-authenticated PDF access

### Phase 6.1: High-Fidelity PDF Reports ✅ COMPLETE

- Brand header with dark slate banner and CONFIDENTIAL labeling
- Professional typography hierarchy (16pt title, 12pt subtitle, 10pt body)
- Visual risk thermometer with color-coded markers (Green/Yellow/Red)
- Professional data tables with zebra striping for category/domain scores
- Recommendations checklist with warning icons for high-risk assessments
- Footer with Entorno35 branding and page numbering
- UTF-8 support with DejaVu Sans fonts for Spanish characters

## Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm start` - Start production server
- `npm run lint` - Run ESLint

## License

Proprietary
