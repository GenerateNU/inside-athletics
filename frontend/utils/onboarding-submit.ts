"use client";

import type { OnboardingData } from "@/utils/onboarding";

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

function buildCreateUserPayload(data: OnboardingData, sessionEmail?: string) {
  const onboardingName =
    data.verification.fullName || data.verification.name || data.account.name;
  const email =
    data.verification.institutionEmail || data.verification.email || sessionEmail || "";
  const { firstName, lastName } = splitName(onboardingName);

  if (!firstName || !data.account.username || !email) {
    return null;
  }

  return {
    first_name: firstName,
    last_name: lastName || firstName,
    email,
    username: data.account.username,
    account_type: data.plan.selectedPlan === "premium",
    verified_athlete_status: data.role.role === "athlete" ? "pending" : "none",
  };
}

export async function submitOnboardingUser(
  data: OnboardingData,
  accessToken: string,
  sessionEmail?: string,
) {
  const payload = buildCreateUserPayload(data, sessionEmail);

  if (!payload) {
    throw new Error("Missing onboarding fields required to create the user.");
  }

  const headers = {
    Authorization: `Bearer ${accessToken}`,
    "Content-Type": "application/json",
  };

  const currentUserResponse = await fetch("/api/v1/user/current", {
    headers,
  });

  if (currentUserResponse.ok) {
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
}
