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
      accent: "from-sky-500/25 to-sky-500/5",
      microcopy: "Brand-ready outputs",
    },
    {
      icon: Users,
      title: "Gestión de Personal",
      description: "Importación masiva de empleados via CSV con validación automática.",
      badge: "Eficiente",
      accent: "from-cyan-500/25 to-cyan-500/5",
      microcopy: "CSV + API ready",
    },
    {
      icon: BarChart3,
      title: "Análisis de Riesgos",
      description: "Visualización de datos con gráficos y mapas de calor por departamento.",
      badge: "Analítico",
      accent: "from-amber-500/25 to-amber-500/5",
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
      accent: "from-orange-500/25 to-orange-500/5",
      microcopy: "Envía y monitorea",
    },
    {
      icon: Smartphone,
      title: "Acceso Móvil",
      description: "Interfaz totalmente responsive optimizada para dispositivos móviles.",
      badge: "Móvil",
      accent: "from-sky-500/25 to-sky-500/5",
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
    <section id="features" className="relative overflow-hidden py-20 lg:py-28">
      <div className="absolute inset-0 pointer-events-none bg-aurora opacity-60 dark:opacity-35" />
      <div className="absolute -left-24 top-10 h-64 w-64 rounded-full bg-[color:var(--nom-bajo-soft)] blur-3xl" />
      <div className="absolute -right-32 -bottom-10 h-80 w-80 rounded-full bg-[color:var(--nom-nulo-soft)] blur-3xl" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 relative">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4 glass-panel px-4 py-2">
            Características Principales
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-4">
            Todo lo que necesitas para cumplir con{" "}
            <span className="gradient-text">NOM-035</span>
          </h2>
          <p className="text-lg text-muted-foreground max-w-3xl mx-auto">
            Plataforma creada para equipos de Seguridad e Higiene que buscan precisión, trazabilidad y
            velocidad al gestionar el riesgo psicosocial.
          </p>
        </div>

        {/* Feature clusters to avoid repetition */}
        <div className="grid gap-8 lg:grid-cols-1 xl:grid-cols-2 mb-10 isolate items-stretch">
          <Card className="portal-surface-strong shimmer-border h-full border-0">
            <CardContent className="p-6 space-y-4 h-full flex flex-col">
              <div className="flex items-center gap-3">
                <span className="text-xs font-semibold uppercase tracking-wide text-[var(--nom-bajo)] bg-[color:var(--nom-bajo-soft)] px-3 py-1 rounded-full">
                  Operación
                </span>
                <p className="text-sm text-muted-foreground">Flujos que reducen tiempo de despliegue</p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 isolate content-start flex-1">
                {features.slice(0, 5).map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div
                      key={feature.title}
                      className="group relative z-0 rounded-xl border border-white/30 dark:border-white/10 bg-gradient-to-br from-white/80 to-white/50 dark:from-white/5 dark:to-white/0 p-4 hover:-translate-y-1 transition-all shadow-sm h-full flex flex-col gap-2 hover:z-20 overflow-hidden transform-gpu"
                    >
                      <div className="flex items-center gap-3 mb-2">
                        <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center`}>
                          <IconComponent className="w-5 h-5 text-foreground" />
                        </div>
                        <div className="text-xs font-medium text-muted-foreground">{feature.microcopy}</div>
                      </div>
                      <div className="flex items-center justify-between mb-1">
                        <h3 className="text-base font-semibold text-foreground">
                          {feature.title}
                        </h3>
                        <Badge variant="secondary" className="text-[11px] px-2">
                          {feature.badge}
                        </Badge>
                      </div>
                      <p className="text-sm text-muted-foreground flex-1">{feature.description}</p>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          <Card className="portal-surface shadow-xl h-full border-0">
            <CardContent className="p-6 space-y-4 h-full flex flex-col">
              <div className="flex items-center gap-3">
                <span className="text-xs font-semibold uppercase tracking-wide text-[var(--nom-medio)] bg-[color:var(--nom-medio-soft)] px-3 py-1 rounded-full">
                  Seguridad & Acceso
                </span>
                <p className="text-sm text-muted-foreground">Confianza en cada interacción</p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 isolate content-start flex-1">
                {features.slice(5, 10).map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div
                      key={feature.title}
                      className="relative z-0 rounded-xl border border-white/40 dark:border-white/10 bg-gradient-to-br from-white/70 to-white/40 dark:from-white/10 dark:to-white/0 p-4 hover:shadow-lg hover:-translate-y-1 transition-all h-full flex flex-col gap-2 hover:z-20 transform-gpu"
                    >
                      <div className="flex items-start gap-3">
                        <div className={`w-10 h-10 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center flex-shrink-0`}>
                          <IconComponent className="w-5 h-5 text-foreground" />
                        </div>
                        <div className="space-y-1 min-w-0 w-full">
                          <div className="flex items-start gap-2 flex-wrap">
                            <h3 className="text-base font-semibold text-foreground min-w-0 break-words">
                              {feature.title}
                            </h3>
                            <Badge variant="secondary" className="text-[11px] px-2 shrink-0">
                              {feature.badge}
                            </Badge>
                          </div>
                          <p className="text-sm text-muted-foreground leading-relaxed">
                            {feature.description}
                          </p>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Soft transition to adjacent sections */}
        <div className="h-16 relative mt-12">
          <div className="absolute inset-0 bg-gradient-to-b from-transparent via-background/70 to-background blur-[2px]" />
        </div>
      </div>
    </section>
  );
}
