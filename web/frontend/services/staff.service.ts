import axiosClient from "@/lib/axios";
import type { ImportResult, PaginatedResponse, Staff, Demographics } from "@/types/backend";

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

export interface CreateStaffRequest {
  full_name: string;
  email?: string;
  curp?: string;
  demographics?: Demographics;
}

export interface UpdateStaffRequest {
  full_name?: string;
  email?: string;
  demographics?: Demographics;
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
   * Create Staff
   * Manually creates a single staff member
   *
   * @param staffData - Staff creation data
   * @returns Promise resolving to created Staff
   * @throws AxiosError on API failure
   */
  async create(staffData: CreateStaffRequest): Promise<Staff> {
    const payload = {
      full_name: staffData.full_name,
      email: staffData.email || undefined,
      curp: staffData.curp || undefined,
      demographics: staffData.demographics || undefined,
    };

    const response = await axiosClient.post<Staff>("/api/v1/staff", payload);
    return response.data;
  }

  /**
   * Update Staff
   * Updates an existing staff member
   *
   * @param id - Staff member UUID
   * @param staffData - Staff update data
   * @returns Promise resolving to updated Staff
   * @throws AxiosError on API failure (404 if not found)
   */
  async update(id: string, staffData: UpdateStaffRequest): Promise<Staff> {
    const payload = {
      full_name: staffData.full_name || undefined,
      email: staffData.email || undefined,
      demographics: staffData.demographics || undefined,
    };

    const response = await axiosClient.put<Staff>(`/api/v1/staff/${id}`, payload);
    return response.data;
  }

  /**
   * Delete Staff
   * Soft deletes a staff member
   *
   * @param id - Staff member UUID
   * @returns Promise resolving when deletion is complete
   * @throws AxiosError on API failure (404 if not found)
   */
  async delete(id: string): Promise<void> {
    await axiosClient.delete(`/api/v1/staff/${id}`);
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

