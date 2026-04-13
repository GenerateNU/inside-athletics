"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { XIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useSession } from "@/utils/SessionContext";
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

type CollegeResponse = {
  id: string;
  name: string;
};

type CollegeListPayload = {
  colleges: CollegeResponse[];
  total: number;
};

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

export default function OnboardingTeamsOfInterestPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const session = useSession();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const [division, setDivision] = useState("");
  const [association, setAssociation] = useState("");
  const [primarySport, setPrimarySport] = useState("");
  const [collegeSearch, setCollegeSearch] = useState("");
  const [selectedUniversities, setSelectedUniversities] = useState<string[]>(
    [],
  );
  const [collegeOptions, setCollegeOptions] = useState<string[]>([]);
  const [isLoadingColleges, setIsLoadingColleges] = useState(false);
  const [collegeError, setCollegeError] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setDivision(data.preferences.division);
    setAssociation(data.preferences.association);
    setPrimarySport(data.preferences.primarySport);
    setSelectedUniversities(data.preferences.selectedUniversities);
  }, [data.preferences, hydrated]);

  useEffect(() => {
    if (!session?.access_token) {
      setIsLoadingColleges(false);
      setCollegeError("You need an active session to load schools.");
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
          const errorText = await response.text();
          const trimmedErrorText = errorText.trim();
          throw new Error(
            trimmedErrorText
              ? `Unable to load schools (${response.status}): ${trimmedErrorText}`
              : `Unable to load schools (${response.status} ${response.statusText}).`,
          );
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
          error instanceof Error ? error.message : "Unable to load schools.",
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

  const filteredCollegeOptions = collegeOptions
    .filter((school) => {
      if (!collegeSearch.trim()) {
        return false;
      }

      return (
        school.toLowerCase().includes(collegeSearch.trim().toLowerCase()) &&
        !selectedUniversities.includes(school)
      );
    })
    .slice(0, 8);

  const selectedPrimarySportLabel =
    primarySports.find((sport) => sport.value === primarySport)?.label ?? "";

  const canContinue = Boolean(
    division && association && primarySport && selectedUniversities.length > 0,
  );

  const addUniversity = (school: string) => {
    if (!school || selectedUniversities.includes(school)) {
      return;
    }

    setSelectedUniversities((current) => [...current, school]);
    setCollegeSearch("");
  };

  const removeUniversity = (school: string) => {
    setSelectedUniversities((current) =>
      current.filter((currentSchool) => currentSchool !== school),
    );
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">Teams of Interest</h1>
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
                  className={`h-12 rounded-xl text-sm font-semibold text-black ${
                    isSelected
                      ? "border-[#7F8C2D] bg-[#D4E94B80]"
                      : "border-[#D4E94B] bg-[#FCFDF1]"
                  }`}
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
          <label
            htmlFor="primary-sport"
            className="block text-sm font-medium text-black"
          >
            Primary Sport
          </label>
          <Select
            value={primarySport}
            onValueChange={(value) => {
              setPrimarySport(value ?? "");
            }}
          >
            <SelectTrigger
              id="primary-sport"
              className="h-10 w-full border-[#3E7DBB] text-sm"
            >
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
          <label
            htmlFor="schools-of-interest"
            className="block text-sm font-medium text-black"
          >
            Schools of Interest
          </label>
          <div className="relative">
            <Input
              id="schools-of-interest"
              value={collegeSearch}
              onChange={(event) => {
                setCollegeSearch(event.target.value);
              }}
              placeholder={
                isLoadingColleges ? "Loading schools..." : "Search for a school"
              }
              className="h-10 border-[#3E7DBB] px-3 text-sm text-black placeholder:text-gray-500"
              autoComplete="off"
            />
            {filteredCollegeOptions.length > 0 ? (
              <div className="absolute left-0 right-0 top-full z-10 mt-2 max-h-56 overflow-y-auto rounded-xl border border-[#D4E94B] bg-white shadow-sm">
                {filteredCollegeOptions.map((school) => (
                  <button
                    key={school}
                    type="button"
                    className="w-full px-3 py-2 text-left text-sm text-black transition-colors hover:bg-[#FCFDF1]"
                    onClick={() => {
                      addUniversity(school);
                    }}
                  >
                    {school}
                  </button>
                ))}
              </div>
            ) : null}
          </div>
          {selectedUniversities.length > 0 ? (
            <div className="flex flex-wrap gap-2">
              {selectedUniversities.map((school) => (
                <Badge
                  key={school}
                  className="h-auto rounded-full border border-[#7F8C2D] bg-[#D4E94B80] px-3 py-1 text-xs font-medium text-black"
                >
                  <button
                    type="button"
                    className="mr-1 inline-flex items-center justify-center text-black"
                    onClick={() => {
                      removeUniversity(school);
                    }}
                    aria-label={`Remove ${school}`}
                  >
                    <XIcon className="size-3" />
                  </button>
                  <span>{school}</span>
                </Badge>
              ))}
            </div>
          ) : null}
          {collegeError ? (
            <p className="text-sm text-red-600" role="alert">
              {collegeError}
            </p>
          ) : null}
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
                  className={`h-12 rounded-xl text-sm font-semibold text-black ${
                    isSelected
                      ? "border-[#7F8C2D] bg-[#D4E94B80]"
                      : "border-[#D4E94B] bg-[#FCFDF1]"
                  }`}
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
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={() => {
            updateSection("preferences", {
              division,
              association,
              primarySport,
              selectedUniversities,
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
