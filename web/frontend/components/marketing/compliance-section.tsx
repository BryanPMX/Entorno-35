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
    <section id="compliance" className="relative overflow-hidden py-20 lg:py-28">
      <div className="absolute inset-0 -z-10 bg-aurora opacity-70 dark:opacity-40" />
      <div className="absolute inset-0 -z-10 soft-grid" />
      <div className="absolute -left-24 top-10 h-64 w-64 rounded-full bg-[color:var(--nom-bajo-soft)] blur-3xl" />
      <div className="absolute -right-32 -bottom-10 h-80 w-80 rounded-full bg-[color:var(--nom-nulo-soft)] blur-3xl" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        {/* Section Header */}
        <div className="text-center mb-16">
          <Badge variant="outline" className="mb-4 glass-panel px-4 py-2">
            Cumplimiento NOM-035
          </Badge>
          <h2 className="text-3xl lg:text-4xl font-bold text-foreground mb-4">
            Norma Oficial Mexicana 035
          </h2>
          <p className="text-lg text-muted-foreground max-w-3xl mx-auto">
            Diseñada para cubrir requisitos legales, mantener trazabilidad y entregar evidencia clara durante auditorías.
          </p>
        </div>

        <div className="grid lg:grid-cols-2 gap-12 items-center mb-16">
          {/* Compliance Requirements */}
          <div className="portal-surface rounded-2xl p-8 border-0 shadow-lg">
            <h3 className="text-2xl font-bold text-foreground mb-6">
              Requisitos cubiertos
            </h3>
            <div className="space-y-6">
              {compliancePoints.map((point) => {
                const IconComponent = point.icon;
                return (
                  <div key={point.title} className="flex items-start gap-4">
                    <div className="flex-shrink-0">
                      <div className="w-10 h-10 rounded-lg bg-gradient-to-br from-[var(--nom-nulo)] to-[var(--nom-bajo)] text-white flex items-center justify-center shadow-md">
                        <IconComponent className="w-5 h-5" />
                      </div>
                    </div>
                    <div className="flex-1">
                      <div className="flex items-center gap-2 mb-1">
                        <h4 className="font-semibold text-foreground">
                          {point.title}
                        </h4>
                        <Badge variant="outline" className="text-xs">
                          {point.requirement}
                        </Badge>
                      </div>
                      <p className="text-muted-foreground">
                        {point.description}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          {/* Benefits */}
          <div className="portal-surface-strong rounded-2xl p-8 border-0 shadow-xl">
            <div className="flex items-center gap-2 text-[var(--nom-nulo)] mb-4">
              <CheckCircle className="w-5 h-5" />
              <p className="text-sm uppercase tracking-wide">Protección integral</p>
            </div>
            <h3 className="text-2xl font-bold text-foreground mb-6">
              Beneficios del cumplimiento
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {benefits.map((benefit) => (
                <div key={benefit} className="flex items-start gap-2">
                  <CheckCircle className="w-4 h-4 text-[var(--nom-nulo)] mt-0.5" />
                  <span className="text-foreground/90 text-sm">{benefit}</span>
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Soft transition to next sections */}
        <div className="h-16 relative mt-8">
          <div className="absolute inset-0 bg-gradient-to-b from-transparent via-background/60 to-background blur-[2px]" />
        </div>
      </div>
    </section>
  );
}
