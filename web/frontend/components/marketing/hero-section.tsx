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
    <section className="relative overflow-hidden bg-gradient-to-b from-white via-blue-50/60 to-white py-20 lg:py-28 dark:from-background dark:via-background dark:to-background">
      <div className="absolute inset-0 -z-10 bg-aurora opacity-70 dark:opacity-40" />
      <div className="absolute inset-0 -z-10 soft-grid" />
      <div className="absolute -left-24 top-10 h-64 w-64 rounded-full bg-gradient-to-br from-blue-300/50 to-purple-400/40 blur-3xl" />
      <div className="absolute -right-32 -bottom-10 h-80 w-80 rounded-full bg-gradient-to-tl from-purple-500/40 to-indigo-400/35 blur-3xl" />

      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid lg:grid-cols-[1.2fr_0.8fr] gap-12 items-center">
          <div className="text-center lg:text-left space-y-8 relative">
            <Badge className="badge-glow px-4 py-2 inline-flex items-center gap-2 hover-shine shadow-lg">
              <Shield className="w-4 h-4" />
              Certificada NOM-035 STPS 2018
            </Badge>

            <div className="space-y-4">
              <h1 className="text-4xl lg:text-6xl font-bold leading-tight text-gray-900 dark:text-white">
                Evalúa y gestiona el{" "}
                <span className="gradient-text block">Riesgo Psicosocial</span>
                con precisión ejecutiva
              </h1>
              <p className="text-lg lg:text-xl text-gray-600 dark:text-gray-300 max-w-2xl leading-relaxed mx-auto lg:mx-0">
                Operacionaliza el cumplimiento NOM-035 con evaluaciones automatizadas,
                reportes listos para auditoría y analítica que combina people analytics con
                seguridad laboral.
              </p>
            </div>

            <div className="flex flex-col sm:flex-row gap-4 justify-center lg:justify-start">
              <Button size="lg" className="px-8 py-4 text-lg hover-shine shadow-lg shadow-blue-500/20" asChild>
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

            <div className="flex flex-col sm:flex-row sm:items-center gap-4 text-sm text-gray-600 dark:text-gray-300">
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <CheckCircle className="w-4 h-4 text-green-500" />
                100% compatible NOM-035
              </div>
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <Users className="w-4 h-4 text-blue-500" />
                500+ empresas activas
              </div>
              <div className="flex items-center gap-2 glass-panel px-3 py-2 rounded-full shadow-sm">
                <FileText className="w-4 h-4 text-purple-500" />
                Reportes PDF listos para STPS
              </div>
            </div>
          </div>

          <Card className="glass-panel-strong shimmer-border shadow-xl">
            <CardContent className="p-6 space-y-6">
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <span className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-gradient-to-br from-blue-600 to-purple-600 text-white shadow-lg animate-[floaty_6s_ease-in-out_infinite]">
                    <Sparkles className="w-5 h-5" />
                  </span>
                  <div>
                    <p className="text-sm text-gray-500 dark:text-gray-300">Tiempo promedio por evaluación</p>
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
                    <p className="text-2xl font-bold text-gray-900 dark:text-white">{item.value}</p>
                    <p className="text-xs text-gray-600 dark:text-gray-300">{item.label}</p>
                  </div>
                ))}
              </div>

              <div className="grid grid-cols-2 gap-3 text-sm text-gray-700 dark:text-gray-200">
                {[
                  "Links seguros con expiración controlada",
                  "Analítica por departamento y antigüedad",
                  "Métricas listadas para auditoría STPS",
                  "Exportación PDF y CSV en un click",
                ].map((benefit) => (
                  <div key={benefit} className="flex items-start gap-2">
                    <CheckCircle className="w-4 h-4 text-green-500 mt-0.5" />
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
