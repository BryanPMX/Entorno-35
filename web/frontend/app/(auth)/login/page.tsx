"use client";

import { FormEvent, useMemo, useState } from "react";
import Image from "next/image";
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
  const [fieldError, setFieldError] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const normalizedIdentifier = useMemo(() => identifier.trim().toUpperCase(), [identifier]);

  const validateRFC = (value: string): string | null => {
    if (!value) return "RFC is required";
    if (value.length < 12) return "RFC must contain at least 12 characters";
    return null;
  };

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const validationError = validateRFC(normalizedIdentifier);
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
      };

      const response = await authService.login(loginRequest);
      useAuthStore.getState().login(response);
      setFeedback({
        type: "success",
        message: "Access validated. Redirecting to dashboard...",
      });
      router.push("/dashboard");
    } catch (error: unknown) {
      void error;
      setFeedback({
        type: "error",
        message: "Invalid credentials. Please try again.",
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
                <Image src="/logo.png" alt="Entorno 35" width={48} height={48} className="h-12 w-auto object-contain" priority />
                <div className="space-y-2">
                  <h1 className="text-3xl font-semibold tracking-tight text-foreground sm:text-4xl">
                    Entorno 35
                  </h1>
                  <p className="max-w-md text-sm leading-6 text-muted-foreground sm:text-base">
                    Centralized NOM-035 access for compliance teams. Securely sign in to manage staff, assessments, and reporting.
                  </p>
                </div>
              </div>

              <div className="space-y-3 text-sm">
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Unified dashboards with risk tracking by department and demographic groups.
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Secure links for staff assessments with real-time progress visibility.
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Export-ready reporting aligned with NOM-035 STPS 2018 requirements.
                </p>
              </div>
            </div>
          </div>

          <div className="px-7 py-9 sm:px-10 sm:py-12">
            <div className="space-y-6">
              <header className="space-y-1">
                <h2 className="text-2xl font-semibold tracking-tight text-foreground">Sign In</h2>
                <p className="text-sm text-muted-foreground">Use your company RFC to enter the platform.</p>
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

                <button
                  type="submit"
                  className="inline-flex h-11 w-full items-center justify-center rounded-lg bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] px-4 text-sm font-semibold text-white shadow-md transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-70"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Signing in..." : "Sign In"}
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
                NOM-035 STPS 2018 compliance platform
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
