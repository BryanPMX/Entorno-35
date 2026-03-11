"use client";

import { FormEvent, useState } from "react";
import { useRouter } from "next/navigation";
import { adminRefundService } from "@/services/admin-refund.service";
import { useAdminAuthStore } from "@/lib/store/admin-auth-store";

type FeedbackState = {
  type: "error" | "success";
  message: string;
} | null;

export default function AdminLoginPage() {
  const router = useRouter();
  const { login } = useAdminAuthStore();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [feedback, setFeedback] = useState<FeedbackState>(null);

  const onSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setFeedback(null);
    setIsSubmitting(true);

    try {
      const response = await adminRefundService.login({
        email: email.trim().toLowerCase(),
        password,
      });
      login(response.token);
      setFeedback({
        type: "success",
        message: "Acceso administrativo validado. Redirigiendo...",
      });
      router.push("/admin/refunds");
    } catch (error: unknown) {
      const fallbackMessage = "Credenciales inválidas o autenticación administrativa no configurada.";
      const apiMessage =
        typeof error === "object" && error && "response" in error
          ? ((error as { response?: { data?: { error?: string } } }).response?.data?.error ?? fallbackMessage)
          : fallbackMessage;
      setFeedback({
        type: "error",
        message: apiMessage,
      });
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="page-background page-background-4 portal-grid relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10 sm:px-6">
      <section className="portal-surface-strong relative z-10 w-full max-w-md rounded-3xl border-0 p-8 shadow-2xl sm:p-10">
        <header className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-muted-foreground">Entorno35 Billing Admin</p>
          <h1 className="text-2xl font-semibold tracking-tight text-foreground">Login Administrativo</h1>
          <p className="text-sm text-muted-foreground">
            Acceso restringido para gestión de solicitudes de reembolso.
          </p>
        </header>

        <form onSubmit={onSubmit} noValidate className="mt-6 space-y-4">
          <div className="space-y-2">
            <label htmlFor="email" className="block text-sm font-medium text-foreground">
              Correo Administrador
            </label>
            <input
              id="email"
              name="email"
              autoComplete="username"
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
              placeholder="admin@entorno35.com"
              disabled={isSubmitting}
              required
            />
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
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="h-11 w-full rounded-lg border border-input bg-background/80 px-3 text-sm shadow-xs outline-none transition focus:border-ring focus:ring-[3px] focus:ring-ring/40"
              placeholder="Contraseña de administrador"
              disabled={isSubmitting}
              required
            />
          </div>

          <button
            type="submit"
            className="inline-flex h-11 w-full items-center justify-center rounded-lg bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] px-4 text-sm font-semibold text-white shadow-md transition hover:brightness-110 disabled:cursor-not-allowed disabled:opacity-70"
            disabled={isSubmitting}
          >
            {isSubmitting ? "Validando..." : "Entrar"}
          </button>
        </form>

        {feedback ? (
          <p
            className={
              feedback.type === "success"
                ? "mt-4 rounded-lg border border-emerald-300/50 bg-emerald-500/10 px-3 py-2 text-xs text-emerald-800"
                : "mt-4 rounded-lg border border-destructive/35 bg-destructive/10 px-3 py-2 text-xs text-destructive"
            }
          >
            {feedback.message}
          </p>
        ) : null}
      </section>
    </div>
  );
}
