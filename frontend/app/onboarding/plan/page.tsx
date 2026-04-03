"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";

const planOptions = [
  {
    label: "Free Plan",
    value: "free",
    description: "Get started with the core Inside Athletics experience.",
  },
  {
    label: "Premium Plan",
    value: "premium",
    description: "Unlock expanded access and premium recruiting tools.",
  },
] as const;

export default function OnboardingPlanPage() {
  const router = useRouter();
  const [selectedPlan, setSelectedPlan] = useState("");

  const canContinue = Boolean(selectedPlan);

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Choose Your Plan</h1>
          <p className="text-sm text-gray-600">
            Pick the plan you want to start with.
          </p>
        </div>

        <div className="grid grid-cols-2 gap-3">
          {planOptions.map((plan) => {
            const isSelected = selectedPlan === plan.value;

            return (
              <Button
                key={plan.value}
                type="button"
                variant="outline"
                className="flex h-full min-h-36 w-full min-w-0 flex-col items-start gap-2 rounded-xl px-4 py-4 text-left whitespace-normal"
                style={{
                  borderColor: "#16A34A",
                  backgroundColor: isSelected ? "#16A34A" : "#FFFFFF",
                  color: isSelected ? "#FFFFFF" : "#000000",
                }}
                onClick={() => {
                  setSelectedPlan(plan.value);
                }}
              >
                <span className="w-full text-sm font-semibold break-words">
                  {plan.label}
                </span>
                <span
                  className="w-full text-sm break-words"
                  style={{ color: isSelected ? "#F0FDF4" : "#4B5563" }}
                >
                  {plan.description}
                </span>
              </Button>
            );
          })}
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            router.push("/");
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
