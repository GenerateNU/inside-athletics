"use client";

import { useEffect, useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { XIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
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

const divisions = [
  { label: "Division I", value: "division-i" },
  { label: "Division II", value: "division-ii" },
  { label: "Division III", value: "division-iii" },
] as const;

const associations = [
  { label: "NCAA", value: "ncaa" },
  { label: "NJCAA", value: "njcaa" },
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
  const searchParams = useSearchParams();
  const session = useSession();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const isAthlete = role === "athlete" || role === "prospective-athlete";
  const [division, setDivision] = useState("");
  const [association, setAssociation] = useState("");
  const [search, setSearch] = useState("");
  const [selectedUniversities, setSelectedUniversities] = useState<string[]>(
    [],
  );
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

    setDivision(data.preferences.division);
    setAssociation(data.preferences.association);
    setSelectedUniversities(data.preferences.selectedUniversities);
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

  const filteredUniversities = useMemo(() => {
    const query = search.trim().toLowerCase();

    return collegeOptions.filter((school) => {
      const matchesQuery =
        query.length === 0 || school.toLowerCase().includes(query);

      return matchesQuery && !selectedUniversities.includes(school);
    });
  }, [collegeOptions, search, selectedUniversities]);

  const canContinue = isAthlete
    ? Boolean(primarySport && program && university)
    : Boolean(division && association && selectedUniversities.length > 0);

  const handleDivisionChange = (value: string | null) => {
    setDivision(value ?? "");
  };

  const handleAssociationChange = (value: string | null) => {
    setAssociation(value ?? "");
  };

  const handlePrimarySportChange = (value: string | null) => {
    setPrimarySport(value ?? "");
  };

  const handleUniversityChange = (value: string | null) => {
    setUniversity(value ?? "");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">
            {isAthlete ? "About you" : "Refine Your Search"}
          </h1>
          <p className="text-sm text-gray-600">
            {isAthlete
              ? "Tell us about yourself!"
              : "Choose your level and add universities you want to keep an eye on."}
          </p>
        </div>

        {isAthlete ? (
          <>
            <div className="space-y-3">
              <label
                htmlFor="primary-sport"
                className="block text-sm font-medium text-black"
              >
                Primary Sport
              </label>
              <Select
                value={primarySport}
                onValueChange={handlePrimarySportChange}
              >
                <SelectTrigger
                  id="primary-sport"
                  className="h-10 w-full text-sm"
                >
                  <SelectValue placeholder="Select a primary sport" />
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
          </>
        ) : (
          <>
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
                        handleDivisionChange(item.value);
                      }}
                    >
                      {item.label}
                    </Button>
                  );
                })}
              </div>
            </div>

            <div className="space-y-3">
              <p className="block text-sm font-medium text-black">
                Association
              </p>
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
                        handleAssociationChange(item.value);
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
                htmlFor="university-search"
                className="block text-sm font-medium text-black"
              >
                Universities of Interest
              </label>
              <Input
                id="university-search"
                type="text"
                value={search}
                placeholder="Search universities"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setSearch(event.target.value);
                }}
              />

              <div className="max-h-56 space-y-2 overflow-y-auto rounded-xl border border-gray-200 bg-white p-2">
                {isLoadingColleges ? (
                  <p className="px-3 py-2 text-sm text-gray-500">
                    Loading universities...
                  </p>
                ) : filteredUniversities.length > 0 ? (
                  filteredUniversities.map((school) => (
                    <button
                      key={school}
                      type="button"
                      className="w-full rounded-lg px-3 py-2 text-left text-sm text-black transition-colors hover:bg-green-50"
                      onClick={() => {
                        setSelectedUniversities((current) => [
                          ...current,
                          school,
                        ]);
                        setSearch("");
                      }}
                    >
                      {school}
                    </button>
                  ))
                ) : (
                  <p className="px-3 py-2 text-sm text-gray-500">
                    No universities match that search.
                  </p>
                )}
              </div>
              {collegeError ? (
                <p className="text-sm text-red-600" role="alert">
                  {collegeError}
                </p>
              ) : null}

              <div className="flex flex-wrap gap-2">
                {selectedUniversities.map((school) => (
                  <Badge
                    key={school}
                    variant="outline"
                    className="h-auto gap-2 rounded-full border-green-600 px-3 py-1 text-xs text-black"
                  >
                    <span>{school}</span>
                    <button
                      type="button"
                      aria-label={`Remove ${school}`}
                      className="flex items-center justify-center text-black"
                      onClick={() => {
                        setSelectedUniversities((current) =>
                          current.filter((item) => item !== school),
                        );
                      }}
                    >
                      <XIcon className="size-3" />
                    </button>
                  </Badge>
                ))}
              </div>
            </div>
          </>
        )}

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("preferences", {
              division,
              association,
              selectedUniversities,
              primarySport,
              program,
              university,
            });
            router.push(`/onboarding/legal?role=${encodeURIComponent(role)}`);
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
