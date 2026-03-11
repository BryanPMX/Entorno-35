import axiosClient from "@/lib/axios";
import type {
  CreateCheckoutSessionRequest,
  CreateCheckoutSessionResponse,
  CreateCustomerPortalSessionResponse,
  PaginatedResponse,
  CreateRefundRequestRequest,
  CreateRefundRequestResponse,
  RefundRequestRecord,
  ExistingCompanyCheckoutSessionRequest,
  VerifyCheckoutSessionResponse,
} from "@/types/backend";

class BillingService {
  async createCheckoutSession(payload: CreateCheckoutSessionRequest): Promise<CreateCheckoutSessionResponse> {
    const response = await axiosClient.post<CreateCheckoutSessionResponse>("/billing/checkout-session", payload);
    return response.data;
  }

  async verifyCheckoutSession(sessionId: string): Promise<VerifyCheckoutSessionResponse> {
    const response = await axiosClient.get<VerifyCheckoutSessionResponse>(`/billing/checkout-session/${sessionId}/verify`);
    return response.data;
  }

  async createExistingCompanyCheckoutSession(
    payload: ExistingCompanyCheckoutSessionRequest
  ): Promise<CreateCheckoutSessionResponse> {
    const response = await axiosClient.post<CreateCheckoutSessionResponse>("/api/v1/billing/checkout-session", payload);
    return response.data;
  }

  async createCustomerPortalSession(): Promise<CreateCustomerPortalSessionResponse> {
    const response = await axiosClient.post<CreateCustomerPortalSessionResponse>("/api/v1/billing/customer-portal", {});
    return response.data;
  }

  async createRefundRequest(payload: CreateRefundRequestRequest): Promise<CreateRefundRequestResponse> {
    const response = await axiosClient.post<CreateRefundRequestResponse>("/api/v1/billing/refund-request", payload);
    return response.data;
  }

  async getRefundRequests(limit = 20, offset = 0): Promise<PaginatedResponse<RefundRequestRecord>> {
    const params = new URLSearchParams({
      limit: String(limit),
      offset: String(offset),
    });
    const response = await axiosClient.get<PaginatedResponse<RefundRequestRecord>>(`/api/v1/billing/refund-requests?${params.toString()}`);
    return response.data;
  }
}

export const billingService = new BillingService();
export default billingService;
