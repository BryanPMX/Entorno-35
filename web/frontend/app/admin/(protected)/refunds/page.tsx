"use client";

import { useEffect, useMemo, useState } from "react";
import { RefreshCcw } from "lucide-react";
import { adminRefundService } from "@/services/admin-refund.service";
import type { RefundRequestRecord, RefundRequestStatus } from "@/types/backend";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

type FeedbackState = {
  type: "error" | "success";
  message: string;
} | null;

const PAGE_SIZE = 20;

function statusLabel(status: RefundRequestStatus): string {
  switch (status) {
    case "requested":
      return "Solicitado";
    case "approved":
      return "Aprobado";
    case "rejected":
      return "Rechazado";
    case "refunded":
      return "Reembolsado";
    default:
      return status;
  }
}

function statusBadgeClass(status: RefundRequestStatus): string {
  switch (status) {
    case "requested":
      return "border-amber-300/60 bg-amber-500/10 text-amber-900";
    case "approved":
      return "border-blue-300/60 bg-blue-500/10 text-blue-900";
    case "rejected":
      return "border-red-300/60 bg-red-500/10 text-red-900";
    case "refunded":
      return "border-emerald-300/60 bg-emerald-500/10 text-emerald-900";
    default:
      return "border-border bg-background text-foreground";
  }
}

