"use client";

import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2, CreditCard, Calendar, AlertTriangle, CheckCircle } from "lucide-react";
import { paymentService } from "@/services/payment.service";
import { toast } from "sonner";
import type { Subscription } from "@/types/backend";

/**
 * Subscription Manager
 *
 * Dashboard component for managing company subscriptions.
 * Shows current subscription status and allows upgrades/cancellations.
 */
export function SubscriptionManager() {
  const [isLoading, setIsLoading] = useState(false);
  const queryClient = useQueryClient();

  // Fetch current subscription
  const {
    data: subscription,
    isLoading: isLoadingSubscription,
    error: subscriptionError,
  } = useQuery({
    queryKey: ["subscription"],
    queryFn: () => paymentService.getSubscription(),
    retry: false, // Don't retry on 404 (no subscription)
  });

  // Cancel subscription mutation
  const cancelMutation = useMutation({
    mutationFn: () => paymentService.cancelSubscription(),
    onSuccess: () => {
      toast.success("Suscripción cancelada", {
        description: "Tu suscripción se cancelará al final del período actual.",
      });
      queryClient.invalidateQueries({ queryKey: ["subscription"] });
    },
    onError: (error: any) => {
      toast.error("Error al cancelar", {
        description: error.response?.data?.error || "No se pudo cancelar la suscripción",
      });
    },
  });

  const handleCancelSubscription = async () => {
    if (!subscription) return;

    const confirmed = window.confirm(
      "¿Estás seguro de que quieres cancelar tu suscripción? Se cancelará al final del período actual."
    );

    if (confirmed) {
      cancelMutation.mutate();
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("es-MX", {
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  };

  const getStatusBadge = (status: string) => {
    switch (status) {
      case "active":
        return <Badge className="bg-green-100 text-green-800">Activa</Badge>;
      case "inactive":
        return <Badge variant="destructive">Inactiva</Badge>;
      case "canceled":
        return <Badge className="bg-yellow-100 text-yellow-800">Cancelada</Badge>;
      case "past_due":
        return <Badge className="bg-red-100 text-red-800">Vencida</Badge>;
      default:
        return <Badge variant="secondary">{status}</Badge>;
    }
  };

  if (isLoadingSubscription) {
    return (
      <Card>
        <CardContent className="flex items-center justify-center py-8">
          <Loader2 className="w-6 h-6 animate-spin mr-2" />
          <span>Cargando suscripción...</span>
        </CardContent>
      </Card>
    );
  }

  if (subscriptionError) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <CreditCard className="w-5 h-5 mr-2" />
            Suscripción
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Alert>
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>
              No tienes una suscripción activa.{" "}
              <a href="/#pricing" className="text-blue-600 hover:underline">
                Ver planes de precios
              </a>
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    );
  }

  if (!subscription) {
    return (
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center">
            <CreditCard className="w-5 h-5 mr-2" />
            Suscripción
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Alert>
            <AlertTriangle className="h-4 w-4" />
            <AlertDescription>
              No se encontró información de suscripción.{" "}
              <a href="/#pricing" className="text-blue-600 hover:underline">
                Ver planes disponibles
              </a>
            </AlertDescription>
          </Alert>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center">
          <CreditCard className="w-5 h-5 mr-2" />
          Suscripción
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        {/* Subscription Status */}
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm text-gray-600">Estado de la suscripción</p>
            <div className="flex items-center space-x-2 mt-1">
              {getStatusBadge(subscription.status)}
              {subscription.cancel_at_period_end && (
                <Badge variant="outline" className="text-orange-600">
                  Se cancelará
                </Badge>
              )}
            </div>
          </div>
        </div>

        {/* Subscription Details */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-gray-600">Plan</p>
            <p className="font-medium">
              {subscription.interval === "month" ? "Mensual" : "Anual"} - $500 MXN
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600">ID de Suscripción</p>
            <p className="font-mono text-sm text-gray-500">
              {subscription.stripe_subscription_id}
            </p>
          </div>
        </div>

        {/* Billing Period */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <p className="text-sm text-gray-600 flex items-center">
              <Calendar className="w-4 h-4 mr-1" />
              Período actual
            </p>
            <p className="font-medium">
              {formatDate(subscription.current_period_start)} - {formatDate(subscription.current_period_end)}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-600">Próxima renovación</p>
            <p className="font-medium">
              {formatDate(subscription.current_period_end)}
            </p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 pt-4 border-t">
          {subscription.status === "active" && !subscription.cancel_at_period_end && (
            <Button
              variant="destructive"
              onClick={handleCancelSubscription}
              disabled={cancelMutation.isPending}
              className="flex-1"
            >
              {cancelMutation.isPending ? (
                <>
                  <Loader2 className="w-4 h-4 animate-spin mr-2" />
                  Cancelando...
                </>
              ) : (
                "Cancelar Suscripción"
              )}
            </Button>
          )}

          {subscription.cancel_at_period_end && (
            <div className="flex items-center text-orange-600">
              <AlertTriangle className="w-4 h-4 mr-2" />
              <span className="text-sm">
                Suscripción cancelada - vence el {formatDate(subscription.current_period_end)}
              </span>
            </div>
          )}

          <Button variant="outline" asChild>
            <a href="/#pricing">Ver Planes</a>
          </Button>
        </div>

        {/* Success Message for Cancellation */}
        {cancelMutation.isSuccess && (
          <Alert className="border-green-200 bg-green-50">
            <CheckCircle className="h-4 w-4 text-green-600" />
            <AlertDescription className="text-green-800">
              Tu suscripción ha sido cancelada exitosamente. Podrás usar el servicio hasta el final del período actual.
            </AlertDescription>
          </Alert>
        )}
      </CardContent>
    </Card>
  );
}