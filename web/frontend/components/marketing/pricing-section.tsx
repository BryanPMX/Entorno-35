import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Check, Sparkles } from "lucide-react";
import Link from "next/link";

/** Shared feature list for both subscription tiers */
const PLAN_FEATURES = [
  "Evaluaciones NOM-035 completas",
  "Reportes PDF profesionales",
  "Dashboard de análisis de riesgos",
  "Soporte por email",
  "Importación CSV de empleados",
  "Enlaces seguros para evaluaciones",
  "Gráficos y mapas de calor por departamento",
  "Envío automático de evaluaciones por correo",
];

/**
 * Pricing Section
 *
 * Two subscription options: monthly ($1,000 MXN) and yearly ($6,000 MXN).
 * Stripe integration to be added later.
 */
export function PricingSection() {
  const plans = [
    {
      id: "monthly",
      name: "Mensual",
      price: "$1,000",
      currency: "MXN",
      period: "mes",
      description: "Facturación mensual, cancela cuando quieras.",
      features: PLAN_FEATURES,
      popular: false,
      cta: "Comenzar ahora",
      ctaHref: "/login",
    },
    {
      id: "yearly",
      name: "Anual",
      price: "$6,000",
      currency: "MXN",
      period: "año",
      description: "Paga una vez al año y ahorra el equivalente a 6 meses.",
      savings: "Ahorra $6,000 vs mensual (50%)",
      features: PLAN_FEATURES,
      popular: true,
      cta: "Comenzar ahora",
      ctaHref: "/login",
    },
  ];

  return (
    <section id="pricing" className="relative overflow-hidden py-20 lg:py-28 bg-gradient-to-b from-gray-50 via-white to-gray-50 dark:from-background dark:via-background dark:to-background">
      <div className="absolute inset-0 -z-10 soft-grid opacity-45" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-16 glass-panel rounded-2xl p-8 md:p-10 border border-white/50 dark:border-white/10 shadow-lg">
          <Badge variant="outline" className="mb-4 glass-panel px-4 py-2">
            Planes y Precios
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Elige tu suscripción
          </h2>
          <p className="text-lg text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Precios en pesos mexicanos (MXN). Todos los planes incluyen las mismas
            funcionalidades; elige la periodicidad que mejor se adapte a tu organización.
          </p>
        </div>

        {/* Pricing Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
          {plans.map((plan) => (
            <Card
              key={plan.id}
              className={`relative flex flex-col ${
                plan.popular
                  ? "glass-panel-strong shimmer-border shadow-xl"
                  : "glass-panel shadow-lg"
              }`}
            >
              {plan.popular && (
                <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                  <Badge className="badge-glow text-white px-4 py-1 shadow-md">
                    <Sparkles className="w-3 h-3 mr-1 inline" />
                    50% OFF anual
                  </Badge>
                </div>
              )}

              <CardHeader className="text-center pb-6 pt-8">
                <CardTitle className="text-2xl font-bold text-gray-900">
                  {plan.name}
                </CardTitle>
                <div className="mt-4 flex items-baseline justify-center gap-1">
                  <span className="text-4xl font-bold text-gray-900">
                    {plan.price}
                  </span>
                  <span className="text-gray-500 text-lg">{plan.currency}</span>
                  {plan.period && (
                    <span className="text-gray-600 ml-1">/{plan.period}</span>
                  )}
                </div>
                {plan.savings && (
                  <p className="text-sm font-medium text-green-600 mt-2">
                    {plan.savings}
                  </p>
                )}
                <p className="text-gray-600 mt-2 text-sm">{plan.description}</p>
              </CardHeader>

              <CardContent className="flex-1 flex flex-col space-y-6">
                <p className="text-sm text-gray-600 dark:text-gray-300 text-center">
                  Incluye todas las capacidades NOM-035: evaluaciones completas, reportes PDF, dashboard, envíos automáticos y soporte.
                </p>

                <div className="mt-auto">
                  <Button
                    className={`w-full ${
                      plan.popular ? "bg-gradient-to-r from-blue-600 to-purple-600 hover:brightness-110" : ""
                    }`}
                    variant={plan.popular ? "default" : "outline"}
                    size="lg"
                    asChild
                  >
                    <Link href={plan.ctaHref}>{plan.cta}</Link>
                  </Button>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* Unified feature list (same for ambos planes) */}
        <div className="mt-10">
          <Card className="glass-panel border border-white/50 dark:border-white/10 shadow-md">
            <CardContent className="p-6">
              <p className="text-sm font-semibold text-gray-900 dark:text-white mb-4">
                Todo lo que incluyen ambos planes
              </p>
              <div className="grid sm:grid-cols-2 md:grid-cols-3 gap-3">
                {PLAN_FEATURES.map((feature) => (
                  <div key={feature} className="flex items-start gap-2 text-sm text-gray-700 dark:text-gray-200">
                    <Check className="w-4 h-4 text-green-500 mt-0.5 flex-shrink-0" />
                    <span>{feature}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Billing note - Stripe later */}
        <p className="text-center text-sm text-gray-500 dark:text-gray-300 mt-8">
          Pago seguro. La integración con pasarela de pago se activará próximamente.
        </p>

        {/* FAQ */}
        <div className="mt-16 text-center">
          <h3 className="text-2xl font-bold text-gray-900 mb-8">
            Preguntas Frecuentes
          </h3>
          <div className="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto text-left">
            <div className="glass-panel rounded-xl p-5 border border-white/50 dark:border-white/10">
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">
                ¿Qué diferencia hay entre Mensual y Anual?
              </h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">
                Las funcionalidades son idénticas. El plan anual cuesta $6,000 MXN (equivalente a 6 meses);
                pagar mes a mes durante 12 meses costaría $12,000 MXN. Con el plan anual ahorras $6,000 MXN
                y obtienes un 50% de descuento efectivo.
              </p>
            </div>
            <div className="glass-panel rounded-xl p-5 border border-white/50 dark:border-white/10">
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">
                ¿Puedo cambiar de plan después?
              </h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">
                Sí. Puedes pasar de mensual a anual en cualquier momento para
                aprovechar el ahorro, o mantener la facturación mensual si lo
                prefieres.
              </p>
            </div>
            <div className="glass-panel rounded-xl p-5 border border-white/50 dark:border-white/10">
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">
                ¿Los datos están seguros?
              </h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">
                Sí. Utilizamos encriptación SSL de nivel bancario y cumplimiento
                con NOM-035 en materia de confidencialidad de datos personales.
              </p>
            </div>
            <div className="glass-panel rounded-xl p-5 border border-white/50 dark:border-white/10">
              <h4 className="font-semibold text-gray-900 dark:text-white mb-2">
                ¿Incluyen soporte técnico?
              </h4>
              <p className="text-gray-600 dark:text-gray-300 text-sm">
                Sí. Todos los planes incluyen soporte técnico por correo para
                configuración, importación de personal y uso de la plataforma.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
