import axiosClient from "@/lib/axios";
import type {
  CreateCheckoutSessionRequest,
  CreateCheckoutSessionResponse,
  CreateCustomerPortalSessionResponse,
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
}

export const billingService = new BillingService();
export default billingService;
