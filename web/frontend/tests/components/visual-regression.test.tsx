import { describe, expect, it, vi } from "vitest";
import { render } from "@testing-library/react";
import type { ReactNode } from "react";
import LoginPage from "@/app/(auth)/login/page";
import DashboardLayout from "@/app/dashboard/layout";
import PublicAssessmentLayout from "@/app/(public)/assessment/[token]/layout";
import { MarketingNavigation } from "@/components/marketing/navigation";
import { MarketingFooter } from "@/components/marketing/footer";

const mockPush = vi.fn();
const mockStoreLogin = vi.fn();
const mockStoreLogout = vi.fn();

vi.mock("next/navigation", () => ({
  useRouter: () => ({
    push: mockPush,
  }),
  usePathname: () => "/dashboard/assessments",
}));

vi.mock("@/components/layout/auth-guard", () => ({
  AuthGuard: ({ children }: { children: ReactNode }) => <>{children}</>,
}));

vi.mock("@/lib/store/auth-store", () => ({
  useAuthStore: vi.fn(() => ({
    login: mockStoreLogin,
    logout: mockStoreLogout,
  })),
}));

vi.mock("sonner", () => ({
  toast: {
    success: vi.fn(),
    error: vi.fn(),
  },
}));

describe("Visual Regression", () => {
  it("matches login page snapshot", () => {
    const { container } = render(<LoginPage />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it("matches dashboard shell snapshot", () => {
    const { container } = render(
      <DashboardLayout>
        <div>Dashboard content</div>
      </DashboardLayout>
    );
    expect(container.firstChild).toMatchSnapshot();
  });

  it("matches public assessment shell snapshot", () => {
    const { container } = render(
      <PublicAssessmentLayout>
        <div>Assessment content</div>
      </PublicAssessmentLayout>
    );
    expect(container.firstChild).toMatchSnapshot();
  });

  it("matches marketing header snapshot", () => {
    const { container } = render(<MarketingNavigation />);
    expect(container.firstChild).toMatchSnapshot();
  });

  it("matches marketing footer snapshot", () => {
    const { container } = render(<MarketingFooter />);
    expect(container.firstChild).toMatchSnapshot();
  });
});
