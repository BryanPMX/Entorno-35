import Image from "next/image";
import Link from "next/link";

/**
 * Marketing Footer
 *
 * Professional footer with links, contact info, and legal information.
 */
export function MarketingFooter() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="portal-footer-shell relative overflow-hidden text-white">
      <div className="absolute inset-0 bg-aurora opacity-35 pointer-events-none" />
      <div className="absolute inset-0 soft-grid opacity-20 pointer-events-none" />
      <div className="absolute -left-24 top-0 h-64 w-64 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl pointer-events-none" />
      <div className="absolute -right-24 bottom-0 h-72 w-72 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl pointer-events-none" />
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12 relative">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Company Info */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <Image
                src="/logo.png"
                alt="Entorno 35"
                width={812}
                height={293}
                className="h-10 w-auto object-contain"
              />
            </div>
            <p className="text-white/70 mb-4 max-w-md">
              Plataforma líder en evaluación de riesgos psicosociales conforme a NOM-035 STPS 2018.
              Ayudamos a organizaciones mexicanas a proteger la salud mental de sus colaboradores.
            </p>
          </div>

          {/* Product Links */}
          <div>
            <h3 className="font-semibold mb-4">Producto</h3>
            <ul className="space-y-2">
              <li>
                <Link href="#features" className="text-white/70 hover:text-white transition-colors">
                  Características
                </Link>
              </li>
              <li>
                <Link href="#compliance" className="text-white/70 hover:text-white transition-colors">
                  Cumplimiento NOM-035
                </Link>
              </li>
              <li>
                <Link href="#pricing" className="text-white/70 hover:text-white transition-colors">
                  Precios
                </Link>
              </li>
              <li>
                <Link href="/login" className="text-white/70 hover:text-white transition-colors">
                  Iniciar Sesión
                </Link>
              </li>
            </ul>
          </div>

          {/* Support Links */}
          <div>
            <h3 className="font-semibold mb-4">Soporte</h3>
            <ul className="space-y-2">
              <li>
                <a href="mailto:soporte@entorno35.com" className="text-white/70 hover:text-white transition-colors">
                  Soporte Técnico
                </a>
              </li>
              <li>
                <a href="https://docs.entorno35.com" className="text-white/70 hover:text-white transition-colors">
                  Documentación
                </a>
              </li>
              <li>
                <Link href="#contact" className="text-white/70 hover:text-white transition-colors">
                  Contacto
                </Link>
              </li>
              <li>
                <a href="/privacy" className="text-white/70 hover:text-white transition-colors">
                  Privacidad
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="border-t border-white/10 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
          <p className="text-white/70 text-sm">
            © {currentYear} Entorno 35. Todos los derechos reservados.
          </p>
          <div className="flex space-x-6 mt-4 md:mt-0">
            <a href="/terms" className="text-white/70 hover:text-white text-sm transition-colors">
              Términos de Servicio
            </a>
            <a href="/privacy" className="text-white/70 hover:text-white text-sm transition-colors">
              Política de Privacidad
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
