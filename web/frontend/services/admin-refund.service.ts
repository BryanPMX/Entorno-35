import axios from "axios";
import adminAxiosClient from "@/lib/admin-axios";
import type {
  AdminLoginRequest,
  AdminRefundRequestListResponse,
  AuthResponse,
  RefundRequestRecord,
  ResolveRefundRequestPayload,
  RefundRequestStatus,
} from "@/types/backend";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

class AdminRefundService {
  async login(payload: AdminLoginRequest): Promise<AuthResponse> {
    const response = await axios.post<AuthResponse>("/auth/admin/login", payload, {
      baseURL: API_BASE_URL,
      headers: {
        "Content-Type": "application/json",
      },
    });
    return response.data;
  }

  async listRefundRequests(params?: {
    status?: RefundRequestStatus | "";
    limit?: number;
    offset?: number;
  }): Promise<AdminRefundRequestListResponse> {
    const searchParams = new URLSearchParams();
    if (params?.status) {
      searchParams.set("status", params.status);
    }
    if (typeof params?.limit === "number") {
      searchParams.set("limit", String(params.limit));
    }
    if (typeof params?.offset === "number") {
      searchParams.set("offset", String(params.offset));
    }

    const query = searchParams.toString();
    const endpoint = query ? `/api/v1/admin/billing/refund-requests?${query}` : "/api/v1/admin/billing/refund-requests";
    const response = await adminAxiosClient.get<AdminRefundRequestListResponse>(endpoint);
    return response.data;
  }

  async resolveRefundRequest(requestId: string, payload: ResolveRefundRequestPayload): Promise<RefundRequestRecord> {
    const response = await adminAxiosClient.patch<RefundRequestRecord>(
      `/api/v1/admin/billing/refund-requests/${encodeURIComponent(requestId)}`,
      payload
    );
    return response.data;
  }
}

export const adminRefundService = new AdminRefundService();
export default adminRefundService;
