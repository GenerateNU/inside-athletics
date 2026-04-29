"use client";

import { useState } from "react";

import { useGetApiV1Sports, useGetApiV1SurveyAverages } from "@/api/hooks";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useSession } from "@/utils/SessionContext";
import type { AverageRatingsRow } from "@/api/models/AverageRatingsRow";

type RatingMetric = {
  key: keyof AggregatedRatings;
  label: string;
};

type AggregatedRatings = {
  player_dev: number;
  academics_athletics_priority: number;
  academic_career_resources: number;
  mental_health_priority: number;
  environment: number;
  culture: number;
  transparency: number;
  response_count: number;
};

const ratingMetrics: RatingMetric[] = [
  { key: "player_dev", label: "Long-term Player Development" },
  { key: "academics_athletics_priority", label: "Academics vs Athletics Priority" },
  { key: "academic_career_resources", label: "Academic & Career Support Resources" },
  { key: "mental_health_priority", label: "Mental Health Priority" },
  { key: "environment", label: "Competitive Environment / Winning Culture" },
  { key: "culture", label: "Team Culture" },
  { key: "transparency", label: "Coaching Transparency & Communication" },
];

function formatRating(value: number) {
  return value.toFixed(1).replace(/\.0$/, "");
}

function aggregateRatings(rows: AverageRatingsRow[]): AggregatedRatings | null {
  if (rows.length === 0) return null;

  const totals = rows.reduce(
    (acc, row) => {
      const w = row.response_count;
      acc.player_dev += row.player_dev * w;
      acc.academics_athletics_priority += row.academics_athletics_priority * w;
      acc.academic_career_resources += row.academic_career_resources * w;
      acc.mental_health_priority += row.mental_health_priority * w;
      acc.environment += row.environment * w;
      acc.culture += row.culture * w;
      acc.transparency += row.transparency * w;
      acc.response_count += w;
      return acc;
    },
    { player_dev: 0, academics_athletics_priority: 0, academic_career_resources: 0, mental_health_priority: 0, environment: 0, culture: 0, transparency: 0, response_count: 0 },
  );

  if (totals.response_count === 0) return null;

  return {
    player_dev: totals.player_dev / totals.response_count,
    academics_athletics_priority: totals.academics_athletics_priority / totals.response_count,
    academic_career_resources: totals.academic_career_resources / totals.response_count,
    mental_health_priority: totals.mental_health_priority / totals.response_count,
    environment: totals.environment / totals.response_count,
    culture: totals.culture / totals.response_count,
    transparency: totals.transparency / totals.response_count,
    response_count: totals.response_count,
  };
}

export function RatingPanel({ collegeId }: { collegeId: string }) {
  const [sportProgram, setSportProgram] = useState<string>("");
  const session = useSession();
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;
  const enabled = !!session?.access_token;

  const { data: sportsData, isLoading: isLoadingSports } = useGetApiV1Sports(
    undefined,
    { query: { enabled }, client: { headers: authHeaders } },
  );

  const sports = sportsData?.sports ?? [];

  const averagesParams = sportProgram
    ? { college_id: collegeId, sport_id: sportProgram }
    : { college_id: collegeId };

  const { data: averagesData, isLoading: isLoadingRatings, isError } =
    useGetApiV1SurveyAverages(averagesParams, {
      query: { enabled },
      client: { headers: authHeaders },
    });

  const rows = averagesData?.averages ?? [];
  const aggregatedRatings = aggregateRatings(rows);

  const helperText = !session?.access_token
    ? "Sign in to load survey ratings."
    : isLoadingRatings
      ? "Loading survey ratings..."
      : isError
        ? "Unable to load survey ratings."
        : !aggregatedRatings
          ? "No survey ratings yet for this college."
          : `${aggregatedRatings.response_count} survey response${aggregatedRatings.response_count === 1 ? "" : "s"} included${sportProgram ? " for this sport" : ""}.`;

  return (
    <div className="w-full space-y-4">
      <Card className="border border-black/10 bg-white py-0 shadow-[0_20px_60px_rgba(0,0,0,0.06)] rounded-2xl">
        <CardContent className="space-y-8 px-6 py-6">
          <Select value={sportProgram} onValueChange={(v) => setSportProgram(v ?? "")}>
            <SelectTrigger className="h-12 w-full rounded-md border-black/10 bg-zinc-50 px-4 text-sm text-black">
              <SelectValue placeholder={isLoadingSports ? "Loading sport programs..." : "Select Sport Program"}>
                {sportProgram ? (sports.find((s) => s.id === sportProgram)?.name ?? "") : null}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              {sports.map((sport) => (
                <SelectItem key={sport.id} value={sport.id}>
                  {sport.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>

          <p className="text-sm text-black/65">{helperText}</p>

          <div className="grid gap-4 md:grid-cols-2">
            {ratingMetrics.map((row) => (
              <div
                key={row.label}
                className="flex items-center gap-4 rounded-xl bg-[#E8F1FA] p-4"
              >
                <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-[#3E7DBB] text-base font-semibold text-white">
                  {aggregatedRatings ? `${formatRating(aggregatedRatings[row.key])}/5` : "--"}
                </div>
                <div className="text-sm font-medium leading-6 text-black">
                  {row.label}
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
