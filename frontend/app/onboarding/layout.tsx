import { Suspense } from "react";
import { Spinner } from "@/components/ui/spinner";

export const dynamic = "force-dynamic";

export default function OnboardingLayout({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<div className="flex items-center justify-center min-h-screen"><Spinner className="size-6" /></div>}>
      {children}
    </Suspense>
  );
}