import { Suspense } from "react";
import OnboardingBillingPage from "./OnboardingBillingPage";

export const dynamic = "force-dynamic";

export default function BillingPage() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <OnboardingBillingPage />
    </Suspense>
  );
}