"use client";

import { useState } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const genderOptions = [
  "Women's",
  "Men's",
  "Co-ed",
] as const;

const sportOptions = [
  "Basketball",
  "Soccer",
  "Track & Field",
  "Volleyball",
  "Tennis",
] as const;

const ratingRows = [
  { label: "Program Strength", value: 92 },
  { label: "Coaching", value: 88 },
  { label: "Athlete Development", value: 84 },
];

export function RatingPanel() {
  const [programGender, setProgramGender] = useState<string>("");
  const [sportProgram, setSportProgram] = useState<string>("");

  const handleProgramGenderChange = (value: string | null) => {
    setProgramGender(value ?? "");
  };

  const handleSportProgramChange = (value: string | null) => {
    setSportProgram(value ?? "");
  };

  return (
    <Card className="w-full max-w-3xl border border-black/10 bg-white py-0 shadow-[0_20px_60px_rgba(0,0,0,0.06)]">
      <CardHeader className="border-b border-black/8 px-6 py-5">
        <CardTitle className="text-left text-2xl font-semibold text-black">
          Rating
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-8 px-6 py-6">
        <div className="grid gap-4 md:grid-cols-2">
          <Select value={programGender} onValueChange={handleProgramGenderChange}>
            <SelectTrigger className="h-12 w-full rounded-xl border-black/10 bg-zinc-50 px-4 text-sm text-black">
              <SelectValue placeholder="Select Program Gender" />
            </SelectTrigger>
            <SelectContent>
              {genderOptions.map((option) => (
                <SelectItem key={option} value={option}>
                  {option}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <Select value={sportProgram} onValueChange={handleSportProgramChange}>
            <SelectTrigger className="h-12 w-full rounded-xl border-black/10 bg-zinc-50 px-4 text-sm text-black">
              <SelectValue placeholder="Select Sport Program" />
            </SelectTrigger>
            <SelectContent>
              {sportOptions.map((option) => (
                <SelectItem key={option} value={option}>
                  {option}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        <div className="grid gap-6 lg:grid-cols-[1.05fr_0.95fr]">
          <div className="rounded-2xl bg-zinc-950 px-6 py-7 text-white">
            <p className="text-sm font-medium uppercase tracking-[0.18em] text-white/70">
              Overall Score
            </p>
            <div className="mt-4 flex items-end gap-3">
              <span className="text-6xl font-semibold leading-none">8.7</span>
              <span className="pb-1 text-sm text-white/70">out of 10</span>
            </div>
            <p className="mt-4 max-w-sm text-sm leading-6 text-white/75">
              Use the program filters above to compare rating snapshots across
              sports and roster groups.
            </p>
          </div>

          <div className="space-y-4">
            {ratingRows.map((row) => (
              <div key={row.label} className="space-y-2">
                <div className="flex items-center justify-between text-sm text-black">
                  <span className="font-medium">{row.label}</span>
                  <span className="text-zinc-500">{row.value}/100</span>
                </div>
                <div className="h-2.5 overflow-hidden rounded-full bg-zinc-100">
                  <div
                    className="h-full rounded-full bg-zinc-900"
                    style={{ width: `${row.value}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
