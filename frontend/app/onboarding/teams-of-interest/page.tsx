"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";

const divisions = [
  { label: "Division I", value: "division-i" },
  { label: "Division II", value: "division-ii" },
  { label: "Division III", value: "division-iii" },
] as const;

const associations = [
  { label: "NCAA", value: "ncaa" },
  { label: "NJCAA", value: "njcaa" },
] as const;

export default function OnboardingTeamsOfInterestPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [division, setDivision] = useState("");
  const [association, setAssociation] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setDivision(data.preferences.division);
    setAssociation(data.preferences.association);
  }, [data.preferences.association, data.preferences.division, hydrated]);

  const canContinue = Boolean(division && association);

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Teams of Interest</h1>
          <p className="text-sm text-gray-600">
            Tell us what you are looking to pursue
          </p>
        </div>

        <div className="space-y-3">
          <p className="block text-sm font-medium text-black">Division</p>
          <div className="grid grid-cols-3 gap-3">
            {divisions.map((item) => {
              const isSelected = division === item.value;

              return (
                <Button
                  key={item.value}
                  type="button"
                  variant="outline"
                  className="h-12 rounded-xl text-sm font-semibold"
                  style={{
                    borderColor: "#16A34A",
                    backgroundColor: isSelected ? "#16A34A" : "#FFFFFF",
                    color: isSelected ? "#FFFFFF" : "#000000",
                  }}
                  onClick={() => {
                    setDivision(item.value);
                  }}
                >
                  {item.label}
                </Button>
              );
            })}
          </div>
        </div>

        <div className="space-y-3">
          <p className="block text-sm font-medium text-black">Association</p>
          <div className="grid grid-cols-2 gap-3">
            {associations.map((item) => {
              const isSelected = association === item.value;

              return (
                <Button
                  key={item.value}
                  type="button"
                  variant="outline"
                  className="h-12 rounded-xl text-sm font-semibold"
                  style={{
                    borderColor: "#16A34A",
                    backgroundColor: isSelected ? "#16A34A" : "#FFFFFF",
                    color: isSelected ? "#FFFFFF" : "#000000",
                  }}
                  onClick={() => {
                    setAssociation(item.value);
                  }}
                >
                  {item.label}
                </Button>
              );
            })}
          </div>
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("preferences", {
              division,
              association,
            });
            router.push(
              role
                ? `/onboarding/plan?role=${encodeURIComponent(role)}`
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
