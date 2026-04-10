"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { submitOnboardingUser } from "@/utils/onboarding-submit";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";

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
  const session = useSession();
  const { data, hydrated, reset, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [selectedPlan, setSelectedPlan] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

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
                className="rounded-xl p-[2px] transition-all"
                style={{
                  background: isPremium
                    ? "#7F8C2D"
                    : "#D4E94B",
                  boxShadow: isSelected
                    ? isPremium
                      ? "4px 6px 14px rgba(127, 140, 45, 0.24)"
                      : "4px 6px 14px rgba(44, 100, 154, 0.18)"
                    : "none",
                  transform: isSelected ? "translateY(-2px)" : "none",
                }}
              >
                <Button
                  type="button"
                  variant="outline"
                  className="flex h-full min-h-72 w-full min-w-0 flex-col items-start rounded-[calc(0.75rem-2px)] px-0 py-0 text-left whitespace-normal"
                  style={{
                    borderColor: "transparent",
                    background: isPremium
                      ? "#E9F4A5"
                      : "#FCFDF1",
                    color: "#000000",
                  }}
                  onClick={() => {
                    setSelectedPlan(plan.value);
                  }}
                >
                  <div className="w-full px-4 py-4">
                    {isPremium ? (
                      <div
                        className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border-[1.5px] px-3 py-3 text-sm font-semibold"
                        style={{
                          backgroundColor: "#FCFDF1",
                          borderColor: "#7F8C2D",
                        }}
                      >
                        {plan.label}
                      </div>
                    ) : (
                      <div
                        className="mx-3 flex min-h-16 w-[calc(100%-1.5rem)] items-start rounded-md border px-3 py-3 text-sm font-semibold"
                        style={{
                          backgroundColor: "#FFFFFF",
                          borderColor: "#D4E94B",
                        }}
                      >
                        {plan.label}
                      </div>
                    )}
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "#7F8C2D"
                          : "#D4E94B",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#000000" : "#4B5563",
                    }}
                  >
                    {plan.price}
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "#7F8C2D"
                          : "#D4E94B",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#000000" : "#4B5563",
                    }}
                  >
                    Feature 1
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "#7F8C2D"
                          : "#D4E94B",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#000000" : "#4B5563",
                    }}
                  >
                    Feature 2
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "#7F8C2D"
                          : "#D4E94B",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#000000" : "#4B5563",
                    }}
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
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={async () => {
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

            if (!session?.access_token) {
              setError(
                "You need an active session before finishing onboarding.",
              );
              return;
            }

            setIsSubmitting(true);
            setError("");

            try {
              await submitOnboardingUser(
                {
                  ...data,
                  plan: {
                    selectedPlan,
                  },
                },
                session.access_token,
                session.user.email,
              );
              reset();
              router.push("/");
            } catch (submissionError) {
              setError(
                submissionError instanceof Error
                  ? submissionError.message
                  : "Unable to finish onboarding.",
              );
            } finally {
              setIsSubmitting(false);
            }
          }}
          disabled={!canContinue || isSubmitting}
        >
          Continue
        </Button>
        {error ? (
          <p className="text-sm text-red-600" role="alert">
            {error}
          </p>
        ) : null}
      </div>
    </div>
  );
}
