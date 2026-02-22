"use client";

import { FormEvent, useMemo, useState } from "react";
import Image from "next/image";
import Link from "next/link";
import { useRouter } from "next/navigation";
type FeedbackState = {
  type: "error" | "success";
  message: string;
} | null;

/**
 * Login Page
 *
 * Authenticates company administrators and redirects to dashboard on success.
 * - Company login: Requires RFC identifier
 */
export default function LoginPage() {
  const router = useRouter();
  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [fieldError, setFieldError] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const normalizedIdentifier = useMemo(() => identifier.trim().toUpperCase(), [identifier]);

  const validateLogin = (value: string, passwordValue: string): string | null => {
    if (!value) return "El RFC es obligatorio.";
    if (value.length < 12) return "El RFC debe contener al menos 12 caracteres.";
    if (!passwordValue) return "La contraseña es obligatoria.";
    return null;
  };

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const validationError = validateLogin(normalizedIdentifier, password);
    if (validationError) {
      setFieldError(validationError);
      setFeedback(null);
      return;
    }

    setFieldError(null);
    setFeedback(null);
    setIsSubmitting(true);

    try {
      const [{ authService }, { useAuthStore }] = await Promise.all([
        import("@/services/auth.service"),
        import("@/lib/store/auth-store"),
      ]);

      const loginRequest = {
        identifier: normalizedIdentifier,
        type: "COMPANY" as const,
        password,
      };

      const response = await authService.login(loginRequest);
      useAuthStore.getState().login(response);
      setFeedback({
        type: "success",
        message: "Acceso validado. Redirigiendo al panel de control...",
      });
      router.push("/dashboard");
    } catch (error: unknown) {
      void error;
      setFeedback({
        type: "error",
        message: "Credenciales inválidas. Por favor, intente nuevamente.",
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="page-background page-background-4 portal-grid relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10 sm:px-6">
      <div className="pointer-events-none absolute -left-28 top-6 h-72 w-72 rounded-full bg-[color:var(--brand-start-soft)] blur-2xl md:blur-3xl" aria-hidden="true" />
      <div className="pointer-events-none absolute -right-24 bottom-0 h-72 w-72 rounded-full bg-[color:var(--brand-end-soft)] blur-2xl md:blur-3xl" aria-hidden="true" />

      <section className="portal-surface-strong login-surface-optimized relative z-10 w-full max-w-5xl overflow-hidden rounded-3xl border-0 shadow-2xl">
        <div className="grid lg:grid-cols-[1.1fr_1fr]">
          <div className="relative overflow-hidden border-b border-white/20 px-7 py-9 lg:border-b-0 lg:border-r lg:px-10 lg:py-12">
            <div className="pointer-events-none absolute -right-24 -top-24 h-64 w-64 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" aria-hidden="true" />
            <div className="pointer-events-none absolute -left-24 -bottom-24 h-64 w-64 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" aria-hidden="true" />
            <div className="relative space-y-8">
              <div className="space-y-5">
                <Image src="/logo.png" alt="Entorno 35" width={812} height={293} className="h-16 w-auto object-contain md:h-20" priority />
                <p className="max-w-md text-sm leading-6 text-muted-foreground sm:text-base">
                  Acceso centralizado a NOM-035 para equipos de cumplimiento. Inicie sesión de forma segura para gestionar personal, evaluaciones y reportes.
                </p>
              </div>

              <div className="space-y-3 text-sm">
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Paneles unificados con seguimiento de riesgos por departamento y grupo demográfico.
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Enlaces seguros para evaluaciones del personal con visibilidad del avance en tiempo real.
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Reportes listos para exportación, alineados con los requisitos de la NOM-035 STPS 2018.
                </p>
              </div>
            </div>
          </div>

          <div className="px-7 py-9 sm:px-10 sm:py-12">
            <div className="space-y-6">
              <header className="space-y-1">
                <h2 className="text-2xl font-semibold tracking-tight text-foreground">Iniciar sesión</h2>
                <p className="text-sm text-muted-foreground">Utilice el RFC de su empresa para ingresar a la plataforma.</p>
              </header>

              <form onSubmit={onSubmit} noValidate className="space-y-4">
                <div className="space-y-2">
                  <label htmlFor="identifier" className="block text-sm font-medium text-foreground">
                    RFC
                  </label>
                  <input
                    id="identifier"
                    name="identifier"
                    autoComplete="username"
                    placeholder="ABC123456789"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={identifier}
                    onChange={(event) => {
                      setIdentifier(event.target.value.toUpperCase());
                      if (fieldError) setFieldError(null);
                      if (feedback?.type === "error") setFeedback(null);
                    }}
                    aria-invalid={Boolean(fieldError)}
                    aria-describedby={fieldError ? "identifier-error" : undefined}
                    disabled={isSubmitting}
                  />
                  {fieldError ? (
                    <p id="identifier-error" className="text-xs text-destructive">
                      {fieldError}
                    </p>
                  ) : null}
                </div>

                <div className="space-y-2">
                  <label htmlFor="password" className="block text-sm font-medium text-foreground">
                    Contraseña
                  </label>
                  <input
                    id="password"
                    name="password"
                    autoComplete="current-password"
                    type="password"
                    placeholder="Tu contraseña de administrador"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={password}
                    onChange={(event) => {
                      setPassword(event.target.value);
                      if (fieldError) setFieldError(null);
                      if (feedback?.type === "error") setFeedback(null);
                    }}
                    disabled={isSubmitting}
                  />
                </div>

                <button
                  type="submit"
                  className="inline-flex h-11 w-full items-center justify-center rounded-lg bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] px-4 text-sm font-semibold text-white shadow-md transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-70"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Iniciando sesión..." : "Iniciar sesión"}
                </button>
              </form>

              {feedback ? (
                <p
                  className={
                    feedback.type === "success"
                      ? "rounded-lg border border-emerald-300/50 bg-emerald-500/10 px-3 py-2 text-xs text-emerald-800"
                      : "rounded-lg border border-destructive/35 bg-destructive/10 px-3 py-2 text-xs text-destructive"
                  }
                >
                  {feedback.message}
                </p>
              ) : null}

              <p className="text-center text-xs text-muted-foreground">
                Plataforma de cumplimiento NOM-035 STPS 2018.{" "}
                <Link href="/register" className="font-medium text-foreground hover:underline">
                  Crear cuenta
                </Link>
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
