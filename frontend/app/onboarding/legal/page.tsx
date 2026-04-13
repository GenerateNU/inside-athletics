"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";

export default function OnboardingLegalPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [accepted, setAccepted] = useState(false);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setAccepted(data.legal.accepted);
  }, [data.legal.accepted, hydrated]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-2xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">
            Privacy Policy & Terms
          </h1>
          <p className="text-sm text-gray-600">
            Review this note before continuing to verification.
          </p>
        </div>

        <div className="space-y-6 rounded-xl border border-gray-200 bg-gray-50 p-6">
          <section className="space-y-2">
            <h2 className="text-lg font-semibold text-[#001F3E]">Privacy Policy</h2>
            <p className="text-sm leading-6 text-gray-700">
              Privacy note: Content shared on Inside Athletics may be used in
              aggregated and anonymized form to generate insights, reports, and
              educational or social media content about college athletics. We do
              not publicly share individual posts in a way that can identify a
              user.
            </p>
          </section>

          <section className="space-y-2">
            <h2 className="text-lg font-semibold text-[#001F3E]">
              Terms & Conditions
            </h2>
            <p className="text-sm leading-6 text-gray-700">
              By continuing, you acknowledge this privacy note and agree to
              Inside Athletics&apos; terms and conditions for platform use.
            </p>
          </section>
        </div>

        <label className="flex items-start gap-3 rounded-xl border border-gray-200 p-4 text-sm text-gray-700">
          <input
            type="checkbox"
            checked={accepted}
            className="mt-1 size-4 rounded border-gray-300"
            onChange={(event) => {
              setAccepted(event.target.checked);
            }}
          />
          <span>
            I have read the privacy policy and terms and conditions, and I agree
            to continue.
          </span>
        </label>

        <div className="flex gap-3">
          <Button
            type="button"
            variant="outline"
            className="h-10 flex-1 rounded-xl text-sm font-semibold"
            onClick={() => {
              router.back();
            }}
          >
            Back
          </Button>
          <Button
            type="button"
            className="h-10 flex-1 rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
            onClick={() => {
              updateSection("legal", {
                accepted,
              });
              router.push(
                `/onboarding/verification?role=${encodeURIComponent(role)}`,
              );
            }}
            disabled={!accepted}
          >
            Continue
          </Button>
        </div>
      </div>
    </div>
  );
}
