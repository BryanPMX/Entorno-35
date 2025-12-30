import type { User } from "@/types/backend";

/**
 * Decode JWT token payload (without verification)
 * This is a simple base64 decode - for production, use a proper JWT library
 */
export function decodeJWT(token: string): User | null {
  try {
    const parts = token.split(".");
    if (parts.length !== 3) {
      return null;
    }

    // Decode payload (second part)
    const payload = parts[1];
    const decoded = JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));

    return {
      company_id: decoded.company_id || "",
      staff_id: decoded.staff_id || undefined,
      email: decoded.email || undefined,
      role: decoded.role === "staff" ? "staff" : "company",
    };
  } catch (error) {
    console.error("Failed to decode JWT:", error);
    return null;
  }
}
