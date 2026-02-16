import axiosClient from "@/lib/axios";
import type { Assessment, Question, PaginatedResponse } from "@/types/backend";

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

  /**
   * Fetch Public Assessment by Token
   * Retrieves assessment details and questions for public access
   *
   * @param token - Assessment link token
   * @returns Promise resolving to assessment and questions
   * @throws AxiosError on API failure
   */
  async fetchByToken(token: string): Promise<{ assessment: Assessment; questions: Question[] }> {
    const response = await axiosClient.get(`/api/v1/assessments/public/${token}`);
    const payload = response.data as { assessment?: Assessment; questions?: Question[] };

    if (!payload?.assessment) {
      throw new Error("Assessment not found for this link.");
    }

    if (!Array.isArray(payload.questions) || payload.questions.length === 0) {
      throw new Error("No questions are configured for this assessment.");
    }

    return {
      assessment: payload.assessment,
      questions: payload.questions,
    };
  }

  /**
   * Submit Assessment Responses
   * Submits staff responses for a public assessment
   *
   * @param token - Assessment link token
   * @param responses - Array of question responses
   * @returns Promise resolving to success message
   * @throws AxiosError on API failure
   */
  async submitResponses(token: string, responses: { question_id: number; value: number }[]): Promise<{ message: string; submitted_at: string }> {
    const response = await axiosClient.post(`/api/v1/assessments/public/${token}/submit`, {
      responses,
    });
    return response.data;
  }

  /**
   * Delete Assessment
   * Deletes a pending assessment (only pending assessments can be deleted)
   *
   * @param id - Assessment UUID
   * @returns Promise resolving when deletion is complete
   * @throws AxiosError on API failure (404 if not found, 400 if not pending)
   */
  async delete(id: string): Promise<void> {
    await axiosClient.delete(`/api/v1/assessments/${id}`);
  }

  /**
   * Send Assessment Email
   * Sends an assessment link via email to the staff member
   * The email address is securely retrieved from the staff member's record in the database
   *
   * @param assessmentId - Assessment UUID
   * @returns Promise resolving to success message
   * @throws AxiosError on API failure
   */
  async sendEmail(assessmentId: string): Promise<{ message: string }> {
    const response = await axiosClient.post(`/api/v1/assessments/${assessmentId}/send-email`);
    return response.data;
  }
}

// Export singleton instance
export const assessmentService = new AssessmentService();
export default assessmentService;
