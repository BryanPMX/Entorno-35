import Link from "next/link";

/**
 * Marketing Footer
 *
 * Professional footer with links, contact info, and legal information.
 */
export function MarketingFooter() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gradient-to-b from-gray-900 via-gray-950 to-black text-white">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Company Info */}
          <div className="col-span-1 md:col-span-2">
            <div className="flex items-center space-x-2 mb-4">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center shadow-md">
                <span className="text-white font-bold text-sm">35</span>
              </div>
              <span className="font-semibold text-xl">Entorno 35</span>
            </div>
            <p className="text-gray-400 mb-4 max-w-md">
              Plataforma líder en evaluación de riesgos psicosociales conforme a NOM-035 STPS 2018.
              Ayudamos a organizaciones mexicanas a proteger la salud mental de sus colaboradores.
            </p>
            <div className="flex space-x-4">
              <a
                href="https://linkedin.com/company/entorno35"
                className="text-gray-400 hover:text-white transition-colors"
                target="_blank"
                rel="noopener noreferrer"
              >
                LinkedIn
              </a>
              <a
                href="mailto:contacto@entorno35.com"
                className="text-gray-400 hover:text-white transition-colors"
              >
                Email
              </a>
            </div>
          </div>

          {/* Product Links */}
          <div>
            <h3 className="font-semibold mb-4">Producto</h3>
            <ul className="space-y-2">
              <li>
                <Link href="#features" className="text-gray-400 hover:text-white transition-colors">
                  Características
                </Link>
              </li>
              <li>
                <Link href="#compliance" className="text-gray-400 hover:text-white transition-colors">
                  Cumplimiento NOM-035
                </Link>
              </li>
              <li>
                <Link href="#pricing" className="text-gray-400 hover:text-white transition-colors">
                  Precios
                </Link>
              </li>
              <li>
                <Link href="/login" className="text-gray-400 hover:text-white transition-colors">
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
                <a href="mailto:soporte@entorno35.com" className="text-gray-400 hover:text-white transition-colors">
                  Soporte Técnico
                </a>
              </li>
              <li>
                <a href="https://docs.entorno35.com" className="text-gray-400 hover:text-white transition-colors">
                  Documentación
                </a>
              </li>
              <li>
                <Link href="#contact" className="text-gray-400 hover:text-white transition-colors">
                  Contacto
                </Link>
              </li>
              <li>
                <a href="/privacy" className="text-gray-400 hover:text-white transition-colors">
                  Privacidad
                </a>
              </li>
            </ul>
          </div>
        </div>

        {/* Bottom Bar */}
        <div className="border-t border-white/10 mt-12 pt-8 flex flex-col md:flex-row justify-between items-center">
          <p className="text-gray-400 text-sm">
            © {currentYear} Entorno 35. Todos los derechos reservados.
          </p>
          <div className="flex space-x-6 mt-4 md:mt-0">
            <a href="/terms" className="text-gray-400 hover:text-white text-sm transition-colors">
              Términos de Servicio
            </a>
            <a href="/privacy" className="text-gray-400 hover:text-white text-sm transition-colors">
              Política de Privacidad
            </a>
            <span className="text-gray-600 text-sm">
              Certificado NOM-035 STPS 2018
            </span>
          </div>
        </div>
      </div>
    </footer>
  );
}
