import type { User } from "@/types/backend";
import { jwtDecode, type JwtPayload } from "jwt-decode";

interface EntornoJWTPayload extends JwtPayload {
  company_id?: string;
  staff_id?: string;
  email?: string;
  role?: string;
}

interface ParsedJWT {
  user: User;
  exp: number | null;
}

function decodePayload(token: string): EntornoJWTPayload | null {
  try {
    return jwtDecode<EntornoJWTPayload>(token);
  } catch {
    return null;
  }
}

function parseJWT(token: string): ParsedJWT | null {
  const payload = decodePayload(token);
  if (!payload || typeof payload.company_id !== "string" || payload.company_id.trim() === "") {
    return null;
  }

  return {
    user: {
      company_id: payload.company_id,
      staff_id: typeof payload.staff_id === "string" ? payload.staff_id : undefined,
      email: typeof payload.email === "string" ? payload.email : undefined,
      role: payload.role === "staff" ? "staff" : "company",
    },
    exp: typeof payload.exp === "number" ? payload.exp : null,
  };
}

/**
 * Decodes user claims from a JWT token.
 * Signature verification is done by the backend on each protected API request.
 */
export function decodeJWT(token: string): User | null {
  return parseJWT(token)?.user ?? null;
}

/**
 * Returns true when the token is malformed, missing exp, or already expired.
 */
export function isTokenExpired(token: string, nowMs = Date.now()): boolean {
  const parsed = parseJWT(token);
  if (!parsed || parsed.exp === null) {
    return true;
  }

  const nowSeconds = Math.floor(nowMs / 1000);
  return parsed.exp <= nowSeconds;
}
