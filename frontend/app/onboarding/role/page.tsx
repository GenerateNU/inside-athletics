"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

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

const universities = [
  { label: "Northeastern University", value: "northeastern" },
  { label: "Boston College", value: "boston-college" },
  { label: "Boston University", value: "boston-university" },
  { label: "Harvard University", value: "harvard" },
  { label: "University of Connecticut", value: "uconn" },
  { label: "University of Michigan", value: "michigan" },
  { label: "University of Notre Dame", value: "notre-dame" },
  { label: "Stanford University", value: "stanford" },
] as const;

export default function OnboardingRolePage() {
  const router = useRouter();
  const [role, setRole] = useState("");
  const [primarySport, setPrimarySport] = useState("");
  const [program, setProgram] = useState("");
  const [university, setUniversity] = useState("");

  const canContinue = Boolean(role && primarySport && program && university);

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Choose Your Role</h1>
          <p className="text-sm text-gray-600">
            Select the role that best describes how you&apos;ll use Inside
            Athletics.
          </p>
        </div>

        <div className="space-y-3">
          <label htmlFor="role" className="block text-sm font-medium text-black">
            Role
          </label>
          <Select value={role} onValueChange={setRole}>
            <SelectTrigger id="role" className="h-10 w-full text-sm">
              <SelectValue placeholder="Select a role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="athlete">Athlete</SelectItem>
              <SelectItem value="coach">Coach</SelectItem>
              <SelectItem value="parent">Parent</SelectItem>
              <SelectItem value="recruiter">Recruiter</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-3">
          <label
            htmlFor="primary-sport"
            className="block text-sm font-medium text-black"
          >
            Primary Sport
          </label>
          <Select value={primarySport} onValueChange={setPrimarySport}>
            <SelectTrigger id="primary-sport" className="h-10 w-full text-sm">
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
          <Select value={university} onValueChange={setUniversity}>
            <SelectTrigger id="university" className="h-10 w-full text-sm">
              <SelectValue placeholder="Select a university" />
            </SelectTrigger>
            <SelectContent>
              {universities.map((school) => (
                <SelectItem key={school.value} value={school.value}>
                  {school.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
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
