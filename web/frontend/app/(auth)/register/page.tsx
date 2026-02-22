"use client";

import { FormEvent, Suspense, useEffect, useMemo, useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { useSearchParams } from "next/navigation";

type FeedbackState = {
  type: "error" | "info";
  message: string;
} | null;

type PlanType = "monthly" | "yearly";

function isPlan(value: string | null): value is PlanType {
  return value === "monthly" || value === "yearly";
}

function RegisterPageContent() {
  const searchParams = useSearchParams();

  const [companyName, setCompanyName] = useState("");
  const [rfc, setRFC] = useState("");
  const [address, setAddress] = useState("");
  const [employeeCount, setEmployeeCount] = useState("");
  const [adminEmail, setAdminEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [plan, setPlan] = useState<PlanType>("yearly");
  const [fieldError, setFieldError] = useState<string | null>(null);
  const [feedback, setFeedback] = useState<FeedbackState>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const planParam = searchParams.get("plan");
    if (isPlan(planParam)) {
      setPlan(planParam);
    }
    if (searchParams.get("canceled") === "1") {
      setFeedback({
        type: "info",
        message: "Checkout canceled. You can resume registration whenever you are ready.",
      });
    }
  }, [searchParams]);

  const normalizedRFC = useMemo(() => rfc.trim().toUpperCase(), [rfc]);

  const validateForm = (): string | null => {
    if (!companyName.trim()) return "Company name is required";
    if (!normalizedRFC) return "RFC is required";
    if (normalizedRFC.length < 12) return "RFC must contain at least 12 characters";
    if (!adminEmail.trim()) return "Admin email is required";
    if (!password) return "Password is required";
    if (password.length < 8) return "Password must contain at least 8 characters";
    if (password !== confirmPassword) return "Passwords do not match";
    if (employeeCount !== "" && Number(employeeCount) < 0) return "Employee count cannot be negative";
    return null;
  };

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const validationError = validateForm();
    if (validationError) {
      setFieldError(validationError);
      setFeedback(null);
      return;
    }

    setFieldError(null);
    setFeedback(null);
    setIsSubmitting(true);

    try {
      const { billingService } = await import("@/services/billing.service");
      const response = await billingService.createCheckoutSession({
        rfc: normalizedRFC,
        company_name: companyName.trim(),
        address: address.trim(),
        admin_email: adminEmail.trim(),
        password,
        plan,
        employee_count: employeeCount === "" ? undefined : Number(employeeCount),
      });

      if (!response.checkout_url) {
        throw new Error("Stripe checkout URL missing");
      }

      window.location.assign(response.checkout_url);
    } catch (error: unknown) {
      const fallbackMessage = "Unable to start checkout. Please verify your details and try again.";
      const message =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;

      setFeedback({
        type: "error",
        message,
      });
      setIsSubmitting(false);
    }
  };

  return (
    <div className="page-background page-background-4 portal-grid relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10 sm:px-6">
      <div className="pointer-events-none absolute -left-28 top-6 h-72 w-72 rounded-full bg-[color:var(--brand-start-soft)] blur-2xl md:blur-3xl" aria-hidden="true" />
      <div className="pointer-events-none absolute -right-24 bottom-0 h-72 w-72 rounded-full bg-[color:var(--brand-end-soft)] blur-2xl md:blur-3xl" aria-hidden="true" />

      <section className="portal-surface-strong login-surface-optimized relative z-10 w-full max-w-6xl overflow-hidden rounded-3xl border-0 shadow-2xl">
        <div className="grid lg:grid-cols-[1.1fr_1fr]">
          <div className="relative overflow-hidden border-b border-white/20 px-7 py-9 lg:border-b-0 lg:border-r lg:px-10 lg:py-12">
            <div className="pointer-events-none absolute -right-24 -top-24 h-64 w-64 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" aria-hidden="true" />
            <div className="pointer-events-none absolute -left-24 -bottom-24 h-64 w-64 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" aria-hidden="true" />
            <div className="relative space-y-8">
              <div className="space-y-5">
                <Image src="/logo.png" alt="Entorno 35" width={812} height={293} className="h-16 w-auto object-contain md:h-20" priority />
                <p className="max-w-md text-sm leading-6 text-muted-foreground sm:text-base">
                  Create your organization account, choose your plan, and complete secure Stripe checkout to activate access.
                </p>
              </div>

              <div className="space-y-3 text-sm">
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Monthly plan: $1,000 MXN per month.
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Yearly plan: $6,000 MXN per year (50% effective savings).
                </p>
                <p className="rounded-lg border border-white/25 bg-background/55 px-4 py-3 text-foreground/95">
                  Credentials are activated immediately after successful payment confirmation.
                </p>
              </div>
            </div>
          </div>

          <div className="px-7 py-9 sm:px-10 sm:py-12">
            <div className="space-y-6">
              <header className="space-y-1">
                <h2 className="text-2xl font-semibold tracking-tight text-foreground">Create Account</h2>
                <p className="text-sm text-muted-foreground">Register your company and proceed to checkout.</p>
              </header>

              <form onSubmit={onSubmit} noValidate className="space-y-4">
                <div className="space-y-2">
                  <label htmlFor="company_name" className="block text-sm font-medium text-foreground">
                    Company Name
                  </label>
                  <input
                    id="company_name"
                    name="company_name"
                    placeholder="Mi Empresa SA de CV"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={companyName}
                    onChange={(event) => setCompanyName(event.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <label htmlFor="rfc" className="block text-sm font-medium text-foreground">
                    RFC
                  </label>
                  <input
                    id="rfc"
                    name="rfc"
                    autoComplete="organization"
                    placeholder="ABC123456789"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={rfc}
                    onChange={(event) => setRFC(event.target.value.toUpperCase())}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="space-y-2">
                  <label htmlFor="admin_email" className="block text-sm font-medium text-foreground">
                    Admin Email
                  </label>
                  <input
                    id="admin_email"
                    name="admin_email"
                    autoComplete="email"
                    type="email"
                    placeholder="admin@empresa.com"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={adminEmail}
                    onChange={(event) => setAdminEmail(event.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <label htmlFor="password" className="block text-sm font-medium text-foreground">
                      Password
                    </label>
                    <input
                      id="password"
                      name="password"
                      autoComplete="new-password"
                      type="password"
                      placeholder="At least 8 characters"
                      className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                      value={password}
                      onChange={(event) => setPassword(event.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>

                  <div className="space-y-2">
                    <label htmlFor="confirm_password" className="block text-sm font-medium text-foreground">
                      Confirm Password
                    </label>
                    <input
                      id="confirm_password"
                      name="confirm_password"
                      autoComplete="new-password"
                      type="password"
                      placeholder="Repeat password"
                      className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                      value={confirmPassword}
                      onChange={(event) => setConfirmPassword(event.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <label htmlFor="employee_count" className="block text-sm font-medium text-foreground">
                      Employee Count (optional)
                    </label>
                    <input
                      id="employee_count"
                      name="employee_count"
                      inputMode="numeric"
                      placeholder="50"
                      className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                      value={employeeCount}
                      onChange={(event) => setEmployeeCount(event.target.value)}
                      disabled={isSubmitting}
                    />
                  </div>

                  <div className="space-y-2">
                    <span className="block text-sm font-medium text-foreground">Plan</span>
                    <div className="grid grid-cols-2 gap-2">
                      <button
                        type="button"
                        className={`h-11 rounded-lg border px-3 text-sm font-medium transition ${
                          plan === "monthly"
                            ? "border-transparent bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] text-white"
                            : "border-input bg-background/80 text-foreground"
                        }`}
                        onClick={() => setPlan("monthly")}
                        disabled={isSubmitting}
                      >
                        Monthly
                      </button>
                      <button
                        type="button"
                        className={`h-11 rounded-lg border px-3 text-sm font-medium transition ${
                          plan === "yearly"
                            ? "border-transparent bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] text-white"
                            : "border-input bg-background/80 text-foreground"
                        }`}
                        onClick={() => setPlan("yearly")}
                        disabled={isSubmitting}
                      >
                        Yearly
                      </button>
                    </div>
                  </div>
                </div>

                <div className="space-y-2">
                  <label htmlFor="address" className="block text-sm font-medium text-foreground">
                    Address (optional)
                  </label>
                  <input
                    id="address"
                    name="address"
                    autoComplete="street-address"
                    placeholder="Av. Ejemplo 123, CDMX"
                    className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
                    value={address}
                    onChange={(event) => setAddress(event.target.value)}
                    disabled={isSubmitting}
                  />
                </div>

                {fieldError ? (
                  <p className="text-xs text-destructive">{fieldError}</p>
                ) : null}

                <button
                  type="submit"
                  className="inline-flex h-11 w-full items-center justify-center rounded-lg bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] px-4 text-sm font-semibold text-white shadow-md transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-70"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Redirecting to checkout..." : "Continue to Secure Checkout"}
                </button>
              </form>

              {feedback ? (
                <p
                  className={
                    feedback.type === "info"
                      ? "rounded-lg border border-blue-300/50 bg-blue-500/10 px-3 py-2 text-xs text-blue-900"
                      : "rounded-lg border border-destructive/35 bg-destructive/10 px-3 py-2 text-xs text-destructive"
                  }
                >
                  {feedback.message}
                </p>
              ) : null}

              <p className="text-center text-xs text-muted-foreground">
                Already have credentials?{" "}
                <Link href="/login" className="font-medium text-foreground hover:underline">
                  Sign in
                </Link>
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}

export default function RegisterPage() {
  return (
    <Suspense
      fallback={
        <div className="page-background page-background-4 portal-grid flex min-h-screen items-center justify-center">
          <p className="text-sm text-muted-foreground">Loading registration...</p>
        </div>
      }
    >
      <RegisterPageContent />
    </Suspense>
  );
}
