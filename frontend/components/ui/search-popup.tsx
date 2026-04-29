"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Search } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import { useGetApiV1Sports, useListApiV1TagTypeByType } from "@/api/hooks";
import { useSession } from "@/utils/SessionContext";
import { Tag, TagButton } from "./tag-button";

function TagSection({
  header,
  tags,
  activeTags,
  maxTagNum,
  onToggle,
}: {
  header: string;
  tags: Tag[];
  activeTags: Tag[];
  maxTagNum: number;
  onToggle: (tag: Tag) => void;
}) {
  return (
    <div className="mt-4">
      <label className="block text-xl text-black font-bold">{header}</label>
      <div className="flex flex-wrap gap-2 mt-2">
        {tags.map((tag) => (
          <TagButton
            key={tag.id}
            tag={tag}
            active={activeTags.some((t) => t.id === tag.id)}
            onClick={() => onToggle(tag)}
          />
        ))}
      </div>
    </div>
  );
}

type TagType =
  | "sports"
  | "divisions"
  | "athletics_performance"
  | "health_wellness"
  | "student_athlete_life"
  | "recruiting_logistics";

type SearchPopupProps = {
  activeTags: Tag[];
  setActiveTagsAction: React.Dispatch<React.SetStateAction<Tag[]>>;
  onBackAction: () => void;
};

const TAG_SECTIONS: {
  header: string;
  type: TagType;
  max: number;
  group?: string;
}[] = [
  { header: "Sports", type: "sports", max: 1 },
  { header: "Divisions", type: "divisions", max: 3 },
  {
    header: "Athletics & Performance",
    type: "athletics_performance",
    max: 5,
    group: "Other Tags",
  },
  { header: "Health & Wellness", type: "health_wellness", max: 5 },
  { header: "Student Athlete Life", type: "student_athlete_life", max: 5 },
  { header: "Recruiting Logistics", type: "recruiting_logistics", max: 5 },
];

// One hook call per tag type — all run in parallel
function useAllTagSections() {
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const sports = useGetApiV1Sports(
    {},
    {
      query: { enabled },
      client: { headers: authHeaders },
    },
  );
  const divisions = useListApiV1TagTypeByType("divisions", {
    query: { enabled },
    client: { headers: authHeaders },
  });
  const athleticsPerf = useListApiV1TagTypeByType("athletics_performance", {
    query: { enabled },
    client: { headers: authHeaders },
  });
  const healthWellness = useListApiV1TagTypeByType("health_wellness", {
    query: { enabled },
    client: { headers: authHeaders },
  });
  const studentAthleteLife = useListApiV1TagTypeByType("student_athlete_life", {
    query: { enabled },
    client: { headers: authHeaders },
  });
  const recruitingLogistic = useListApiV1TagTypeByType("recruiting_logistics", {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const results: Record<TagType, Tag[]> = {
    sports: (sports.data?.sports ?? []).map((s) => ({
      id: s.id,
      name: s.name,
      type: "sports",
    })),
    divisions: divisions.data ?? [],
    athletics_performance: athleticsPerf.data ?? [],
    health_wellness: healthWellness.data ?? [],
    student_athlete_life: studentAthleteLife.data ?? [],
    recruiting_logistics: recruitingLogistic.data ?? [],
  };

  const loading =
    sports.isLoading ||
    sports.data === undefined ||
    [
      divisions,
      athleticsPerf,
      healthWellness,
      studentAthleteLife,
      recruitingLogistic,
    ].some((r) => r.isLoading || r.data === undefined);

  return { tagsByType: results, loading };
}

export default function SearchPopup({
  activeTags,
  setActiveTagsAction,
  onBackAction,
}: SearchPopupProps) {
  const [search, setSearch] = useState("");
  const { tagsByType, loading } = useAllTagSections();

  const filter = (tags: Tag[]) =>
    Array.isArray(tags)
      ? tags.filter((tag) =>
          tag.name.toLowerCase().includes(search.toLowerCase()),
        )
      : [];

  const toggleTag = (tag: Tag, max: number) => {
    setActiveTagsAction((prev) => {
      const exists = prev.find((t) => t.id === tag.id);
      if (exists) return prev.filter((t) => t.id !== tag.id);
      const sameType = prev.filter((t) => t.type === tag.type);
      if (sameType.length >= max) return prev;
      return [...prev, tag];
    });
  };

  const removeTag = (tag: Tag) => {
    setActiveTagsAction((prev) => prev.filter((t) => t.id !== tag.id));
  };

  const otherTagsIndex = TAG_SECTIONS.findIndex(
    (s) => s.group === "Other Tags",
  );

  return (
    <div className="flex w-[60vw] bg-white rounded-2xl justify-center py-10 px-10 max-h-[60vh] overflow-y-auto">
      <div className="w-full flex flex-col space-y-3 h-full">
        <div className="flex gap-4">
          <label className="block text-2xl text-black font-bold">Add Tag</label>
        </div>

        <div className="relative">
          <Search
            size={16}
            className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400"
          />
          <Input
            id="search_tags"
            name="search_tags"
            type="text"
            placeholder="Search Tags"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>

        {activeTags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {activeTags.map((tag) => (
              <TagButton
                key={tag.id}
                tag={tag}
                active={true}
                onClick={() => removeTag(tag)}
              />
            ))}
          </div>
        )}

        {loading ? (
          <div className="flex items-center justify-center min-h-[400px]">
            <Spinner className="size-6 text-gray-400" />
          </div>
        ) : (
          TAG_SECTIONS.map(({ header, type, max }, i) => (
            <div key={type}>
              {i === otherTagsIndex && (
                <div className="mt-8">
                  <label className="block text-2xl text-black font-bold">
                    Other Tags
                  </label>
                  <label className="block text-md text-[#001225]">
                    Select Max 5
                  </label>
                </div>
              )}

              <TagSection
                header={header}
                tags={filter(tagsByType[type])}
                activeTags={activeTags}
                maxTagNum={max}
                onToggle={(tag) => toggleTag(tag, max)}
              />

              {i < otherTagsIndex && (
                <hr className="border-t border-gray-300 mt-4" />
              )}
            </div>
          ))
        )}

        <div className="sticky bottom-0 flex justify-end pt-4">
          <Button
            variant="ghost"
            className="rounded-full bg-[#2C649A] hover:bg-[#245580] hover:text-[#F4F8FA] text-[#F4F8FA] flex items-center gap-2 px-4 py-1"
            onClick={onBackAction}
          >
            Done
          </Button>
        </div>
      </div>
    </div>
  );
}
