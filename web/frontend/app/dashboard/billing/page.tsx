"use client";

import { useState } from "react";
import { CreditCard, RefreshCcw, ShieldCheck, ExternalLink } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { translations } from "@/lib/translations";
import { billingService } from "@/services/billing.service";

type FeedbackState =
  | {
      type: "success" | "error" | "info";
      message: string;
    }
  | null;

type LoadingAction = "portal" | "monthly" | "yearly" | null;

export default function BillingPage() {
  const [loadingAction, setLoadingAction] = useState<LoadingAction>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);

  const handleExistingCompanyCheckout = async (plan: "monthly" | "yearly") => {
    setFeedback({
      type: "info",
      message: translations.billing.checkoutRedirect,
    });
    setLoadingAction(plan);

    try {
      const response = await billingService.createExistingCompanyCheckoutSession({ plan });
      if (!response.checkout_url) {
        throw new Error("Checkout URL missing");
      }
      window.location.assign(response.checkout_url);
    } catch (error: unknown) {
      const fallbackMessage = "No se pudo iniciar el checkout de facturacion. Intenta de nuevo.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;

      setFeedback({
        type: "error",
        message: apiMessage,
      });
      setLoadingAction(null);
    }
  };

  const handleOpenPortal = async () => {
    setFeedback({
      type: "info",
      message: translations.billing.portalRedirect,
    });
    setLoadingAction("portal");

    try {
      const response = await billingService.createCustomerPortalSession();
      if (!response.url) {
        throw new Error("Portal URL missing");
      }
      window.location.assign(response.url);
    } catch (error: unknown) {
      const fallbackMessage = "No se pudo abrir el portal de facturacion. Verifica que tu cliente de Stripe este configurado.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;

      setFeedback({
        type: "error",
        message: apiMessage,
      });
      setLoadingAction(null);
    }
  };

  return (
    <div className="page-content-shell page-content-shell-a space-y-8 animate-in fade-in slide-in-from-bottom-4">
      <div className="space-y-2">
        <h1 className="heading-1">{translations.billing.title}</h1>
        <p className="label-muted">{translations.billing.subtitle}</p>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        <Card className="portal-surface border-0">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <ShieldCheck className="h-5 w-5" />
              <span>{translations.billing.portalTitle}</span>
            </CardTitle>
            <CardDescription>{translations.billing.portalDescription}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <p className="text-sm text-muted-foreground">{translations.billing.activePlanHint}</p>
            <Button
              onClick={handleOpenPortal}
              disabled={loadingAction !== null}
              className="w-full sm:w-auto"
            >
              <ExternalLink className="mr-2 h-4 w-4" />
              {loadingAction === "portal" ? translations.billing.portalRedirect : translations.billing.openPortal}
            </Button>
          </CardContent>
        </Card>

        <Card className="portal-surface border-0">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <RefreshCcw className="h-5 w-5" />
              <span>{translations.billing.reactivationTitle}</span>
            </CardTitle>
            <CardDescription>{translations.billing.reactivationDescription}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid gap-3 sm:grid-cols-2">
              <Button
                variant="outline"
                onClick={() => handleExistingCompanyCheckout("monthly")}
                disabled={loadingAction !== null}
                className="justify-start"
              >
                <CreditCard className="mr-2 h-4 w-4" />
                {loadingAction === "monthly" ? translations.billing.checkoutRedirect : translations.billing.monthlyAction}
              </Button>
              <Button
                onClick={() => handleExistingCompanyCheckout("yearly")}
                disabled={loadingAction !== null}
                className="justify-start"
              >
                <CreditCard className="mr-2 h-4 w-4" />
                {loadingAction === "yearly" ? translations.billing.checkoutRedirect : translations.billing.yearlyAction}
              </Button>
            </div>
            <p className="text-sm text-muted-foreground">{translations.billing.supportHint}</p>
          </CardContent>
        </Card>
      </div>

      {feedback ? (
        <div
          className={
            feedback.type === "success"
              ? "rounded-lg border border-emerald-300/50 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-900"
              : feedback.type === "info"
                ? "rounded-lg border border-blue-300/50 bg-blue-500/10 px-4 py-3 text-sm text-blue-900"
                : "rounded-lg border border-destructive/35 bg-destructive/10 px-4 py-3 text-sm text-destructive"
          }
        >
          {feedback.message}
        </div>
      ) : null}
    </div>
  );
}
