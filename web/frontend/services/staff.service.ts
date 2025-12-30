import axiosClient from "@/lib/axios";
import type { ImportResult, PaginatedResponse, Staff } from "@/types/backend";

/**
 * Staff Service
 * 
 * Service layer for staff management operations.
 * Components should NEVER call axios directly - they call these service methods.
 * This centralizes error handling and type validation.
 */

export interface StaffListParams {
  limit?: number; // Default: 50, Max: 100
  offset?: number; // Default: 0
}

class StaffService {
  /**
   * Get All Staff
   * Lists staff members for the authenticated company with pagination
   * 
   * @param params - Pagination parameters (limit, offset)
   * @returns Promise resolving to PaginatedResponse<Staff>
   * @throws AxiosError on API failure
   */
  async getAll(params?: StaffListParams): Promise<PaginatedResponse<Staff>> {
    const queryParams = new URLSearchParams();
    
    if (params?.limit !== undefined) {
      queryParams.append("limit", params.limit.toString());
    }
    if (params?.offset !== undefined) {
      queryParams.append("offset", params.offset.toString());
    }

    const queryString = queryParams.toString();
    const url = `/api/v1/staff${queryString ? `?${queryString}` : ""}`;

    const response = await axiosClient.get<PaginatedResponse<Staff>>(url);
    return response.data;
  }

  /**
   * Get Staff by ID
   * Retrieves a single staff member by ID
   * 
   * @param id - Staff member UUID
   * @returns Promise resolving to Staff
   * @throws AxiosError on API failure (404 if not found)
   */
  async getById(id: string): Promise<Staff> {
    const response = await axiosClient.get<Staff>(`/api/v1/staff/${id}`);
    return response.data;
  }

  /**
   * Upload CSV
   * Bulk imports staff members from a CSV file
   * 
   * @param file - CSV file to upload
   * @returns Promise resolving to ImportResult with import statistics
   * @throws AxiosError on API failure
   */
  async uploadCSV(file: File): Promise<ImportResult> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await axiosClient.post<ImportResult>(
      "/api/v1/staff/import",
      formData,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      }
    );

    return response.data;
  }
}

// Export singleton instance
export const staffService = new StaffService();
export default staffService;

