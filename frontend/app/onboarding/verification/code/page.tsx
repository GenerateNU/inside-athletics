"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOnboarding } from "@/utils/onboarding";

export default function OnboardingVerificationCodePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [code, setCode] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setCode(data.verification.code);
  }, [data.verification.code, hydrated]);

  const canContinue = Boolean(code.trim());

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Verification Code</h1>
          <p className="text-sm text-gray-600">
            Enter the verification code sent to your email.
          </p>
        </div>

        <div className="space-y-3">
          <label htmlFor="verification-code" className="block text-sm font-medium text-black">
            Code
          </label>
          <Input
            id="verification-code"
            type="text"
            value={code}
            placeholder="Enter your code"
            className="h-10 rounded-xl px-3 text-sm"
            onChange={(event) => {
              setCode(event.target.value);
            }}
          />
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("verification", {
              code,
            });
            router.push(
              role === "athlete"
                ? `/onboarding/teams-of-interest?role=${encodeURIComponent(role)}`
                : "/onboarding/plan",
            );
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
