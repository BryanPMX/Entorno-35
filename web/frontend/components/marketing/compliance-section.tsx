import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { CheckCircle, AlertTriangle, Users, FileCheck, TrendingUp } from "lucide-react";

/**
 * Compliance Section
 *
 * Details NOM-035 compliance benefits and requirements coverage.
 */
export function ComplianceSection() {
  const compliancePoints = [
    {
      icon: Users,
      title: "Evaluación Inicial",
      description: "Identificación de factores de riesgo psicosocial en el centro de trabajo.",
      requirement: "Artículo 6"
    },
    {
      icon: AlertTriangle,
      title: "Evaluación de Cambio",
      description: "Evaluaciones cuando cambian condiciones que afectan la salud mental.",
      requirement: "Artículo 7"
    },
    {
      icon: FileCheck,
      title: "Plan de Prevención",
      description: "Desarrollo de planes para atender riesgos identificados.",
      requirement: "Artículo 8"
    },
    {
      icon: TrendingUp,
      title: "Seguimiento Continuo",
      description: "Monitoreo y actualización constante del programa de prevención.",
      requirement: "Artículo 9"
    }
  ];

  const benefits = [
    "Cumple 100% con NOM-035 STPS 2018",
    "Reportes aceptados por autoridad laboral",
    "Reducción de riesgos de multas y sanciones",
    "Mejora del clima laboral y productividad",
    "Protección legal para la organización",
    "Datos para toma de decisiones estratégicas"
  ];

  return (
    <section id="compliance" className="py-20 lg:py-28 bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4">
            Cumplimiento NOM-035
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 mb-4">
            Norma Oficial Mexicana 035
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Nuestra plataforma está diseñada específicamente para cubrir todos los
            requisitos de la NOM-035 en materia de factores de riesgo psicosocial.
          </p>
        </div>

        <div className="grid lg:grid-cols-2 gap-12 items-center mb-16">
          {/* Compliance Requirements */}
          <div>
            <h3 className="text-2xl font-bold text-gray-900 mb-8">
              Requisitos Cubiertos
            </h3>
            <div className="space-y-6">
              {compliancePoints.map((point, index) => {
                const IconComponent = point.icon;
                return (
                  <div key={index} className="flex items-start space-x-4">
                    <div className="flex-shrink-0">
                      <div className="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                        <IconComponent className="w-5 h-5 text-green-600" />
                      </div>
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-1">
                        <h4 className="font-semibold text-gray-900">
                          {point.title}
                        </h4>
                        <Badge variant="outline" className="text-xs">
                          {point.requirement}
                        </Badge>
                      </div>
                      <p className="text-gray-600">
                        {point.description}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Benefits */}
          <div>
            <h3 className="text-2xl font-bold text-gray-900 mb-8">
              Beneficios del Cumplimiento
            </h3>
            <Card className="border-green-200 bg-green-50">
              <CardHeader>
                <CardTitle className="flex items-center text-green-800">
                  <CheckCircle className="w-5 h-5 mr-2" />
                  Protección Integral
                </CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-3">
                  {benefits.map((benefit, index) => (
                    <li key={index} className="flex items-start">
                      <CheckCircle className="w-4 h-4 text-green-600 mt-0.5 mr-3 flex-shrink-0" />
                      <span className="text-green-800">{benefit}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          </div>
        </div>

        {/* Statistics */}
        <div className="bg-gradient-to-r from-blue-600 to-purple-600 rounded-2xl p-8 text-white text-center">
          <h3 className="text-2xl font-bold mb-4">
            Más de 500 Organizaciones ya Cumplen con NOM-035
          </h3>
          <p className="text-blue-100 mb-6 max-w-2xl mx-auto">
            Empresas de todos los tamaños confían en Entorno 35 para mantener
            el cumplimiento con la norma y proteger la salud mental de sus colaboradores.
          </p>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-3xl mx-auto">
            <div>
              <div className="text-3xl font-bold mb-2">98%</div>
              <div className="text-blue-100">Tasa de Adopción</div>
            </div>
            <div>
              <div className="text-3xl font-bold mb-2">15min</div>
              <div className="text-blue-100">Tiempo Promedio por Evaluación</div>
            </div>
            <div>
              <div className="text-3xl font-bold mb-2">24/7</div>
              <div className="text-blue-100">Soporte Técnico</div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}