import type { Metadata } from "next";

/**
 * Marketing Layout
 *
 * Clean, professional layout for marketing pages.
 * Separate from dashboard layout for different branding and navigation.
 */
export const metadata: Metadata = {
  title: "Entorno 35 - NOM-035 Compliance Platform",
  description: "Professional psychosocial risk assessment platform compliant with NOM-035 STPS 2018. Streamline your organization's mental health compliance with automated assessments and comprehensive reporting.",
  keywords: ["NOM-035", "psychosocial risk", "workplace assessment", "compliance", "mental health", "Mexico", "STPS"],
  openGraph: {
    title: "Entorno 35 - NOM-035 Compliance Platform",
    description: "Professional psychosocial risk assessment platform for Mexican organizations",
    type: "website",
  },
};

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-white">
      {children}
    </div>
  );
}