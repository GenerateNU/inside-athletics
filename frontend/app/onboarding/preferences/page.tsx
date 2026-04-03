"use client";

import { useMemo, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { XIcon } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const divisions = [
  { label: "Division I", value: "division-i" },
  { label: "Division II", value: "division-ii" },
  { label: "Division III", value: "division-iii" },
] as const;

const associations = [
  { label: "NCAA", value: "ncaa" },
  { label: "NJCAA", value: "njcaa" },
] as const;

const universityOptions = [
  "Boston College",
  "Boston University",
  "Duke University",
  "Harvard University",
  "Northeastern University",
  "Ohio State University",
  "Penn State University",
  "Stanford University",
  "Syracuse University",
  "Texas A&M University",
  "University of Connecticut",
  "University of Florida",
  "University of Michigan",
  "University of North Carolina",
  "University of Notre Dame",
  "University of Southern California",
] as const;

export default function OnboardingPreferencesPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const role = searchParams.get("role") ?? "";
  const [division, setDivision] = useState("");
  const [association, setAssociation] = useState("");
  const [search, setSearch] = useState("");
  const [selectedUniversities, setSelectedUniversities] = useState<string[]>(
    [],
  );

  const filteredUniversities = useMemo(() => {
    const query = search.trim().toLowerCase();

    return universityOptions.filter((school) => {
      const matchesQuery =
        query.length === 0 || school.toLowerCase().includes(query);

      return matchesQuery && !selectedUniversities.includes(school);
    });
  }, [search, selectedUniversities]);

  const canContinue = Boolean(
    division && association && selectedUniversities.length > 0,
  );

  const handleDivisionChange = (value: string | null) => {
    setDivision(value ?? "");
  };

  const handleAssociationChange = (value: string | null) => {
    setAssociation(value ?? "");
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Refine Your Search</h1>
          <p className="text-sm text-gray-600">
            Choose your level and add universities you want to keep an eye on.
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
            {filteredUniversities.length > 0 ? (
              filteredUniversities.map((school) => (
                <button
                  key={school}
                  type="button"
                  className="w-full rounded-lg px-3 py-2 text-left text-sm text-black transition-colors hover:bg-green-50"
                  onClick={() => {
                    setSelectedUniversities((current) => [...current, school]);
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

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            router.push(
              `/onboarding/verification?role=${encodeURIComponent(role)}`,
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
