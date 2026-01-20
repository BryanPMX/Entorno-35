import { MarketingNavigation } from "@/components/marketing/navigation";
import { HeroSection } from "@/components/marketing/hero-section";
import { FeaturesSection } from "@/components/marketing/features-section";
import { ComplianceSection } from "@/components/marketing/compliance-section";
import { PricingSection } from "@/components/marketing/pricing-section";
import { MarketingFooter } from "@/components/marketing/footer";

/**
 * Root Page - Marketing Landing Page
 *
 * Professional landing page for Entorno-35 NOM-035 compliance platform.
 * Publicly accessible page designed to attract companies and convert them into customers.
 */
export default function HomePage() {
  return (
    <div className="min-h-screen">
      <MarketingNavigation />
      <main>
        <HeroSection />
        <FeaturesSection />
        <ComplianceSection />
        <PricingSection />
      </main>
      <MarketingFooter />
    </div>
  );
}
