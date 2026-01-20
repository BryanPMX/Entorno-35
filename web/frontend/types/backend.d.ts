/**
 * Backend Type Definitions
 * These interfaces strictly match the Go Domain models and API responses
 */

// ============================================================================
// Enums
// ============================================================================

export type SubscriptionStatus = "active" | "inactive";

export type QuestionType = "binary" | "likert";

export type QuestionPolarity = "POSITIVE" | "NEGATIVE";

export type GuideType = "I" | "II" | "III";

export type AssessmentStatus = "pending" | "completed" | "cancelled";

export type RiskLevel = "nulo" | "bajo" | "medio" | "alto" | "muy_alto";

export type SubscriptionStatus = "active" | "inactive" | "canceled" | "past_due";

export type SubscriptionInterval = "month" | "year";

/**
 * Subscription represents a company subscription
 */
export interface Subscription {
  id: string; // UUID
  company_id: string; // UUID
  stripe_subscription_id: string;
  stripe_price_id: string;
  status: SubscriptionStatus;
  interval: SubscriptionInterval;
  current_period_start: string; // ISO timestamp
  current_period_end: string; // ISO timestamp
  cancel_at_period_end: boolean;
  created_at: string; // ISO timestamp
}

/**
 * Question represents a NOM-035 assessment question
 */
export interface Question {
  id: number;
  question_number: number;
  guide_type: GuideType;
  type: QuestionType;
  text: string;
  polarity?: QuestionPolarity;
  section?: string;
  subsection?: string;
  category_id?: number;
  domain_id?: number;
  dimension_id?: number;
  order_index: number;
  created_at: string;
  updated_at: string;
  // Relationships
  category?: {
    id: number;
    name: string;
    created_at: string;
    updated_at: string;
  };
}

// ============================================================================
// Core Models
// ============================================================================

/**
 * Company represents a tenant company in the multi-tenant system
 */
export interface Company {
  id: string; // UUID
  rfc: string;
  name: string;
  address?: string;
  subscription_status: SubscriptionStatus;
  employee_count: number;
  subscription_start_date?: string; // ISO date string
  subscription_end_date?: string; // ISO date string
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
  deleted_at?: string; // ISO timestamp (soft delete)
}

/**
 * Demographics represents flexible demographic data (JSONB)
 */
export interface Demographics {
  gender?: string; // e.g., "masculino", "femenino", "otro"
  age_range?: string; // e.g., "18-25", "26-35", "36-45", "46-55", "56+"
  marital_status?: string; // e.g., "soltero", "casado", "divorciado", "viudo"
  education_level?: string; // e.g., "secundaria", "preparatoria", "universidad", "postgrado"
  time_in_position?: string; // e.g., "<1 año", "1-3 años", "3-5 años", "5+ años"
  shift_type?: string; // e.g., "diurno", "nocturno", "mixto"
  shift_rotation?: string; // e.g., "fijo", "rotativo"
  total_work_experience?: string; // e.g., "<1 año", "1-5 años", "5-10 años", "10+ años"
  department?: string;
  role?: string;
}

/**
 * Staff represents a staff member (user) belonging to a company
 */
export interface Staff {
  id: string; // UUID
  company_id: string; // UUID
  curp?: string; // Made optional
  employee_id?: string | null; // sql.NullString - use string when Valid=true, null when Valid=false
  full_name: string;
  email?: string;
  demographics: Demographics;
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
  deleted_at?: string; // ISO timestamp (soft delete)
}

/**
 * Assessment represents a completed or in-progress assessment
 */
export interface Assessment {
  id: string; // UUID
  staff_id: string; // UUID
  company_id: string; // UUID
  period: number; // e.g., 2025
  guide_type: GuideType;
  status: AssessmentStatus;
  total_score?: number;
  risk_level?: RiskLevel;
  requires_medical_attention: boolean;
  completed_at?: string; // ISO timestamp
  created_at: string; // ISO timestamp
  updated_at: string; // ISO timestamp
  deleted_at?: string; // ISO timestamp (soft delete)

  // Relationships
  staff?: Staff;
  company?: Company;
}

// ============================================================================
// API Request/Response Types
// ============================================================================

/**
 * Login Request
 */
export interface LoginRequest {
  identifier: string; // RFC for COMPANY
  type: "COMPANY";
}

/**
 * Auth Response (Login Response)
 */
export interface AuthResponse {
  token: string; // JWT token
}

/**
 * User context (decoded from JWT, not directly from API)
 * This matches the auth.Context structure in Go
 */
export interface User {
  company_id: string; // UUID
  staff_id?: string; // UUID (optional, only for staff users)
  email?: string;
  role: "company" | "staff";
}

/**
 * Paginated Response wrapper
 */
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  limit: number;
  offset: number;
}

/**
 * Staff Import Result
 */
export interface ImportResult {
  total_processed: number;
  success_count: number;
  skipped_count: number;
  errors: string[];
}

