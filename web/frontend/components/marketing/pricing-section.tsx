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
      name: "Básico",
      price: "$299",
      period: "mes",
      description: "Perfecto para pequeñas empresas",
      employees: "Hasta 50 empleados",
      features: [
        "Evaluaciones NOM-035 completas",
        "Reportes PDF profesionales",
        "Dashboard básico",
        "Soporte por email",
        "Importación CSV de empleados",
        "Enlaces seguros para evaluaciones"
      ],
      popular: false,
      cta: "Comenzar Prueba Gratuita"
    },
    {
      name: "Profesional",
      price: "$599",
      period: "mes",
      description: "Para organizaciones en crecimiento",
      employees: "Hasta 200 empleados",
      features: [
        "Todo lo del plan Básico",
        "Análisis avanzado de riesgos",
        "Gráficos y reportes detallados",
        "Mapas de calor por departamento",
        "Evaluaciones de seguimiento",
        "Soporte prioritario",
        "API para integraciones",
        "Reportes personalizados"
      ],
      popular: true,
      cta: "Comenzar Prueba Gratuita"
    },
    {
      name: "Empresarial",
      price: "Personalizado",
      period: "",
      description: "Para grandes organizaciones",
      employees: "Sin límite de empleados",
      features: [
        "Todo lo del plan Profesional",
        "Implementación dedicada",
        "Soporte 24/7 telefónico",
        "Integración con sistemas HR",
        "Reportes ejecutivos avanzados",
        "Capacitación del equipo",
        "SLA garantizado",
        "Cuenta dedicada"
      ],
      popular: false,
      cta: "Contactar Ventas"
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
            Elige el Plan Perfecto para tu Organización
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Precios transparentes sin costos ocultos. Todos los planes incluyen
            30 días de prueba gratuita.
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
                ¿Qué incluye la prueba gratuita?
              </h4>
              <p className="text-gray-600 text-sm">
                Acceso completo a todas las funcionalidades por 30 días.
                Sin límite de evaluaciones durante el período de prueba.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Puedo cambiar de plan en cualquier momento?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, puedes actualizar o cambiar tu plan en cualquier momento.
                Los cambios se reflejan inmediatamente en tu próxima factura.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Los datos están seguros?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, utilizamos encriptación SSL de nivel bancario y cumplimiento
                con NOM-035 en materia de confidencialidad de datos personales.
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-gray-900 mb-2">
                ¿Ofrecen soporte técnico?
              </h4>
              <p className="text-gray-600 text-sm">
                Sí, todos los planes incluyen soporte técnico. El plan Empresarial
                incluye soporte telefónico 24/7 con SLA garantizado.
              </p>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}