export default function AdminRefundRequestsPage() {
  const [requests, setRequests] = useState<RefundRequestRecord[]>([]);
  const [total, setTotal] = useState(0);
  const [offset, setOffset] = useState(0);
  const [statusFilter, setStatusFilter] = useState<RefundRequestStatus | "">("requested");
  const [notes, setNotes] = useState<Record<string, string>>({});
  const [refundIDs, setRefundIDs] = useState<Record<string, string>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [actionLoadingID, setActionLoadingID] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);

  const hasPrev = offset > 0;
  const hasNext = offset + requests.length < total;
  const pageLabel = useMemo(() => `${offset + 1}-${Math.min(offset + requests.length, total)} de ${total}`, [offset, requests.length, total]);

  const loadRequests = async (nextOffset = offset) => {
    setIsLoading(true);
    setFeedback(null);
    try {
      const response = await adminRefundService.listRefundRequests({
        status: statusFilter,
        limit: PAGE_SIZE,
        offset: nextOffset,
      });
      setRequests(response.data);
      setTotal(response.total);
      setOffset(response.offset);
    } catch (error: unknown) {
      const fallbackMessage = "No se pudieron cargar las solicitudes de reembolso.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;
      setFeedback({ type: "error", message: apiMessage });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    void loadRequests(0);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [statusFilter]);

  const runDecision = async (request: RefundRequestRecord, target: Extract<RefundRequestStatus, "approved" | "rejected" | "refunded">) => {
    const note = (notes[request.id] ?? "").trim();
    const stripeRefundID = (refundIDs[request.id] ?? "").trim();

    if ((target === "rejected" || target === "refunded") && note.length < 10) {
      setFeedback({ type: "error", message: "Agrega una nota de resolución de al menos 10 caracteres." });
      return;
    }
    if (target === "refunded" && stripeRefundID === "") {
      setFeedback({ type: "error", message: "Debes indicar el ID de reembolso de Stripe para marcar como reembolsado." });
      return;
    }

    setActionLoadingID(request.id);
    setFeedback(null);
    try {
      await adminRefundService.resolveRefundRequest(request.id, {
        status: target,
        resolution_note: note || undefined,
        stripe_refund_id: target === "refunded" ? stripeRefundID : undefined,
      });
      setFeedback({
        type: "success",
        message: `Solicitud ${request.id.slice(0, 8)} actualizada a ${statusLabel(target).toLowerCase()}.`,
      });
      await loadRequests(offset);
    } catch (error: unknown) {
      const fallbackMessage = "No se pudo actualizar la solicitud.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;
      setFeedback({ type: "error", message: apiMessage });
    } finally {
      setActionLoadingID(null);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 className="heading-1">Solicitudes de Reembolso</h1>
          <p className="label-muted">Panel interno para revisión y resolución de reembolsos.</p>
        </div>
        <Button type="button" variant="outline" onClick={() => void loadRequests(offset)} disabled={isLoading || actionLoadingID !== null}>
          <RefreshCcw className="mr-2 h-4 w-4" />
          Actualizar
        </Button>
      </div>

      <Card className="portal-surface border-0">
        <CardHeader className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <CardTitle>Bandeja</CardTitle>
            <CardDescription>Filtra y gestiona solicitudes pendientes o históricas.</CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <label htmlFor="status-filter" className="text-sm text-muted-foreground">
              Estado
            </label>
            <select
              id="status-filter"
              value={statusFilter}
              onChange={(event) => setStatusFilter(event.target.value as RefundRequestStatus | "")}
              className="h-10 rounded-md border border-input bg-background px-3 text-sm"
              disabled={isLoading || actionLoadingID !== null}
            >
              <option value="">Todos</option>
              <option value="requested">Solicitado</option>
              <option value="approved">Aprobado</option>
              <option value="rejected">Rechazado</option>
              <option value="refunded">Reembolsado</option>
            </select>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {feedback ? (
            <div
              className={
                feedback.type === "success"
                  ? "rounded-lg border border-emerald-300/50 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-900"
                  : "rounded-lg border border-destructive/35 bg-destructive/10 px-4 py-3 text-sm text-destructive"
              }
            >
              {feedback.message}
            </div>
          ) : null}

          {isLoading ? (
            <p className="text-sm text-muted-foreground">Cargando solicitudes...</p>
          ) : requests.length === 0 ? (
            <p className="text-sm text-muted-foreground">No hay solicitudes para este filtro.</p>
          ) : (
            <div className="space-y-4">
              {requests.map((request) => {
                const isActionRunning = actionLoadingID === request.id;
                return (
                  <article key={request.id} className="rounded-xl border border-border/60 bg-background/60 p-4 space-y-4">
                    <div className="flex flex-wrap items-center justify-between gap-2">
                      <div>
                        <p className="text-sm font-semibold text-foreground">
                          {request.company_name ?? "Empresa"} - {request.company_rfc ?? request.company_id.slice(0, 8)}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          Request ID: {request.id} • Creada: {new Date(request.created_at).toLocaleString()}
                        </p>
                      </div>
                      <Badge className={statusBadgeClass(request.status)}>{statusLabel(request.status)}</Badge>
                    </div>

                    <div className="rounded-md border border-border/60 bg-background/70 px-3 py-2 text-sm text-foreground">
                      {request.reason}
                    </div>

                    <div className="grid gap-3 md:grid-cols-2">
                      <div className="space-y-2">
                        <label htmlFor={`note-${request.id}`} className="text-xs font-medium text-muted-foreground">
                          Nota de resolución
                        </label>
                        <textarea
                          id={`note-${request.id}`}
                          value={notes[request.id] ?? request.resolution_note ?? ""}
                          onChange={(event) =>
                            setNotes((current) => ({
                              ...current,
                              [request.id]: event.target.value,
                            }))
                          }
                          className="min-h-24 w-full rounded-md border border-input bg-background px-3 py-2 text-sm shadow-sm outline-none transition focus-visible:ring-2 focus-visible:ring-ring"
                          disabled={isActionRunning}
                        />
                      </div>

                      <div className="space-y-2">
                        <label htmlFor={`stripe-refund-${request.id}`} className="text-xs font-medium text-muted-foreground">
                          Stripe Refund ID (solo para reembolsado)
                        </label>
                        <input
                          id={`stripe-refund-${request.id}`}
                          type="text"
                          value={refundIDs[request.id] ?? request.stripe_refund_id ?? ""}
                          onChange={(event) =>
                            setRefundIDs((current) => ({
                              ...current,
                              [request.id]: event.target.value,
                            }))
                          }
                          className="h-10 w-full rounded-md border border-input bg-background px-3 text-sm shadow-sm outline-none transition focus-visible:ring-2 focus-visible:ring-ring"
                          placeholder="re_xxx"
                          disabled={isActionRunning}
                        />
                        {request.reviewed_by_email ? (
                          <p className="text-xs text-muted-foreground">
                            Revisado por: {request.reviewed_by_email}
                            {request.reviewed_at ? ` (${new Date(request.reviewed_at).toLocaleString()})` : ""}
                          </p>
                        ) : null}
                      </div>
                    </div>

                    <div className="flex flex-wrap gap-2">
                      <Button
                        type="button"
                        variant="outline"
                        disabled={isActionRunning || request.status !== "requested"}
                        onClick={() => void runDecision(request, "approved")}
                      >
                        Aprobar
                      </Button>
                      <Button
                        type="button"
                        variant="outline"
                        disabled={isActionRunning || (request.status !== "requested" && request.status !== "approved")}
                        onClick={() => void runDecision(request, "rejected")}
                      >
                        Rechazar
                      </Button>
                      <Button
                        type="button"
                        disabled={isActionRunning || (request.status !== "requested" && request.status !== "approved")}
                        onClick={() => void runDecision(request, "refunded")}
                      >
                        {isActionRunning ? "Guardando..." : "Marcar Reembolsado"}
                      </Button>
                    </div>
                  </article>
                );
              })}
            </div>
          )}

          <div className="flex flex-wrap items-center justify-between gap-3 border-t border-border/60 pt-3">
            <p className="text-xs text-muted-foreground">{pageLabel}</p>
            <div className="flex gap-2">
              <Button
                type="button"
                variant="outline"
                onClick={() => void loadRequests(Math.max(offset - PAGE_SIZE, 0))}
                disabled={!hasPrev || isLoading || actionLoadingID !== null}
              >
                Anterior
              </Button>
              <Button
                type="button"
                variant="outline"
                onClick={() => void loadRequests(offset + PAGE_SIZE)}
                disabled={!hasNext || isLoading || actionLoadingID !== null}
              >
                Siguiente
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
