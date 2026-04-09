"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import { submitOnboardingUser } from "@/utils/onboarding-submit";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";

const surveyQuestions = [
  "To what extent does your program prioritize long-term player development (athletic, personal, and leadership growth over your college career)?",
  "When conflicts arise, how often are academics prioritized over athletics within your program?",
  "How would you rate the quality and accessibility of academic and career support resources available to athletes at your school (e.g., advisors, career guidance, networking)?",
  "How seriously does your program treat athlete mental health, including coach awareness, available resources, and overall support?",
  "How would you describe the overall competitiveness of your program (e.g., performance standards, emphasis on winning, playing time)?",
  "How would you rate the overall team culture within your program (including trust, accountability, support and camaraderie among teammates, and alignment with team values)?",
  "How clear and transparent are coaches about expectations, roles, and feedback?",
] as const;

const ratingOptions = ["1", "2", "3", "4", "5"] as const;

export default function OnboardingAthleticProgramSurveyPage() {
  const router = useRouter();
  const session = useSession();
  const { data, hydrated, reset } = useOnboarding();
  const [responses, setResponses] = useState<Record<number, string>>({});
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setResponses(data.survey.responses);
  }, [data.survey.responses, hydrated]);

  const finishOnboarding = async () => {
    if (!session?.access_token) {
      setError("You need an active session before finishing onboarding.");
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      await submitOnboardingUser(
        {
          ...data,
          survey: {
            responses,
          },
        },
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
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-3xl space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-3 text-center">
          <h1 className="text-4xl font-bold text-black">
            Athletic Program Survey
          </h1>
          <div className="flex flex-col items-center justify-center gap-3 sm:flex-row">
            <p className="text-sm text-gray-600">
              You can skip this survey and return later via your profile page
            </p>
            <Button
              type="button"
              variant="outline"
              className="rounded-xl px-4 text-sm font-semibold"
              style={{
                borderColor: "#2C649A",
                color: "#2C649A",
              }}
              onClick={finishOnboarding}
              disabled={isSubmitting}
            >
              Skip for now
            </Button>
          </div>
          <p className="text-sm text-gray-600">
            Your responses to these questions will live and be visible to others
            on your profile page. Please rank the following qualities from 1 - 5
            regarding your sports team.
          </p>
          {error ? (
            <p className="text-sm text-red-600" role="alert">
              {error}
            </p>
          ) : null}
        </div>

        <div className="space-y-6">
          {surveyQuestions.map((question, index) => (
            <div
              key={question}
              className="space-y-4 rounded-xl border border-gray-200 px-5 py-5"
            >
              <p className="text-sm font-medium text-black">
                {index + 1}. {question}
              </p>
              <div className="flex flex-wrap gap-3">
                {ratingOptions.map((option) => {
                  const isSelected = responses[index] === option;

                  return (
                    <Button
                      key={option}
                      type="button"
                      variant="outline"
                      className="h-11 min-w-16 rounded-xl px-6 text-sm font-semibold"
                      style={{
                        borderColor: isSelected ? "transparent" : "#00804D",
                        background: isSelected
                          ? "linear-gradient(180deg, #00804D 0%, #043D26 100%)"
                          : "#FFFFFF",
                        color: isSelected ? "#FFFFFF" : "#000000",
                      }}
                      onClick={() => {
                        setResponses((current) => ({
                          ...current,
                          [index]: option,
                        }));
                      }}
                    >
                      {option}
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
          onClick={finishOnboarding}
          disabled={isSubmitting}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
