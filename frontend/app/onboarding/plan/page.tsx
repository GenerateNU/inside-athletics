"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";

const planOptions = [
  {
    label: "Premium Plan",
    value: "premium",
    price: "$9.99/mo",
  },
  {
    label: "Standard Plan",
    value: "free",
    price: "$0/mo",
  },
] as const;

export default function OnboardingPlanPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [selectedPlan, setSelectedPlan] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setSelectedPlan(data.plan.selectedPlan);
  }, [data.plan.selectedPlan, hydrated]);

  const canContinue = Boolean(selectedPlan);

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-3xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Choose Your Plan</h1>
          <p className="text-sm text-gray-600">
            Pick the plan you want to start with.
          </p>
        </div>

        <div className="grid grid-cols-2 gap-5">
          {planOptions.map((plan) => {
            const isSelected = selectedPlan === plan.value;
            const isPremium = plan.value === "premium";

            return (
              <div
                key={plan.value}
                className={`rounded-xl p-[2px] transition-all ${
                  isPremium ? "bg-[#7F8C2D]" : "bg-[#D4E94B]"
                } ${
                  isSelected
                    ? isPremium
                      ? "-translate-y-0.5 shadow-[4px_6px_14px_rgba(127,140,45,0.24)]"
                      : "-translate-y-0.5 shadow-[4px_6px_14px_rgba(44,100,154,0.18)]"
                    : ""
                }`}
              >
                <Button
                  type="button"
                  variant="outline"
                  className={`flex h-full min-h-72 w-full min-w-0 flex-col items-start rounded-[calc(0.75rem-2px)] border-transparent px-0 py-0 text-left text-black whitespace-normal ${
                    isPremium ? "bg-[#E9F4A5]" : "bg-[#FCFDF1]"
                  }`}
                  onClick={() => {
                    setSelectedPlan(plan.value);
                  }}
                >
                  <div className="w-full px-4 py-4">
                    {isPremium ? (
                      <div
                        className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border-[1.5px] border-[#7F8C2D] bg-[#FCFDF1] px-3 py-3 text-sm font-semibold"
                      >
                        {plan.label}
                      </div>
                    ) : (
                      <div
                        className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border border-[#D4E94B] bg-white px-3 py-3 text-sm font-semibold"
                      >
                        {plan.label}
                      </div>
                    )}
                  </div>
                  <div className="w-full px-4">
                    <div
                      className={`mx-3 border-t ${
                        isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"
                      }`}
                    />
                  </div>
                  <div
                    className={`w-full px-4 py-4 text-sm ${
                      isPremium ? "text-black" : "text-gray-600"
                    }`}
                  >
                    {plan.price}
                  </div>
                  <div className="w-full px-4">
                    <div
                      className={`mx-3 border-t ${
                        isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"
                      }`}
                    />
                  </div>
                  <div
                    className={`w-full px-4 py-4 text-sm ${
                      isPremium ? "text-black" : "text-gray-600"
                    }`}
                  >
                    Feature 1
                  </div>
                  <div className="w-full px-4">
                    <div
                      className={`mx-3 border-t ${
                        isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"
                      }`}
                    />
                  </div>
                  <div
                    className={`w-full px-4 py-4 text-sm ${
                      isPremium ? "text-black" : "text-gray-600"
                    }`}
                  >
                    Feature 2
                  </div>
                  <div className="w-full px-4">
                    <div
                      className={`mx-3 border-t ${
                        isPremium ? "border-[#7F8C2D]" : "border-[#D4E94B]"
                      }`}
                    />
                  </div>
                  <div
                    className={`w-full px-4 py-4 text-sm ${
                      isPremium ? "text-black" : "text-gray-600"
                    }`}
                  >
                    Feature 3
                  </div>
                </Button>
              </div>
            );
          })}
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={() => {
            updateSection("plan", {
              selectedPlan,
            });

            if (selectedPlan === "free") {
              router.push(
                role
                  ? `/onboarding/topic-tags?role=${encodeURIComponent(role)}`
                  : "/onboarding/topic-tags",
              );
              return;
            }

            router.push(
              role
                ? `/onboarding/billing?role=${encodeURIComponent(role)}`
                : "/onboarding/billing",
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
