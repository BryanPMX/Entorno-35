"use client";

import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { authService } from "@/services/auth.service";
import { useAuthStore } from "@/lib/store/auth-store";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
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
const loginSchema = z
  .object({
    identifier: z.string().min(1, "Identifier is required"),
    type: z.enum(["COMPANY", "STAFF"]),
    company_id: z.string().optional(),
  })
  .refine(
    (data) => {
      // company_id is required when type is STAFF
      if (data.type === "STAFF") {
        return data.company_id && data.company_id.trim().length > 0;
      }
      return true;
    },
    {
      message: "Company ID is required for staff login",
      path: ["company_id"],
    }
  );

type LoginFormValues = z.infer<typeof loginSchema>;

/**
 * Login Page
 * 
 * Authenticates users (Company or Staff) and redirects to dashboard on success.
 * - Company login: Requires RFC identifier
 * - Staff login: Requires CURP identifier + company_id
 */
export default function LoginPage() {
  const router = useRouter();
  const { login: storeLogin } = useAuthStore();

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      identifier: "",
      type: "COMPANY",
      company_id: "",
    },
  });

  const loginType = form.watch("type");

  const onSubmit = async (values: LoginFormValues) => {
    try {
      // Prepare login request (match backend API)
      const loginRequest = {
        identifier: values.identifier.trim(),
        type: values.type,
        ...(values.type === "STAFF" && values.company_id
          ? { company_id: values.company_id.trim() }
          : {}),
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
            Sign in to your account
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
              {/* Login Type Selection */}
              <FormField
                control={form.control}
                name="type"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Login Type</FormLabel>
                    <FormControl>
                      <div className="flex gap-4">
                        <Button
                          type="button"
                          variant={field.value === "COMPANY" ? "default" : "outline"}
                          onClick={() => {
                            field.onChange("COMPANY");
                            form.setValue("company_id", "");
                          }}
                          className="flex-1"
                        >
                          Company
                        </Button>
                        <Button
                          type="button"
                          variant={field.value === "STAFF" ? "default" : "outline"}
                          onClick={() => field.onChange("STAFF")}
                          className="flex-1"
                        >
                          Staff
                        </Button>
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Identifier Field (RFC for COMPANY, CURP for STAFF) */}
              <FormField
                control={form.control}
                name="identifier"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>
                      {loginType === "COMPANY" ? "RFC" : "CURP"}
                    </FormLabel>
                    <FormControl>
                      <Input
                        placeholder={
                          loginType === "COMPANY"
                            ? "ABC123456789"
                            : "CURP12345678901234"
                        }
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {/* Company ID Field (only for STAFF) */}
              {loginType === "STAFF" && (
                <FormField
                  control={form.control}
                  name="company_id"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Company ID (UUID)</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="550e8400-e29b-41d4-a716-446655440000"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                      <p className="text-sm text-gray-500">
                        Required for staff login
                      </p>
                    </FormItem>
                  )}
                />
              )}

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

