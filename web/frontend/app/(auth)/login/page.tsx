"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { authService } from "@/services/auth.service";
import { useAuthStore } from "@/lib/store/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { toast } from "sonner";

// Login form schema - matches backend API requirements
const loginSchema = z.object({
  identifier: z.string().min(1, "RFC is required"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

/**
 * Login Page
 *
 * Authenticates company administrators and redirects to dashboard on success.
 * - Company login: Requires RFC identifier
 */
export default function LoginPage() {
  const router = useRouter();
  const { login: storeLogin } = useAuthStore();

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      identifier: "",
    },
  });

  const onSubmit = async (values: LoginFormValues) => {
    try {
      // Prepare login request (match backend API)
      const loginRequest = {
        identifier: values.identifier.trim(),
        type: "COMPANY" as const,
      };

      // Call auth service
      const response = await authService.login(loginRequest);

      // Update auth store
      storeLogin(response);

      // Redirect to dashboard
      router.push("/dashboard");
      toast.success("Login successful!");
    } catch (error: unknown) {
      console.error("Login error:", error);
      toast.error("Invalid credentials. Please try again.");
    }
  };

  return (
    <div className="auth-shell portal-grid relative flex min-h-screen items-center justify-center overflow-hidden px-4 py-10">
      <div className="pointer-events-none absolute -left-28 top-6 h-72 w-72 rounded-full bg-[color:var(--brand-start-soft)] blur-3xl" />
      <div className="pointer-events-none absolute -right-24 bottom-0 h-72 w-72 rounded-full bg-[color:var(--brand-end-soft)] blur-3xl" />

      <Card className="portal-surface-strong w-full max-w-md border-0">
        <CardHeader className="space-y-4 pb-4">
          <div className="mx-auto flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-[var(--gradient-start)] to-[var(--gradient-end)] text-lg font-bold text-white shadow-lg">
            35
          </div>
          <div className="space-y-1 text-center">
            <CardTitle className="text-2xl font-bold tracking-tight">Entorno 35</CardTitle>
            <CardDescription>
              Sign in with your company RFC
            </CardDescription>
          </div>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-5">
              {/* RFC Field */}
              <FormField
                control={form.control}
                name="identifier"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>RFC</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="ABC123456789"
                        className="bg-background/80"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button
                type="submit"
                className="w-full bg-gradient-to-r from-[var(--gradient-start)] to-[var(--gradient-end)] text-white shadow-md hover:brightness-110"
                disabled={form.formState.isSubmitting}
              >
                {form.formState.isSubmitting ? "Signing in..." : "Sign In"}
              </Button>
              <p className="text-center text-xs text-muted-foreground">
                NOM-035 STPS 2018 compliance platform
              </p>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
