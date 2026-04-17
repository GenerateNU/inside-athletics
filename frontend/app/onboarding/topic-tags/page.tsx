"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { XIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";
import { submitOnboardingUser } from "@/utils/onboarding-submit";
import { useSession } from "@/utils/SessionContext";

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
  const session = useSession();
  const { data, hydrated, updateSection, reset } = useOnboarding();
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setSelectedTags(data.topicTags.selectedTags);
  }, [data.topicTags.selectedTags, hydrated]);

  const toggleTag = (tag: string) => {
    setSelectedTags((current) =>
      current.includes(tag)
        ? current.filter((item) => item !== tag)
        : [...current, tag],
    );
  };

  const canContinue = selectedTags.length > 0 && !isSubmitting;

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-3xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-3 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">Choose Topic Tags</h1>
          <p className="text-sm text-gray-600">
            Select what you&apos;d like to see on your home feed
          </p>
          <div className="flex min-h-12 flex-wrap justify-center gap-2 rounded-xl border border-gray-200 px-4 py-3">
            {selectedTags.length > 0 ? (
              selectedTags.map((tag) => (
                <button
                  key={tag}
                  type="button"
                  className="inline-flex items-center gap-2 rounded-full border border-[#7F8C2D] bg-[#D4E94B80] px-3 py-1 text-sm font-medium text-black"
                  onClick={() => {
                    toggleTag(tag);
                  }}
                >
                  <XIcon className="size-3" />
                  <span>{tag}</span>
                </button>
              ))
            ) : (
              <span className="text-sm text-gray-500">
                No tags selected yet
              </span>
            )}
          </div>
        </div>

        <div className="space-y-6">
          {topicCategories.map((category) => (
            <div key={category.name} className="space-y-3">
              <h2 className="text-center text-sm font-semibold text-[#001F3E]">
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
                      className={`h-12 rounded-xl px-4 text-sm font-semibold text-black ${
                        isSelected
                          ? "border-[#7F8C2D] bg-[#D4E94B80]"
                          : "border-[#D4E94B] bg-[#FCFDF1]"
                      }`}
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

        {error ? (
          <p className="text-center text-sm text-red-600" role="alert">
            {error}
          </p>
        ) : null}
        <Button
          type="button"
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={async () => {
            updateSection("topicTags", { selectedTags });

            if (data.role.role !== "prospective-athlete") {
              router.push("/onboarding/athletic-program-survey");
              return;
            }

            if (!session?.access_token) {
              setError("You need an active session before finishing onboarding.");
              return;
            }

            setIsSubmitting(true);
            setError("");

            try {
              await submitOnboardingUser(
                { ...data, topicTags: { selectedTags } },
                session.access_token,
                session.user.email,
              );
              reset();
              router.push("/");
            } catch (submissionError) {
              setError(
                submissionError instanceof Error
                  ? submissionError.message
                  : "Unable to finish onboarding.",
              );
            } finally {
              setIsSubmitting(false);
            }
          }}
          disabled={!canContinue}
        >
          {isSubmitting ? "Finishing..." : "Continue"}
        </Button>
      </div>
    </div>
  );
}
