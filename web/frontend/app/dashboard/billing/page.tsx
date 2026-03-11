"use client";

import { useEffect, useState } from "react";
import { CreditCard, RefreshCcw, ShieldCheck, ExternalLink, FileText } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { translations } from "@/lib/translations";
import { billingService } from "@/services/billing.service";
import type { RefundRequestRecord } from "@/types/backend";

type FeedbackState =
  | {
      type: "success" | "error" | "info";
      message: string;
    }
  | null;

type LoadingAction = "portal" | "monthly" | "yearly" | "refund" | null;

export default function BillingPage() {
  const [loadingAction, setLoadingAction] = useState<LoadingAction>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);
  const [refundReason, setRefundReason] = useState("");
  const [refundHistory, setRefundHistory] = useState<RefundRequestRecord[]>([]);
  const [historyLoading, setHistoryLoading] = useState(true);

  const loadRefundHistory = async () => {
    setHistoryLoading(true);
    try {
      const response = await billingService.getRefundRequests(10, 0);
      setRefundHistory(response.data);
    } catch {
      setRefundHistory([]);
    } finally {
      setHistoryLoading(false);
    }
  };

  useEffect(() => {
    void loadRefundHistory();
  }, []);

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

  const handleCreateRefundRequest = async () => {
    const reason = refundReason.trim();
    if (reason.length < 20) {
      setFeedback({
        type: "error",
        message: "El motivo del reembolso debe tener al menos 20 caracteres.",
      });
      return;
    }

    setFeedback({
      type: "info",
      message: translations.billing.refundSubmitting,
    });
    setLoadingAction("refund");

    try {
      await billingService.createRefundRequest({ reason });
      setRefundReason("");
      await loadRefundHistory();
      setFeedback({
        type: "success",
        message: translations.billing.refundSuccess,
      });
    } catch (error: unknown) {
      const fallbackMessage = "No se pudo enviar la solicitud de reembolso. Intenta de nuevo.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;
      const normalizedMessage =
        apiMessage.toLowerCase().includes("already exists") || apiMessage.toLowerCase().includes("open refund request")
          ? translations.billing.refundDuplicate
          : apiMessage;

      setFeedback({
        type: "error",
        message: normalizedMessage,
      });
    } finally {
      setLoadingAction(null);
    }
  };

  return (
    <div className="page-content-shell page-content-shell-a space-y-8 animate-in fade-in slide-in-from-bottom-4">
      <div className="space-y-2">
        <h1 className="heading-1">{translations.billing.title}</h1>
        <p className="label-muted">{translations.billing.subtitle}</p>
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
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

        <Card className="portal-surface border-0">
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              <span>{translations.billing.refundTitle}</span>
            </CardTitle>
            <CardDescription>{translations.billing.refundDescription}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <label htmlFor="refund-reason" className="text-sm font-medium text-foreground">
              {translations.billing.refundReasonLabel}
            </label>
            <textarea
              id="refund-reason"
              value={refundReason}
              onChange={(event) => setRefundReason(event.target.value)}
              placeholder={translations.billing.refundReasonPlaceholder}
              className="min-h-28 w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground shadow-sm outline-none transition placeholder:text-muted-foreground focus-visible:ring-2 focus-visible:ring-ring"
              maxLength={2000}
              disabled={loadingAction !== null}
            />
            <div className="flex items-center justify-between gap-3">
              <p className="text-xs text-muted-foreground">{refundReason.trim().length}/2000</p>
              <Button
                onClick={handleCreateRefundRequest}
                disabled={loadingAction !== null || refundReason.trim().length < 20}
                className="w-full sm:w-auto"
              >
                {loadingAction === "refund" ? translations.billing.refundSubmitting : translations.billing.refundSubmit}
              </Button>
            </div>
            <div className="space-y-2 rounded-md border border-border/60 bg-background/50 p-3">
              <p className="text-xs font-medium text-muted-foreground">Historial de solicitudes</p>
              {historyLoading ? (
                <p className="text-xs text-muted-foreground">Cargando...</p>
              ) : refundHistory.length === 0 ? (
                <p className="text-xs text-muted-foreground">Sin solicitudes registradas.</p>
              ) : (
                <div className="space-y-2">
                  {refundHistory.map((request) => (
                    <div key={request.id} className="rounded border border-border/50 bg-background/80 px-2 py-1.5 text-xs">
                      <p className="font-medium text-foreground">
                        {request.status.toUpperCase()} - {new Date(request.created_at).toLocaleDateString()}
                      </p>
                      {request.resolution_note ? <p className="text-muted-foreground">{request.resolution_note}</p> : null}
                    </div>
                  ))}
                </div>
              )}
            </div>
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
