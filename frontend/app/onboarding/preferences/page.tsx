"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

type CollegeResponse = {
  id: string;
  name: string;
};

type CollegeListPayload = {
  colleges: CollegeResponse[];
  total: number;
};

const primarySports = [
  { label: "Basketball", value: "basketball" },
  { label: "Soccer", value: "soccer" },
  { label: "Track & Field", value: "track-and-field" },
  { label: "Volleyball", value: "volleyball" },
  { label: "Tennis", value: "tennis" },
  { label: "Swimming", value: "swimming" },
  { label: "Softball", value: "softball" },
  { label: "Baseball", value: "baseball" },
] as const;

const programOptions = [
  { label: "Women's", value: "womens" },
  { label: "Men's", value: "mens" },
] as const;

function unwrapBody<T>(value: unknown): T | undefined {
  let current = value;

  for (let depth = 0; depth < 5; depth += 1) {
    if (!current || typeof current !== "object") {
      return current as T | undefined;
    }

    if ("body" in current && current.body !== undefined) {
      current = current.body;
      continue;
    }

    if ("Body" in current && current.Body !== undefined) {
      current = current.Body;
      continue;
    }

    return current as T | undefined;
  }

  return current as T | undefined;
}

export default function OnboardingPreferencesPage() {
  const router = useRouter();
  const session = useSession();
  const { data, hydrated, updateSection } = useOnboarding();
  const [primarySport, setPrimarySport] = useState("");
  const [program, setProgram] = useState("");
  const [university, setUniversity] = useState("");
  const [collegeOptions, setCollegeOptions] = useState<string[]>([]);
  const [isLoadingColleges, setIsLoadingColleges] = useState(true);
  const [collegeError, setCollegeError] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setPrimarySport(data.preferences.primarySport);
    setProgram(data.preferences.program);
    setUniversity(data.preferences.university);
  }, [data.preferences, hydrated]);

  useEffect(() => {
    if (!session?.access_token) {
      setIsLoadingColleges(false);
      setCollegeOptions([]);
      setCollegeError("You need an active session to load colleges.");
      return;
    }

    let cancelled = false;

    async function loadColleges() {
      try {
        setCollegeError("");
        setIsLoadingColleges(true);

        const response = await fetch("/api/v1/college/?limit=500", {
          headers: {
            Authorization: `Bearer ${session.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error("Unable to load colleges.");
        }

        const payload = unwrapBody<CollegeListPayload>(await response.json());
        const colleges = (payload?.colleges ?? [])
          .map((college) => college.name.trim())
          .filter(Boolean);

        if (cancelled) {
          return;
        }

        setCollegeOptions(colleges);
      } catch (error) {
        if (cancelled) {
          return;
        }

        setCollegeError(
          error instanceof Error ? error.message : "Unable to load colleges.",
        );
      } finally {
        if (!cancelled) {
          setIsLoadingColleges(false);
        }
      }
    }

    loadColleges();

    return () => {
      cancelled = true;
    };
  }, [session?.access_token]);

  const canContinue = Boolean(primarySport && program && university);

  const handlePrimarySportChange = (value: string | null) => {
    setPrimarySport(value ?? "");
  };

  const selectedPrimarySportLabel =
    primarySports.find((sport) => sport.value === primarySport)?.label ?? "";

  const handleUniversityChange = (value: string | null) => {
    setUniversity(value ?? "");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">About you</h1>
          <p className="text-sm text-gray-600">Tell us about yourself!</p>
        </div>

        <div className="space-y-3">
          <label
            htmlFor="primary-sport"
            className="block text-sm font-medium text-black"
          >
            Primary Sport
          </label>
          <Select value={primarySport} onValueChange={handlePrimarySportChange}>
            <SelectTrigger id="primary-sport" className="h-10 w-full text-sm">
              <SelectValue placeholder="Select a primary sport">
                {selectedPrimarySportLabel}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              {primarySports.map((sport) => (
                <SelectItem key={sport.value} value={sport.value}>
                  {sport.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-3">
          <p className="block text-sm font-medium text-black">
            Which program would you join
          </p>
          <div className="grid grid-cols-2 gap-3">
            {programOptions.map((option) => {
              const isSelected = program === option.value;

              return (
                <Button
                  key={option.value}
                  type="button"
                  variant="outline"
                  className="h-12 rounded-xl text-sm font-semibold"
                  style={{
                    borderColor: "#16A34A",
                    backgroundColor: isSelected ? "#16A34A" : "#FFFFFF",
                    color: isSelected ? "#FFFFFF" : "#000000",
                  }}
                  onClick={() => {
                    setProgram(option.value);
                  }}
                >
                  {option.label}
                </Button>
              );
            })}
          </div>
        </div>

        <div className="space-y-3">
          <label
            htmlFor="university"
            className="block text-sm font-medium text-black"
          >
            University
          </label>
          <Select value={university} onValueChange={handleUniversityChange}>
            <SelectTrigger id="university" className="h-10 w-full text-sm">
              <SelectValue
                placeholder={
                  isLoadingColleges
                    ? "Loading universities..."
                    : "Select a university"
                }
              />
            </SelectTrigger>
            <SelectContent>
              {collegeOptions.map((school) => (
                <SelectItem key={school} value={school}>
                  {school}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          {collegeError ? (
            <p className="text-sm text-red-600" role="alert">
              {collegeError}
            </p>
          ) : null}
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("preferences", {
              primarySport,
              program,
              university,
            });
            router.push(
              `/onboarding/legal?role=${encodeURIComponent(data.role.role)}`,
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
