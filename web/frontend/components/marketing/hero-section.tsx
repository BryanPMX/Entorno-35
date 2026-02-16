import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { CheckCircle, Shield, Users, FileText, Sparkles } from "lucide-react";

/**
 * Hero Section
 *
 * Main landing section with compelling value proposition and CTAs.
 */
export function HeroSection() {
  return (
    <section className="relative overflow-hidden py-20 lg:py-28">
      <div className="absolute inset-0 -z-10 bg-aurora opacity-70 dark:opacity-40" />
      <div className="absolute inset-0 -z-10 soft-grid" />
      <div className="absolute -left-24 top-10 h-64 w-64 rounded-full bg-[color:var(--nom-bajo-soft)] blur-3xl" />
      <div className="absolute -right-32 -bottom-10 h-80 w-80 rounded-full bg-[color:var(--nom-nulo-soft)] blur-3xl" />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid lg:grid-cols-[1.2fr_0.8fr] gap-12 items-center">
          <div className="text-center lg:text-left space-y-8 relative">
            <Badge className="badge-glow px-4 py-2 inline-flex items-center gap-2 hover-shine shadow-lg">
              <Shield className="w-4 h-4" />
              Certificada NOM-035 STPS 2018
            </Badge>

            <div className="space-y-4">
              <h1 className="text-4xl lg:text-6xl font-bold leading-tight text-foreground">
                Evalúa y gestiona el{" "}
                <span className="gradient-text block">Riesgo Psicosocial</span>
                con precisión ejecutiva
              </h1>
              <p className="text-lg lg:text-xl text-muted-foreground max-w-2xl leading-relaxed mx-auto lg:mx-0">
                Operacionaliza el cumplimiento NOM-035 con evaluaciones automatizadas,
                reportes listos para auditoría y analítica que combina people analytics con
                seguridad laboral.
              </p>
            </div>

            <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
              <Button size="lg" className="px-8 py-4 text-lg hover-shine bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] text-white shadow-lg shadow-[rgb(2_132_199_/_24%)]" asChild>
                <Link href="#pricing">
                  Ver planes y precios
                  <span className="ml-2">→</span>
                </Link>
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="px-8 py-4 text-lg glass-panel border border-white/40 text-primary hover:border-primary/60"
                asChild
              >
                <Link href="#features">Explorar características</Link>
              </Button>
            </div>

            <div className="flex flex-col sm:flex-row sm:items-center gap-4 text-sm text-muted-foreground">
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <CheckCircle className="w-4 h-4 text-[var(--nom-nulo)]" />
                100% compatible NOM-035
              </div>
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <Users className="w-4 h-4 text-[var(--nom-bajo)]" />
                500+ empresas activas
              </div>
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <FileText className="w-4 h-4 text-[var(--nom-medio)]" />
                Reportes PDF listos para STPS
              </div>
            </div>
          </div>

          <Card className="portal-surface-strong shimmer-border border-0 shadow-xl">
            <CardContent className="p-6 space-y-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-[var(--gradient-start)] to-[var(--gradient-end)] text-white shadow-lg animate-[floaty_6s_ease-in-out_infinite]">
                    <Sparkles className="w-5 h-5" />
                  </span>
                  <div>
                    <p className="text-sm text-muted-foreground">Tiempo promedio por evaluación</p>
                    <p className="text-xl font-semibold">15 minutos</p>
                  </div>
                </div>
                <Badge variant="secondary" className="px-3 py-1 text-xs uppercase tracking-wide">
                  Express
                </Badge>
              </div>

              <div className="grid grid-cols-3 gap-4">
                {[
                  { label: "Tasa de adopción", value: "98%" },
                  { label: "Satisfacción", value: "4.9★" },
                  { label: "Disponibilidad", value: "24/7" },
                ].map((item) => (
                  <div key={item.label} className="glass-panel rounded-xl p-4 text-center space-y-2 border border-white/20">
                    <p className="text-2xl font-bold text-foreground">{item.value}</p>
                    <p className="text-xs text-muted-foreground">{item.label}</p>
                  </div>
                ))}
              </div>

              <div className="grid grid-cols-2 gap-3 text-sm text-foreground/90">
                {[
                  "Links seguros con expiración controlada",
                  "Analítica por departamento y antigüedad",
                  "Métricas listadas para auditoría STPS",
                  "Exportación PDF y CSV en un click",
                ].map((benefit) => (
                  <div key={benefit} className="flex items-start gap-2">
                    <CheckCircle className="w-4 h-4 text-[var(--nom-nulo)] mt-0.5" />
                    <span>{benefit}</span>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}
