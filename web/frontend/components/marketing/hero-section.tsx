import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CheckCircle, Shield, Users, FileText } from "lucide-react";

/**
 * Hero Section
 *
 * Main landing section with compelling value proposition and CTAs.
 */
export function HeroSection() {
  return (
    <section className="relative bg-gradient-to-br from-blue-50 via-white to-blue-50 py-20 lg:py-28">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="text-center">
          {/* Badge */}
          <Badge variant="secondary" className="mb-6 px-4 py-2">
            <Shield className="w-4 h-4 mr-2" />
            Plataforma Certificada NOM-035 STPS 2018
          </Badge>

          {/* Headline */}
          <h1 className="text-4xl lg:text-6xl font-bold text-gray-900 mb-6 leading-tight">
            Evalúa y Gestiona el
            <span className="text-blue-600 block">Riesgo Psicosocial</span>
            en tu Organización
          </h1>

          {/* Subheadline */}
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto leading-relaxed">
            Simplifica el cumplimiento con NOM-035 mediante evaluaciones automatizadas,
            reportes profesionales en PDF y análisis detallados del riesgo psicosocial laboral.
          </p>

          {/* Social Proof */}
          <div className="flex items-center justify-center space-x-6 mb-8 text-sm text-gray-500">
            <div className="flex items-center">
              <CheckCircle className="w-4 h-4 text-green-500 mr-1" />
              <span>100% Compatible con NOM-035</span>
            </div>
            <div className="flex items-center">
              <Users className="w-4 h-4 text-blue-500 mr-1" />
              <span>+500 Empresas Confían en Nosotros</span>
            </div>
            <div className="flex items-center">
              <FileText className="w-4 h-4 text-purple-500 mr-1" />
              <span>Reportes PDF Profesionales</span>
            </div>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
            <Button size="lg" className="px-8 py-4 text-lg" asChild>
              <Link href="#pricing">
                Ver Planes y Precios
                <span className="ml-2">→</span>
              </Link>
            </Button>
            <Button size="lg" variant="outline" className="px-8 py-4 text-lg" asChild>
              <Link href="#features">
                Ver Características
              </Link>
            </Button>
          </div>

          {/* Trust Indicators */}
          <div className="text-sm text-gray-500">
            <p>Desde $1,000 MXN/mes • Configuración en minutos • Soporte técnico incluido</p>
          </div>
        </div>
      </div>

      {/* Background decoration */}
      <div className="absolute top-0 right-0 -z-10 opacity-10">
        <div className="w-96 h-96 bg-gradient-to-br from-blue-400 to-purple-500 rounded-full blur-3xl"></div>
      </div>
    </section>
  );
}