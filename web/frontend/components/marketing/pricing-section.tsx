import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Check, Star } from "lucide-react";
import Link from "next/link";

/**
 * Pricing Section
 *
 * Professional pricing tiers for different company sizes.
 */
export function PricingSection() {
  const plans = [
    {
      name: "Mensual",
      price: "$500",
      period: "mes",
      description: "Pago mensual flexible",
      employees: "Sin límite de empleados",
      features: [
        "Evaluaciones NOM-035 completas",
        "Reportes PDF profesionales",
        "Dashboard completo con analytics",
        "Soporte prioritario por email",
        "Importación masiva de empleados",
        "Enlaces seguros para evaluaciones",
        "Análisis avanzado de riesgos",
        "Mapas de calor por departamento",
        "API para integraciones",
        "Encriptación de datos",
        "Cumplimiento NOM-035 garantizado"
      ],
      popular: false,
      cta: "Suscribirse Mensual"
    },
    {
      name: "Anual",
      price: "$5,000",
      period: "año",
      description: "Ahorra 17% vs mensual",
      employees: "Sin límite de empleados",
      features: [
        "Todo lo del plan Mensual",
        "Descuento del 17% anual",
        "Implementación prioritaria",
        "Soporte telefónico incluido",
        "Reportes ejecutivos avanzados",
        "Capacitación del equipo incluida",
        "SLA garantizado",
        "Actualizaciones prioritarias",
        "Backup y recuperación de datos",
        "Integración con sistemas HR"
      ],
      popular: true,
      cta: "Suscribirse Anual"
    }
  ];

  return (
    <section id="pricing" className="py-20 lg:py-28 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4">
            Planes y Precios
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 mb-4">
            Planes Transparentes para tu Empresa
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Precios simples y transparentes en pesos mexicanos. Sin costos ocultos,
            sin contratos forzosos.
          </p>
        </div>

        {/* Pricing Cards */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 max-w-5xl mx-auto">
          {plans.map((plan, index) => (
            <Card
              key={index}
              className={`relative ${
                plan.popular
                  ? 'border-blue-500 shadow-xl scale-105'
                  : 'border-gray-200 shadow-lg'
              }`}
            >
              {plan.popular && (
                <div className="absolute -top-4 left-1/2 transform -translate-x-1/2">
                  <Badge className="bg-blue-600 text-white px-4 py-1">
                    <Star className="w-3 h-3 mr-1" />
                    Más Popular
                  </Badge>
                </div>
              )}

              <CardHeader className="text-center pb-8">
                <CardTitle className="text-2xl font-bold text-gray-900">
                  {plan.name}
                </CardTitle>
                <div className="mt-4">
                  <span className="text-4xl font-bold text-gray-900">
                    {plan.price}
                  </span>
                  {plan.period && (
                    <span className="text-gray-600 ml-1">/{plan.period}</span>
                  )}
                </div>
                <p className="text-gray-600 mt-2">{plan.description}</p>
                <p className="text-sm font-medium text-blue-600 mt-2">
                  {plan.employees}
                </p>
              </CardHeader>

              <CardContent>
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature, featureIndex) => (
                    <li key={featureIndex} className="flex items-start">
                      <Check className="w-4 h-4 text-green-500 mt-0.5 mr-3 flex-shrink-0" />
                      <span className="text-gray-700 text-sm">{feature}</span>
                    </li>
                  ))}
                </ul>

                <Button
                  className={`w-full ${
                    plan.popular
                      ? 'bg-blue-600 hover:bg-blue-700'
                      : ''
                  }`}
                  variant={plan.popular ? 'default' : 'outline'}
                  asChild
                >
                  <Link href={plan.cta === "Contactar Ventas" ? "#contact" : "/login"}>
                    {plan.cta}
                  </Link>
                </Button>
              </CardContent>
            </Card>
          ))}
        </div>

        {/* FAQ Section */}
        <div className="mt-16 text-center">
          <h3 className="text-2xl font-bold text-gray-900 mb-8">
            Preguntas Frecuentes
          </h3>
          <div className="grid md:grid-cols-2 gap-8 max-w-4xl mx-auto text-left">
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Puedo cambiar entre planes mensual y anual?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, puedes cambiar entre planes en cualquier momento desde tu
                dashboard. Los cambios se prorratean automáticamente.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Qué métodos de pago aceptan?
              </h4>
              <p className="text-gray-600 text-sm">
                Aceptamos tarjetas de crédito y débito (Visa, Mastercard, American Express)
                y transferencias bancarias. Procesamos pagos de forma segura con Stripe.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Los datos están seguros?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, utilizamos encriptación SSL de nivel bancario, cumplimiento PCI DSS,
                y protección de datos conforme a NOM-035 y leyes mexicanas de privacidad.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Ofrecen soporte técnico?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, todos los planes incluyen soporte técnico por email.
                El plan Anual incluye soporte telefónico prioritario.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}