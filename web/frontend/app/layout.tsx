import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Entorno 35 - NOM-035 Compliance Platform",
  description: "Professional psychosocial risk assessment platform compliant with NOM-035 STPS 2018. Streamline your organization's mental health compliance with automated assessments and comprehensive reporting.",
  keywords: ["NOM-035", "psychosocial risk", "workplace assessment", "compliance", "mental health", "Mexico", "STPS"],
  icons: {
    icon: [{ url: "/tab-logo.png", type: "image/png" }],
    shortcut: [{ url: "/tab-logo.png", type: "image/png" }],
    apple: [{ url: "/tab-logo.png", type: "image/png" }],
  },
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
        {children}
      </body>
    </html>
  );
}
