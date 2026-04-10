"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { useGetApiV1Sports } from "@/api/hooks";
import { Card, CardContent } from "@/components/ui/card";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import type { GetAllSportsResponse } from "@/api/models/GetAllSportsResponse";
import { useSession } from "@/utils/SessionContext";

const genderOptions = ["Women's", "Men's", "Co-ed"] as const;

type AverageRatingsRow = {
  sport_id: string;
  college_id: string;
  player_dev: number;
  academics_athletics_priority: number;
  academic_career_resources: number;
  mental_health_priority: number;
  environment: number;
  culture: number;
  transparency: number;
  response_count: number;
};

type AverageRatingsResponse = {
  averages: AverageRatingsRow[];
};

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
  {
    key: "academics_athletics_priority",
    label: "Academics vs Athletics Priority",
  },
  {
    key: "academic_career_resources",
    label: "Academic & Career Support Resources",
  },
  { key: "mental_health_priority", label: "Mental Health Priority" },
  { key: "environment", label: "Competitive Environment / Winning Culture" },
  { key: "culture", label: "Team Culture" },
  { key: "transparency", label: "Coaching Transparency & Communication" },
];

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

function formatRating(value: number) {
  return value.toFixed(1).replace(/\.0$/, "");
}

async function getAverageRatings({
  accessToken,
  collegeId,
  sportId,
}: {
  accessToken: string;
  collegeId: string;
  sportId: string;
}) {
  const searchParams = new URLSearchParams({
    college_id: collegeId,
    sport_id: sportId,
  });
  const response = await fetch(
    `/api/v1/survey/averages?${searchParams.toString()}`,
    {
      headers: {
        Authorization: `Bearer ${accessToken}`,
        Accept: "application/json",
      },
    },
  );

  if (!response.ok) {
    throw new Error(`Unable to load ratings (${response.status})`);
  }

  const payload = (await response.json()) as unknown;
  const data = unwrapBody<AverageRatingsResponse>(payload);

  return data?.averages ?? [];
}

function aggregateRatings(rows: AverageRatingsRow[]): AggregatedRatings | null {
  if (rows.length === 0) {
    return null;
  }

  const totals = rows.reduce(
    (acc, row) => {
      const weight = row.response_count;

      acc.player_dev += row.player_dev * weight;
      acc.academics_athletics_priority +=
        row.academics_athletics_priority * weight;
      acc.academic_career_resources += row.academic_career_resources * weight;
      acc.mental_health_priority += row.mental_health_priority * weight;
      acc.environment += row.environment * weight;
      acc.culture += row.culture * weight;
      acc.transparency += row.transparency * weight;
      acc.response_count += weight;

      return acc;
    },
    {
      player_dev: 0,
      academics_athletics_priority: 0,
      academic_career_resources: 0,
      mental_health_priority: 0,
      environment: 0,
      culture: 0,
      transparency: 0,
      response_count: 0,
    },
  );

  if (totals.response_count === 0) {
    return null;
  }

  return {
    player_dev: totals.player_dev / totals.response_count,
    academics_athletics_priority:
      totals.academics_athletics_priority / totals.response_count,
    academic_career_resources:
      totals.academic_career_resources / totals.response_count,
    mental_health_priority:
      totals.mental_health_priority / totals.response_count,
    environment: totals.environment / totals.response_count,
    culture: totals.culture / totals.response_count,
    transparency: totals.transparency / totals.response_count,
    response_count: totals.response_count,
  };
}

export function RatingPanel({ collegeId }: { collegeId: string }) {
  const [programGender, setProgramGender] = useState<string>("");
  const [sportProgram, setSportProgram] = useState<string>("");
  const session = useSession();
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const handleProgramGenderChange = (value: string | null) => {
    setProgramGender(value ?? "");
  };

  const handleSportProgramChange = (value: string | null) => {
    setSportProgram(value ?? "");
  };

  const { data: sportsResponse, isLoading: isLoadingSports } =
    useGetApiV1Sports(undefined, {
      query: { enabled: !!session?.access_token },
      client: { headers: authHeaders },
    });

  const sports = unwrapBody<GetAllSportsResponse>(sportsResponse)?.sports ?? [];

  const {
    data: surveyRows = [],
    isLoading: isLoadingRatings,
    isError,
  } = useQuery({
    queryKey: ["survey-averages", collegeId, sportProgram],
    queryFn: () =>
      getAverageRatings({
        accessToken: session!.access_token,
        collegeId,
        sportId: sportProgram,
      }),
    enabled: !!session?.access_token && !!collegeId && !!sportProgram,
  });

  const aggregatedRatings = aggregateRatings(surveyRows);
  const helperText = !session?.access_token
    ? "Sign in to load survey ratings."
    : !sportProgram
      ? "Select a sport program to load ratings."
      : isLoadingRatings
        ? "Loading survey ratings..."
        : isError
          ? "Unable to load survey ratings."
          : !aggregatedRatings
            ? "No survey ratings yet for this sport."
            : `${aggregatedRatings.response_count} survey response${aggregatedRatings.response_count === 1 ? "" : "s"} included.`;

  return (
    <div className="w-full max-w-3xl space-y-4">
      <h2 className="text-left text-2xl font-semibold text-black">Rating</h2>
      <Card className="border border-black/10 bg-[#E8F1FA] py-0 shadow-[0_20px_60px_rgba(0,0,0,0.06)]">
        <CardContent className="space-y-8 px-6 py-6">
          <div className="grid gap-4 md:grid-cols-2">
            <Select
              value={programGender}
              onValueChange={handleProgramGenderChange}
            >
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

            <Select
              value={sportProgram}
              onValueChange={handleSportProgramChange}
            >
              <SelectTrigger className="h-12 w-full rounded-md border-black/10 bg-zinc-50 px-4 text-sm text-black">
                <SelectValue
                  placeholder={
                    isLoadingSports
                      ? "Loading sport programs..."
                      : "Select Sport Program"
                  }
                />
              </SelectTrigger>
              <SelectContent>
                {sports.map((option) => (
                  <SelectItem key={option.id} value={option.id}>
                    {option.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <p className="text-sm text-black/65">{helperText}</p>

          <div className="grid gap-4 md:grid-cols-2">
            {ratingMetrics.map((row) => (
              <div
                key={row.label}
                className="flex items-center gap-4 rounded-xl bg-[#E8F1FA] p-4"
              >
                <div className="flex h-16 w-16 shrink-0 items-center justify-center rounded-2xl bg-[#3E7DBB] text-base font-semibold text-white">
                  {aggregatedRatings
                    ? `${formatRating(aggregatedRatings[row.key])}/5`
                    : "--"}
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
