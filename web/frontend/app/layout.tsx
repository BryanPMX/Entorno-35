import type { Metadata } from "next";
import "./globals.css";
import { QueryProvider } from "@/components/providers/query-provider";
import { Toaster } from "sonner";

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

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className="antialiased"
        suppressHydrationWarning={true}
      >
        <QueryProvider>
          {children}
          <Toaster position="top-center" />
        </QueryProvider>
      </body>
    </html>
  );
}
