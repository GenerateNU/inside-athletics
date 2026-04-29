import { Suspense } from "react";
import OnboardingBillingPage from "./OnboardingBillingPage";
import { Spinner } from "@/components/ui/spinner";

export const dynamic = "force-dynamic";

export default function BillingPage() {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-screen"><Spinner className="size-6" /></div>}>
      <OnboardingBillingPage />
    </Suspense>
  );
}