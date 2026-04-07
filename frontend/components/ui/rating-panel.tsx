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
  { label: "Long-term Player Development", value: "5" },
  { label: "Academics cs Athletics Priority", value: "4" },
  { label: "Academic & Career Support Resources", value: "4" },
  { label: "Mental Health Priority", value: "4" },
  { label: "Competitive Environment / Winning Culture", value: "5" },
  { label: "Team Culture", value: "4" },
  { label: "Coaching Transparency & communication", value: "4" },
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
    <Card className="w-full max-w-3xl border border-black/10 bg-[#E8F1FA] py-0 shadow-[0_20px_60px_rgba(0,0,0,0.06)]">
      <CardHeader className="border-b border-black/8 px-6 py-5">
        <CardTitle className="text-left text-2xl font-semibold text-black">
          Rating
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-8 px-6 py-6">
        <div className="grid gap-4 md:grid-cols-2">
          <Select value={programGender} onValueChange={handleProgramGenderChange}>
            <SelectTrigger className="h-12 w-full rounded-md border-black/10 bg-zinc-50 px-4 text-sm text-black">
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
            <SelectTrigger className="h-12 w-full rounded-md border-black/10 bg-zinc-50 px-4 text-sm text-black">
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

        <div className="grid gap-4 md:grid-cols-2">
          {ratingRows.map((row) => (
            <div
              key={row.label}
              className="flex items-center gap-4 rounded-xl bg-[#E8F1FA] p-4"
            >
              <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-[#3E7DBB] text-base font-semibold text-white">
                {row.value}/5
              </div>
              <div className="text-sm font-medium leading-6 text-black">
                {row.label}
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
