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
                    ? "linear-gradient(180deg, #2C649A 0%, #5F95C7 100%)"
                    : "linear-gradient(180deg, #2C649A 0%, #16A34A 100%)",
                  boxShadow: isSelected
                    ? isPremium
                      ? "0 0 0 3px rgba(44, 100, 154, 0.28), 0 14px 30px rgba(44, 100, 154, 0.22)"
                      : "0 0 0 3px rgba(22, 163, 74, 0.22), 0 14px 30px rgba(44, 100, 154, 0.12)"
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
                      ? "linear-gradient(180deg, #2C649A 0%, #76A7D4 100%)"
                      : "#FFFFFF",
                    color: isPremium ? "#FFFFFF" : "#000000",
                  }}
                  onClick={() => {
                    setSelectedPlan(plan.value);
                  }}
                >
                  <div className="w-full px-4 py-4">
                    <span className="block w-full text-sm font-semibold break-words">
                      {plan.label}
                    </span>
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "rgba(234, 244, 255, 0.35)"
                          : "rgba(44, 100, 154, 0.2)",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#EAF4FF" : "#4B5563",
                    }}
                  >
                    {plan.price}
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "rgba(234, 244, 255, 0.35)"
                          : "rgba(44, 100, 154, 0.2)",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#EAF4FF" : "#4B5563",
                    }}
                  >
                    Feature 1
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "rgba(234, 244, 255, 0.35)"
                          : "rgba(44, 100, 154, 0.2)",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#EAF4FF" : "#4B5563",
                    }}
                  >
                    Feature 2
                  </div>
                  <div className="w-full px-4">
                    <div
                      className="mx-3 border-t"
                      style={{
                        borderColor: isPremium
                          ? "rgba(234, 244, 255, 0.35)"
                          : "rgba(44, 100, 154, 0.2)",
                      }}
                    />
                  </div>
                  <div
                    className="w-full px-4 py-4 text-sm"
                    style={{
                      color: isPremium ? "#EAF4FF" : "#4B5563",
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
