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
      microcopy: "Adaptive branching",
    },
    {
      icon: FileText,
      title: "Reportes PDF Profesionales",
      description: "Reportes de cumplimiento NOM-035 listos para presentación ante autoridades.",
      badge: "Profesional",
      accent: "from-blue-500/25 to-blue-500/5",
      microcopy: "Brand-ready outputs",
    },
    {
      icon: Users,
      title: "Gestión de Personal",
      description: "Importación masiva de empleados via CSV con validación automática.",
      badge: "Eficiente",
      accent: "from-indigo-500/25 to-indigo-500/5",
      microcopy: "CSV + API ready",
    },
    {
      icon: BarChart3,
      title: "Análisis de Riesgos",
      description: "Visualización de datos con gráficos y mapas de calor por departamento.",
      badge: "Analítico",
      accent: "from-purple-500/25 to-purple-500/5",
      microcopy: "Heatmaps & cohorts",
    },
    {
      icon: Shield,
      title: "Cumplimiento NOM-035",
      description: "100% compatible con la Norma Oficial Mexicana 035 de seguridad y salud laboral.",
      badge: "Certificado",
      accent: "from-cyan-500/25 to-cyan-500/5",
      microcopy: "Auditoría lista",
    },
    {
      icon: Clock,
      title: "Evaluaciones Express",
      description: "Los empleados completan evaluaciones en menos de 15 minutos.",
      badge: "Rápido",
      accent: "from-amber-500/25 to-amber-500/5",
      microcopy: "15 min promedio",
    },
    {
      icon: Mail,
      title: "Envío Automático",
      description: "Links seguros enviados automáticamente por email a todo el personal.",
      badge: "Automático",
      accent: "from-pink-500/25 to-pink-500/5",
      microcopy: "Envía y monitorea",
    },
    {
      icon: Smartphone,
      title: "Acceso Móvil",
      description: "Interfaz totalmente responsive optimizada para dispositivos móviles.",
      badge: "Móvil",
      accent: "from-blue-500/25 to-blue-500/5",
      microcopy: "Optimizada PWA",
    },
    {
      icon: Globe,
      title: "Aislamiento Seguro",
      description: "Cada empresa tiene su propio espacio seguro con datos completamente aislados.",
      badge: "Seguro",
      accent: "from-sky-500/25 to-sky-500/5",
      microcopy: "Multi-tenant seguro",
    },
    {
      icon: Lock,
      title: "Encriptación SSL",
      description: "Todos los datos transmitidos y almacenados con encriptación de nivel bancario.",
      badge: "Encriptado",
      accent: "from-slate-500/25 to-slate-500/5",
      microcopy: "TLS + at-rest",
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

        {/* Feature clusters to avoid repetition */}
        <div className="grid gap-8 lg:grid-cols-1 xl:grid-cols-[1.1fr_0.9fr] mb-10">
          <Card className="glass-panel-strong shimmer-border h-full">
            <CardContent className="p-6 space-y-4 h-full flex flex-col">
              <div className="flex items-center gap-3">
                <span className="text-xs font-semibold uppercase tracking-wide text-blue-700 dark:text-blue-200 bg-blue-100/70 dark:bg-blue-500/20 px-3 py-1 rounded-full">
                  Operación
                </span>
                <p className="text-sm text-gray-600 dark:text-gray-300">Flujos que reducen tiempo de despliegue</p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {features.slice(0, 5).map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div
                      key={feature.title}
                      className="group rounded-xl border border-white/30 dark:border-white/10 bg-gradient-to-br from-white/80 to-white/50 dark:from-white/5 dark:to-white/0 p-4 hover:translate-y-[-3px] transition-all shadow-sm h-full flex flex-col gap-2"
                    >
                      <div className="flex items-center gap-3 mb-2">
                        <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center`}>
                          <IconComponent className="w-5 h-5 text-gray-900 dark:text-white" />
                        </div>
                        <div className="text-xs font-medium text-gray-500 dark:text-gray-300">{feature.microcopy}</div>
                      </div>
                      <div className="flex items-center justify-between mb-1">
                        <h3 className="text-base font-semibold text-gray-900 dark:text-white">
                          {feature.title}
                        </h3>
                        <Badge variant="secondary" className="text-[11px] px-2">
                          {feature.badge}
                        </Badge>
                      </div>
                      <p className="text-sm text-gray-600 dark:text-gray-300 flex-1">{feature.description}</p>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          <Card className="glass-panel shadow-xl h-full">
            <CardContent className="p-6 space-y-4 h-full flex flex-col">
              <div className="flex items-center gap-3">
                <span className="text-xs font-semibold uppercase tracking-wide text-purple-700 dark:text-purple-200 bg-purple-100/70 dark:bg-purple-500/20 px-3 py-1 rounded-full">
                  Seguridad & Acceso
                </span>
                <p className="text-sm text-gray-600 dark:text-gray-300">Confianza en cada interacción</p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                {features.slice(5, 10).map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div
                      key={feature.title}
                      className="rounded-xl border border-white/40 dark:border-white/10 bg-gradient-to-br from-white/70 to-white/40 dark:from-white/10 dark:to-white/0 p-4 hover:shadow-lg transition-all h-full flex flex-col gap-2"
                    >
                      <div className="flex items-start gap-3">
                        <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center`}>
                          <IconComponent className="w-5 h-5 text-gray-900 dark:text-white" />
                        </div>
                        <div className="space-y-1">
                          <div className="flex items-center gap-2">
                            <h3 className="text-base font-semibold text-gray-900 dark:text-white">
                              {feature.title}
                            </h3>
                            <Badge variant="secondary" className="text-[11px] px-2">
                              {feature.badge}
                            </Badge>
                          </div>
                          <p className="text-sm text-gray-600 dark:text-gray-300 leading-relaxed">
                            {feature.description}
                          </p>
                          <p className="text-xs text-gray-500 dark:text-gray-400">{feature.microcopy}</p>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Bottom CTA */}
        <div className="mt-16 grid gap-6 lg:grid-cols-1 xl:grid-cols-[1.1fr_0.9fr] items-stretch">
          <Card className="glass-panel-strong shimmer-border h-full">
            <CardContent className="p-6 space-y-4 h-full flex flex-col">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm uppercase tracking-wide text-gray-500 dark:text-gray-300 mb-1">
                    Implementación guiada
                  </p>
                  <h3 className="text-xl font-semibold text-gray-900 dark:text-white">
                    En vivo con un especialista NOM-035
                  </h3>
                  <p className="text-sm text-gray-600 dark:text-gray-300">
                    Incluye checklist, dominio seguro, branding y plantillas PDF listas.
                  </p>
                </div>
                <Badge className="badge-glow text-white px-3 py-1 shadow">72h Go-Live</Badge>
              </div>
              <div className="grid sm:grid-cols-3 gap-3 text-sm text-gray-700 dark:text-gray-200">
                {[
                  "Diagnóstico inicial",
                  "Configuración y pruebas",
                  "Capacitación equipos",
                ].map((item) => (
                  <div key={item} className="flex items-center gap-2 glass-panel rounded-lg px-3 py-2 border border-white/40 dark:border-white/10">
                    <CheckCircle className="w-4 h-4 text-green-500" />
                    <span>{item}</span>
                  </div>
                ))}
              </div>
              <Link
                href="#pricing"
                className="inline-flex items-center px-5 py-3 rounded-md text-white bg-gradient-to-r from-blue-600 to-purple-600 shadow-lg hover:shadow-xl transition w-fit"
              >
                Ver precios
                <span className="ml-2">→</span>
              </Link>
            </CardContent>
          </Card>

          <Card className="glass-panel border-dashed border-2 border-white/40 dark:border-white/10 h-full">
            <CardContent className="p-6 space-y-4 text-gray-700 dark:text-gray-200 h-full flex flex-col">
              <div className="flex items-center gap-3">
                <div className="h-12 w-12 rounded-full bg-gradient-to-br from-emerald-500 to-blue-500 opacity-90 shadow-inner" />
                <div>
                  <p className="text-sm text-gray-500 dark:text-gray-300">SLA de soporte</p>
                  <p className="text-lg font-semibold">Respuesta promedio 2h</p>
                </div>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-sm flex-1">
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  Chat y email en horario laboral
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  Escalamiento crítico &lt; 1h
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  Revisiones trimestrales de riesgo
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="w-4 h-4 text-green-500" />
                  Playbooks y runbooks listos
                </div>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}
