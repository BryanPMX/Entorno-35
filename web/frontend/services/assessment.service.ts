import axiosClient from "@/lib/axios";
import type { Assessment, PaginatedResponse } from "@/types/backend";

export interface CreateAssessmentRequest {
  staff_id: string;
  period: number;
}

export interface ListAssessmentsParams {
  staff_id?: string;
  period?: number;
  status?: string;
  limit?: number;
  offset?: number;
}

/**
 * Assessment Service
 *
 * Service layer for assessment operations.
 * Components should NEVER call axios directly - they call these service methods.
 * This centralizes error handling and type validation.
 */
class AssessmentService {
  /**
   * Create Assessment
   * Creates a new assessment for a staff member
   *
   * @param data - Assessment creation data
   * @returns Promise resolving to created Assessment
   * @throws AxiosError on API failure
   */
  async create(data: CreateAssessmentRequest): Promise<Assessment> {
    const response = await axiosClient.post<Assessment>("/api/v1/assessments", data);
    return response.data;
  }

  /**
   * List Assessments
   * Lists assessments with optional filters and pagination
   *
   * @param params - Query parameters for filtering and pagination
   * @returns Promise resolving to PaginatedResponse<Assessment>
   * @throws AxiosError on API failure
   */
  async list(params?: ListAssessmentsParams): Promise<PaginatedResponse<Assessment>> {
    const queryParams = new URLSearchParams();

    if (params?.staff_id) queryParams.append("staff_id", params.staff_id);
    if (params?.period) queryParams.append("period", params.period.toString());
    if (params?.status) queryParams.append("status", params.status);
    if (params?.limit) queryParams.append("limit", params.limit.toString());
    if (params?.offset) queryParams.append("offset", params.offset.toString());

    const queryString = queryParams.toString();
    const url = `/api/v1/assessments${queryString ? `?${queryString}` : ""}`;

    const response = await axiosClient.get<PaginatedResponse<Assessment>>(url);
    return response.data;
  }

  /**
   * Get Assessment
   * Retrieves a single assessment by ID
   *
   * @param id - Assessment UUID
   * @returns Promise resolving to Assessment
   * @throws AxiosError on API failure (404 if not found)
   */
  async get(id: string): Promise<Assessment> {
    const response = await axiosClient.get<Assessment>(`/api/v1/assessments/${id}`);
    return response.data;
  }

  /**
   * Create Assessment Link
   * Generates a secure link for staff to access an assessment
   *
   * @param assessmentId - Assessment UUID
   * @param expiresInDays - Link expiration in days (1-365)
   * @returns Promise resolving to link data
   * @throws AxiosError on API failure
   */
  async createLink(assessmentId: string, expiresInDays: number): Promise<{ token: string; link: string; expires_at: string }> {
    const response = await axiosClient.post(`/api/v1/assessments/${assessmentId}/links`, {
      expires_in_days: expiresInDays,
    });
    return response.data;
  }

  /**
   * Calculate Assessment Score
   * Triggers score calculation for a completed assessment
   *
   * @param assessmentId - Assessment UUID
   * @returns Promise resolving to success message
   * @throws AxiosError on API failure
   */
  async calculateScore(assessmentId: string): Promise<{ message: string; assessment_id: string }> {
    const response = await axiosClient.post(`/api/v1/assessments/${assessmentId}/calculate`);
    return response.data;
  }
}

// Export singleton instance
export const assessmentService = new AssessmentService();
export default assessmentService;