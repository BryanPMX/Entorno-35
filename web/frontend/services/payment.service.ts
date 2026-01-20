import axiosClient from "@/lib/axios";
import type { Subscription } from "@/types/backend";

/**
 * Payment Service
 *
 * Service layer for payment and subscription operations.
 * Components should NEVER call axios directly - they call these service methods.
 * This centralizes error handling and type validation.
 */
class PaymentService {
  /**
   * Create Subscription
   * Creates a new subscription for the company
   *
   * @param priceId - Stripe price ID (monthly or yearly)
   * @returns Promise resolving to subscription data
   * @throws AxiosError on payment failure
   */
  async createSubscription(priceId: string): Promise<Subscription> {
    const response = await axiosClient.post<Subscription>("/payments/subscription", {
      price_id: priceId,
    });

    return response.data;
  }

  /**
   * Get Subscription
   * Retrieves the current subscription for the authenticated company
   *
   * @returns Promise resolving to current subscription
   * @throws AxiosError if no subscription exists
   */
  async getSubscription(): Promise<Subscription> {
    const response = await axiosClient.get<Subscription>("/payments/subscription");
    return response.data;
  }

  /**
   * Cancel Subscription
   * Cancels the current subscription at period end
   *
   * @returns Promise resolving when cancellation is complete
   * @throws AxiosError on cancellation failure
   */
  async cancelSubscription(): Promise<void> {
    await axiosClient.delete("/payments/subscription");
  }
}

// Export singleton instance
export const paymentService = new PaymentService();
export default paymentService;