import Link from "next/link";
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
      badge: "Inteligente",
      accent: "from-emerald-500/25 to-emerald-500/5",
    },
    {
      icon: FileText,
      title: "Reportes PDF Profesionales",
      description: "Reportes de cumplimiento NOM-035 listos para presentación ante autoridades.",
      badge: "Profesional",
      accent: "from-blue-500/25 to-blue-500/5",
    },
    {
      icon: Users,
      title: "Gestión de Personal",
      description: "Importación masiva de empleados via CSV con validación automática.",
      badge: "Eficiente",
      accent: "from-indigo-500/25 to-indigo-500/5",
    },
    {
      icon: BarChart3,
      title: "Análisis de Riesgos",
      description: "Visualización de datos con gráficos y mapas de calor por departamento.",
      badge: "Analítico",
      accent: "from-purple-500/25 to-purple-500/5",
    },
    {
      icon: Shield,
      title: "Cumplimiento NOM-035",
      description: "100% compatible con la Norma Oficial Mexicana 035 de seguridad y salud laboral.",
      badge: "Certificado",
      accent: "from-cyan-500/25 to-cyan-500/5",
    },
    {
      icon: Clock,
      title: "Evaluaciones Express",
      description: "Los empleados completan evaluaciones en menos de 15 minutos.",
      badge: "Rápido",
      accent: "from-amber-500/25 to-amber-500/5",
    },
    {
      icon: Mail,
      title: "Envío Automático",
      description: "Links seguros enviados automáticamente por email a todo el personal.",
      badge: "Automático",
      accent: "from-pink-500/25 to-pink-500/5",
    },
    {
      icon: Smartphone,
      title: "Acceso Móvil",
      description: "Interfaz totalmente responsive optimizada para dispositivos móviles.",
      badge: "Móvil",
      accent: "from-blue-500/25 to-blue-500/5",
    },
    {
      icon: Globe,
      title: "Aislamiento Seguro",
      description: "Cada empresa tiene su propio espacio seguro con datos completamente aislados.",
      badge: "Seguro",
      accent: "from-sky-500/25 to-sky-500/5",
    },
    {
      icon: Lock,
      title: "Encriptación SSL",
      description: "Todos los datos transmitidos y almacenados con encriptación de nivel bancario.",
      badge: "Encriptado",
      accent: "from-slate-500/25 to-slate-500/5",
    }
  ];

  return (
    <section id="features" className="relative overflow-hidden py-20 lg:py-28 bg-gradient-to-b from-gray-50 via-white to-gray-50 dark:from-background dark:via-background dark:to-background">
      <div className="absolute inset-0 pointer-events-none bg-aurora opacity-60 dark:opacity-35" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 relative">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4 glass-panel px-4 py-2">
            Características Principales
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-gray-900 dark:text-white mb-4">
            Todo lo que necesitas para cumplir con{" "}
            <span className="gradient-text">NOM-035</span>
          </h2>
          <p className="text-lg text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Plataforma creada para equipos de Seguridad e Higiene que buscan precisión, trazabilidad y
            velocidad al gestionar el riesgo psicosocial.
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {features.map((feature, index) => {
            const IconComponent = feature.icon;
            return (
              <Card
                key={feature.title}
                className="glass-panel hover:translate-y-[-4px] transition-all duration-300 border border-white/50 dark:border-white/10 shadow-lg"
              >
                <CardContent className="p-6 space-y-4">
                  <div className="flex items-start gap-4">
                    <div className={`w-12 h-12 rounded-xl bg-gradient-to-br ${feature.accent} flex items-center justify-center shadow-inner`}>
                      <IconComponent className="w-6 h-6 text-gray-900 dark:text-white" />
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center justify-between gap-2 mb-1">
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                          {feature.title}
                        </h3>
                        <Badge variant="secondary" className="text-xs px-2 py-1">
                          {feature.badge}
                        </Badge>
                      </div>
                      <p className="text-gray-600 dark:text-gray-300 leading-relaxed">
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
        <div className="mt-16 grid gap-6 lg:grid-cols-[1.2fr_0.8fr] items-center">
          <Card className="glass-panel-strong shimmer-border">
            <CardContent className="p-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div>
                <p className="text-sm uppercase tracking-wide text-gray-500 dark:text-gray-300 mb-1">
                  Implementación guiada
                </p>
                <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                  Acompañamiento de onboarding y revisión de cumplimiento
                </h3>
                <p className="text-sm text-gray-600 dark:text-gray-300">
                  Incluye checklist NOM-035, configuración de dominios y plantillas PDF.
                </p>
              </div>
              <Link
                href="#pricing"
                className="inline-flex items-center px-5 py-3 rounded-md text-white bg-gradient-to-r from-blue-600 to-purple-600 shadow-lg hover:shadow-xl transition"
              >
                Ver precios
                <span className="ml-2">→</span>
              </Link>
            </CardContent>
          </Card>

          <Card className="glass-panel border-dashed border-2 border-white/40 dark:border-white/10">
            <CardContent className="p-6 space-y-3 text-gray-700 dark:text-gray-200">
              <div className="flex items-center gap-2">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-emerald-500 to-blue-500 opacity-90" />
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-300">SLA de soporte</p>
                  <p className="text-lg font-semibold">Respuesta promedio 2h</p>
                </div>
              </div>
              <p className="text-sm">
                Soporte humano en español, playbooks de mitigación y sesiones de revisión trimestral de riesgo.
              </p>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}
