"use client";

import { useEffect, useRef, useState } from "react";

import { submitOnboardingUser } from "@/utils/onboarding-submit";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";

export function OnboardingComplete() {
  const session = useSession();
  const { data, hydrated, reset } = useOnboarding();
  const [error, setError] = useState("");
  const hasStarted = useRef(false);

  useEffect(() => {
    if (!hydrated || !session?.access_token || hasStarted.current) {
      return;
    }

    hasStarted.current = true;

    submitOnboardingUser(data, session.access_token, session.user.email)
      .then(() => {
        setError("");
        reset();
      })
      .catch((submissionError) => {
        hasStarted.current = false;
        setError(
          submissionError instanceof Error
            ? submissionError.message
            : "Unable to finish onboarding.",
        );
      });
  }, [data, hydrated, reset, session]);

  if (!error) {
    return null;
  }

  return (
    <p className="text-sm text-red-600" role="alert">
      {error}
    </p>
  );
}
