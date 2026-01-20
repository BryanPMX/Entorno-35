import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  CheckCircle,
  FileText,
  Users,
  Shield,
  BarChart3,
  Clock,
  Mail,
  Smartphone,
  Globe,
  Lock
} from "lucide-react";

/**
 * Features Section
 *
 * Showcases key features and benefits of the platform.
 */
export function FeaturesSection() {
  const features = [
    {
      icon: CheckCircle,
      title: "Evaluaciones Automatizadas",
      description: "Cuestionarios inteligentes que se adaptan según las respuestas del empleado.",
      badge: "Inteligente"
    },
    {
      icon: FileText,
      title: "Reportes PDF Profesionales",
      description: "Reportes de cumplimiento NOM-035 listos para presentación ante autoridades.",
      badge: "Profesional"
    },
    {
      icon: Users,
      title: "Gestión de Personal",
      description: "Importación masiva de empleados via CSV con validación automática.",
      badge: "Eficiente"
    },
    {
      icon: BarChart3,
      title: "Análisis de Riesgos",
      description: "Visualización de datos con gráficos y mapas de calor por departamento.",
      badge: "Analítico"
    },
    {
      icon: Shield,
      title: "Cumplimiento NOM-035",
      description: "100% compatible con la Norma Oficial Mexicana 035 de seguridad y salud laboral.",
      badge: "Certificado"
    },
    {
      icon: Clock,
      title: "Evaluaciones Express",
      description: "Los empleados completan evaluaciones en menos de 15 minutos.",
      badge: "Rápido"
    },
    {
      icon: Mail,
      title: "Envío Automático",
      description: "Links seguros enviados automáticamente por email a todo el personal.",
      badge: "Automático"
    },
    {
      icon: Smartphone,
      title: "Acceso Móvil",
      description: "Interfaz totalmente responsive optimizada para dispositivos móviles.",
      badge: "Móvil"
    },
    {
      icon: Globe,
      title: "Aislamiento Seguro",
      description: "Cada empresa tiene su propio espacio seguro con datos completamente aislados.",
      badge: "Seguro"
    },
    {
      icon: Lock,
      title: "Encriptación SSL",
      description: "Todos los datos transmitidos y almacenados con encriptación de nivel bancario.",
      badge: "Encriptado"
    }
  ];

  return (
    <section id="features" className="py-20 lg:py-28 bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4">
            Características Principales
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 mb-4">
            Todo lo que Necesitas para Cumplir con NOM-035
          </h2>
          <p className="text-xl text-gray-600 max-w-3xl mx-auto">
            Una plataforma completa diseñada específicamente para organizaciones mexicanas
            que requieren cumplimiento con la Norma Oficial Mexicana 035.
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => {
            const IconComponent = feature.icon;
            return (
              <Card key={index} className="border-0 shadow-lg hover:shadow-xl transition-shadow duration-300">
                <CardContent className="p-6">
                  <div className="flex items-start space-x-4">
                    <div className="flex-shrink-0">
                      <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
                        <IconComponent className="w-6 h-6 text-blue-600" />
                      </div>
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-2">
                        <h3 className="text-lg font-semibold text-gray-900">
                          {feature.title}
                        </h3>
                        <Badge variant="secondary" className="text-xs">
                          {feature.badge}
                        </Badge>
                      </div>
                      <p className="text-gray-600 leading-relaxed">
                        {feature.description}
                      </p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>

        {/* Bottom CTA */}
        <div className="text-center mt-16">
          <p className="text-gray-600 mb-6">
            ¿Listo para simplificar el cumplimiento NOM-035 en tu organización?
          </p>
          <Card className="max-w-md mx-auto border-blue-200 bg-blue-50">
            <CardContent className="p-6 text-center">
              <h3 className="font-semibold text-gray-900 mb-2">
                Prueba Gratuita por 30 Días
              </h3>
              <p className="text-sm text-gray-600 mb-4">
                Sin compromiso. Configura tu cuenta y comienza a evaluar.
              </p>
              <a
                href="/login"
                className="inline-flex items-center px-6 py-3 border border-transparent text-base font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 transition-colors"
              >
                Comenzar Ahora
              </a>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}