"use client";

import { useEffect, useRef, useState } from "react";

import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";

function splitName(fullName: string) {
  const trimmed = fullName.trim();

  if (!trimmed) {
    return { firstName: "", lastName: "" };
  }

  const parts = trimmed.split(/\s+/);
  const firstName = parts[0] ?? "";
  const lastName = parts.slice(1).join(" ");

  return { firstName, lastName };
}

function buildCreateUserPayload(
  onboardingName: string,
  username: string,
  email: string,
  role: string,
  selectedPlan: string,
) {
  const { firstName, lastName } = splitName(onboardingName);

  if (!firstName || !username || !email) {
    return null;
  }

  return {
    first_name: firstName,
    last_name: lastName || firstName,
    email,
    username,
    account_type: selectedPlan === "premium",
    verified_athlete_status: role === "athlete" ? "pending" : "none",
  };
}

export function OnboardingComplete() {
  const session = useSession();
  const { data, hydrated, reset } = useOnboarding();
  const [error, setError] = useState("");
  const hasStarted = useRef(false);

  useEffect(() => {
    if (!hydrated || !session?.access_token || hasStarted.current) {
      return;
    }

    const onboardingName =
      data.verification.fullName || data.verification.name || data.account.name;
    const email =
      data.verification.institutionEmail ||
      data.verification.email ||
      session.user.email ||
      "";

    const payload = buildCreateUserPayload(
      onboardingName,
      data.account.username,
      email,
      data.role.role,
      data.plan.selectedPlan,
    );

    if (!payload) {
      return;
    }

    hasStarted.current = true;

    const headers = {
      Authorization: `Bearer ${session.access_token}`,
      "Content-Type": "application/json",
    };

    const submitOnboarding = async () => {
      const currentUserResponse = await fetch("/api/v1/user/current", {
        headers,
      });

      if (currentUserResponse.ok) {
        reset();
        return;
      }

      if (currentUserResponse.status !== 404) {
        throw new Error("Unable to verify current user.");
      }

      const createUserResponse = await fetch("/api/v1/user", {
        method: "POST",
        headers,
        body: JSON.stringify(payload),
      });

      if (!createUserResponse.ok) {
        throw new Error("Unable to create user.");
      }

      reset();
    };

    submitOnboarding().catch((submissionError) => {
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
