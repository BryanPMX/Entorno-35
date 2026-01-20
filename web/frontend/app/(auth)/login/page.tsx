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
    <div className="min-h-screen flex items-center justify-center bg-gray-50 px-4">
      <Card className="w-full max-w-md">
        <CardHeader className="space-y-1">
          <CardTitle className="text-2xl font-bold text-center">
            Entorno 35
          </CardTitle>
          <CardDescription className="text-center">
            Sign in with your company RFC
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
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
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <Button type="submit" className="w-full" disabled={form.formState.isSubmitting}>
                {form.formState.isSubmitting ? "Signing in..." : "Sign In"}
              </Button>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}

