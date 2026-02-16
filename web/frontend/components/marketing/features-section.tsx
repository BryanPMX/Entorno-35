import { Card, CardContent } from "@/components/ui/card";
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
  const operacionFeatures = [
    {
      icon: CheckCircle,
      title: "Evaluaciones Automatizadas",
      description: "Cuestionarios inteligentes que se adaptan según las respuestas del empleado.",
      accent: "from-emerald-500/25 to-emerald-500/5",
    },
    {
      icon: FileText,
      title: "Reportes PDF Profesionales",
      description: "Reportes de cumplimiento NOM-035 listos para presentación ante autoridades.",
      accent: "from-blue-500/25 to-blue-500/5",
    },
    {
      icon: Users,
      title: "Gestión de Personal",
      description: "Importación masiva de empleados via CSV con validación automática.",
      accent: "from-indigo-500/25 to-indigo-500/5",
    },
    {
      icon: BarChart3,
      title: "Análisis de Riesgos",
      description: "Visualización de datos con gráficos y mapas de calor por departamento.",
      accent: "from-purple-500/25 to-purple-500/5",
    },
    {
      icon: Shield,
      title: "Cumplimiento NOM-035",
      description: "100% compatible con la Norma Oficial Mexicana 035 de seguridad y salud laboral.",
      accent: "from-cyan-500/25 to-cyan-500/5",
    },
  ];

  const seguridadFeatures = [
    {
      icon: Clock,
      title: "Evaluaciones Express",
      description: "Los empleados completan evaluaciones en menos de 15 minutos.",
      accent: "from-amber-500/25 to-amber-500/5",
    },
    {
      icon: Mail,
      title: "Envío Automático",
      description: "Links seguros enviados automáticamente por email a todo el personal.",
      accent: "from-pink-500/25 to-pink-500/5",
    },
    {
      icon: Smartphone,
      title: "Acceso Móvil",
      description: "Interfaz totalmente responsive optimizada para dispositivos móviles.",
      accent: "from-blue-500/25 to-blue-500/5",
    },
    {
      icon: Globe,
      title: "Aislamiento Seguro",
      description: "Cada empresa tiene su propio espacio seguro con datos completamente aislados.",
      accent: "from-sky-500/25 to-sky-500/5",
    },
    {
      icon: Lock,
      title: "Encriptación SSL",
      description: "Todos los datos transmitidos y almacenados con encriptación de nivel bancario.",
      accent: "from-slate-500/25 to-slate-500/5",
    },
  ];

  const sectionCardClass =
    "rounded-xl border border-white/30 dark:border-white/10 bg-gradient-to-br from-white/80 to-white/50 dark:from-white/5 dark:to-white/0 p-6 hover:-translate-y-1 transition-all shadow-sm flex flex-col items-center text-center gap-3";

  return (
    <section id="features" className="glass-section glass-section-alt-b relative overflow-hidden py-20 lg:py-28">
      <div className="absolute -left-24 top-10 h-64 w-64 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" />
      <div className="absolute -right-32 -bottom-10 h-80 w-80 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 relative">
        {/* Section Header */}
        <div className="text-center mb-16">
          <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-4">
            Características Principales
          </h2>
          <p className="text-lg text-muted-foreground max-w-3xl mx-auto">
            Todo lo que necesitas para cumplir con{" "}
            <span className="gradient-text">NOM-035</span>. Plataforma creada para equipos de
            Seguridad e Higiene que buscan precisión, trazabilidad y velocidad al gestionar el
            riesgo psicosocial.
          </p>
        </div>

        {/* Feature clusters - same card format for both sections */}
        <div className="grid gap-8 lg:grid-cols-1 xl:grid-cols-2 mb-10 isolate items-stretch">
          <Card className="portal-surface-strong shimmer-border h-full border-0">
            <CardContent className="p-6 space-y-6 h-full flex flex-col">
              <div className="text-center space-y-2">
                <h3 className="text-xl font-semibold text-foreground">Operación</h3>
                <p className="text-sm text-muted-foreground">
                  Flujos que reducen tiempo de despliegue
                </p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 flex-1 content-start">
                {operacionFeatures.map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div key={feature.title} className={`group relative z-0 ${sectionCardClass}`}>
                      <div
                        className={`w-12 h-12 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center shrink-0`}
                      >
                        <IconComponent className="w-6 h-6 text-foreground" />
                      </div>
                      <h4 className="text-base font-semibold text-foreground">{feature.title}</h4>
                      <p className="text-sm text-muted-foreground">{feature.description}</p>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          <Card className="portal-surface-strong shimmer-border h-full border-0">
            <CardContent className="p-6 space-y-6 h-full flex flex-col">
              <div className="text-center space-y-2">
                <h3 className="text-xl font-semibold text-foreground">Seguridad & Acceso</h3>
                <p className="text-sm text-muted-foreground">
                  Confianza en cada interacción
                </p>
              </div>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 flex-1 content-start">
                {seguridadFeatures.map((feature) => {
                  const IconComponent = feature.icon;
                  return (
                    <div key={feature.title} className={`group relative z-0 ${sectionCardClass}`}>
                      <div
                        className={`w-12 h-12 rounded-lg bg-gradient-to-br ${feature.accent} flex items-center justify-center shrink-0`}
                      >
                        <IconComponent className="w-6 h-6 text-foreground" />
                      </div>
                      <h4 className="text-base font-semibold text-foreground">{feature.title}</h4>
                      <p className="text-sm text-muted-foreground">{feature.description}</p>
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
