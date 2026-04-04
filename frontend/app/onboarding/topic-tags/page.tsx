"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";

const topicCategories = [
  {
    name: "Athletics & Performance",
    tags: ["Training", "Game Highlights", "Skill Development", "Team Culture"],
  },
  {
    name: "Recruiting & Logistics",
    tags: ["Recruiting", "Scholarships", "Visits", "Eligibility"],
  },
  {
    name: "Health & Wellness",
    tags: ["Nutrition", "Recovery", "Mental Health", "Injury Prevention"],
  },
  {
    name: "Student-Athlete Life",
    tags: ["Academics", "Time Management", "Campus Life", "Leadership"],
  },
] as const;

export default function OnboardingTopicTagsPage() {
  const router = useRouter();
  const [selectedTags, setSelectedTags] = useState<string[]>([]);

  const toggleTag = (tag: string) => {
    setSelectedTags((current) =>
      current.includes(tag)
        ? current.filter((item) => item !== tag)
        : [...current, tag],
    );
  };

  const canContinue = selectedTags.length > 0;

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-3xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-3 text-center">
          <h1 className="text-4xl font-bold text-black">Choose Topic Tags</h1>
          <p className="text-sm text-gray-600">
            Select what you&apos;d like to see on your home feed
          </p>
          <div className="flex min-h-12 flex-wrap justify-center gap-2 rounded-xl border border-gray-200 px-4 py-3">
            {selectedTags.length > 0 ? (
              selectedTags.map((tag) => (
                <button
                  key={tag}
                  type="button"
                  className="inline-flex items-center gap-2 rounded-full px-3 py-1 text-sm font-medium text-white"
                  style={{
                    background:
                      "linear-gradient(180deg, #00804D 0%, #043D26 100%)",
                  }}
                  onClick={() => {
                    toggleTag(tag);
                  }}
                >
                  <XIcon className="size-3" />
                  <span>{tag}</span>
                </button>
              ))
            ) : (
              <span className="text-sm text-gray-500">No tags selected yet</span>
            )}
          </div>
        </div>

        <div className="space-y-6">
          {topicCategories.map((category) => (
            <div key={category.name} className="space-y-3">
              <h2 className="text-center text-sm font-semibold text-black">
                {category.name}
              </h2>
              <div className="flex flex-wrap gap-3">
                {category.tags.map((tag) => {
                  const isSelected = selectedTags.includes(tag);

                  return (
                    <Button
                      key={tag}
                      type="button"
                      variant="outline"
                      className="h-12 rounded-xl px-4 text-sm font-semibold"
                      style={{
                        borderColor: isSelected ? "transparent" : "#00804D",
                        background: isSelected
                          ? "linear-gradient(180deg, #00804D 0%, #043D26 100%)"
                          : "#FFFFFF",
                        color: isSelected ? "#FFFFFF" : "#000000",
                      }}
                      onClick={() => {
                        toggleTag(tag);
                      }}
                    >
                      {tag}
                    </Button>
                  );
                })}
              </div>
            </div>
          ))}
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            router.push("/onboarding/athletic-program-survey");
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
