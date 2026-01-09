import axiosClient from "@/lib/axios";

export interface RiskDistribution {
  risk_level: string;
  count: number;
}

export interface DepartmentHeatmap {
  department: string;
  risk_level: string;
  count: number;
}

export interface GeneralReportDTO {
  company_id: string;
  company_name?: string;
  period?: number;
  total_staff: number;
  completed_assessments: number;
  participation_rate: number;
  risk_distribution: RiskDistribution[];
  department_heatmap: DepartmentHeatmap[];
}

export interface IndividualReportDTO {
  assessment_id: string;
  period: number;
  guide_type: string;
  staff_name: string;
  department?: string;
  shift?: string;
  total_score: number;
  risk_level: string;
  category_scores: Record<string, number>;
  category_risk_levels: Record<string, string>;
  category_max_scores?: Record<string, number>;
  domain_scores: Record<string, number>;
  domain_risk_levels: Record<string, string>;
  domain_max_scores?: Record<string, number>;
  requires_medical_attention: boolean;
  completed_at?: string;
  recommendations: string[];
}

/**
 * Report Service
 *
 * Service layer for report operations.
 * Components should NEVER call axios directly - they call these service methods.
 * This centralizes error handling and type validation.
 */
class ReportService {
  /**
   * Get Individual Report
   * Retrieves a detailed report for a specific assessment
   *
   * @param assessmentId - Assessment UUID
   * @returns Promise resolving to IndividualReportDTO
   * @throws AxiosError on API failure
   */
  async getIndividualReport(assessmentId: string): Promise<IndividualReportDTO> {
    const response = await axiosClient.get<IndividualReportDTO>(
      `/api/v1/reports/individual/${assessmentId}`
    );
    return response.data;
  }

  /**
   * Get General Report
   * Retrieves aggregated company-wide compliance metrics
   *
   * @param period - Optional period filter (defaults to current)
   * @returns Promise resolving to GeneralReportDTO
   * @throws AxiosError on API failure
   */
  async getGeneralReport(period?: number): Promise<GeneralReportDTO> {
    const queryParams = period ? `?period=${period}` : "";
    const response = await axiosClient.get<GeneralReportDTO>(
      `/api/v1/reports/general${queryParams}`
    );
    return response.data;
  }

  /**
   * Download Individual Report PDF
   * Downloads a PDF version of an individual assessment report
   *
   * @param assessmentId - Assessment UUID
   * @returns Promise resolving to PDF blob
   * @throws AxiosError on API failure
   */
  async downloadIndividualReportPDF(assessmentId: string): Promise<Blob> {
    const response = await axiosClient.get(
      `/api/v1/reports/individual/${assessmentId}/pdf`,
      {
        responseType: 'blob',
      }
    );
    return response.data;
  }
}

// Export singleton instance
export const reportService = new ReportService();
export default reportService;