"use client";

import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import type { VerifyCheckoutSessionResponse } from "@/types/backend";

type LoadState = "loading" | "ready" | "error";

function RegisterSuccessPageContent() {
  const searchParams = useSearchParams();
  const sessionID = searchParams.get("session_id");

  const [state, setState] = useState<LoadState>("loading");
  const [result, setResult] = useState<VerifyCheckoutSessionResponse | null>(null);
  const [errorMessage, setErrorMessage] = useState("");

  useEffect(() => {
    const verify = async () => {
      if (!sessionID) {
        setState("error");
        setErrorMessage("Missing checkout session ID. Please repeat registration.");
        return;
      }

      try {
        const { billingService } = await import("@/services/billing.service");
        const response = await billingService.verifyCheckoutSession(sessionID);
        setResult(response);
        setState("ready");
      } catch (error: unknown) {
        const fallbackMessage = "Unable to verify your payment right now. Please try again in a minute.";
        const message =
          typeof error === "object" && error && "response" in error
            ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
            : fallbackMessage;
        setState("error");
        setErrorMessage(message);
      }
    };

    void verify();
  }, [sessionID]);

  return (
    <div className="page-background page-background-4 portal-grid relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10 sm:px-6">
      <section className="portal-surface-strong relative z-10 w-full max-w-2xl rounded-3xl border-0 p-8 shadow-2xl sm:p-10">
        <h1 className="text-2xl font-semibold tracking-tight text-foreground">Checkout Confirmation</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          We are validating your subscription and provisioning access credentials.
        </p>

        {state === "loading" ? (
          <p className="mt-6 rounded-lg border border-blue-300/50 bg-blue-500/10 px-4 py-3 text-sm text-blue-900">
            Confirming payment with Stripe...
          </p>
        ) : null}

        {state === "error" ? (
          <>
            <p className="mt-6 rounded-lg border border-destructive/35 bg-destructive/10 px-4 py-3 text-sm text-destructive">
              {errorMessage}
            </p>
            <div className="mt-6 flex flex-wrap gap-3">
              <Link
                href="/register"
                className="inline-flex h-10 items-center justify-center rounded-lg border border-input px-4 text-sm font-medium text-foreground transition hover:bg-accent"
              >
                Back to Registration
              </Link>
            </div>
          </>
        ) : null}

        {state === "ready" && result ? (
          <div className="mt-6 space-y-4">
            {result.active ? (
              <p className="rounded-lg border border-emerald-300/50 bg-emerald-500/10 px-4 py-3 text-sm text-emerald-900">
                Subscription active. Your credentials are ready.
              </p>
            ) : (
              <p className="rounded-lg border border-amber-300/50 bg-amber-500/10 px-4 py-3 text-sm text-amber-900">
                Payment received but activation is still pending. Retry verification in a moment.
              </p>
            )}

            <div className="rounded-xl border border-white/20 bg-background/60 p-4 text-sm text-foreground">
              <p><strong>Company:</strong> {result.company_name}</p>
              <p><strong>Login RFC:</strong> {result.login_identifier}</p>
              {result.admin_email ? <p><strong>Admin Email:</strong> {result.admin_email}</p> : null}
              <p><strong>Password:</strong> The password you defined during registration.</p>
            </div>

            <div className="flex flex-wrap gap-3">
              <Link
                href="/login"
                className="inline-flex h-10 items-center justify-center rounded-lg bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] px-4 text-sm font-semibold text-white shadow-md transition hover:brightness-110"
              >
                Go to Login
              </Link>
              {!result.active ? (
                <Link
                  href={`/register/success?session_id=${encodeURIComponent(result.session_id)}`}
                  className="inline-flex h-10 items-center justify-center rounded-lg border border-input px-4 text-sm font-medium text-foreground transition hover:bg-accent"
                >
                  Retry Verification
                </Link>
              ) : null}
            </div>
          </div>
        ) : null}
      </section>
    </div>
  );
}

export default function RegisterSuccessPage() {
  return (
    <Suspense
      fallback={
        <div className="page-background page-background-4 portal-grid flex min-h-screen items-center justify-center">
          <p className="text-sm text-muted-foreground">Validating checkout...</p>
        </div>
      }
    >
      <RegisterSuccessPageContent />
    </Suspense>
  );
}
