"use client";

import type { OnboardingData } from "@/utils/onboarding";

type TagPayload = {
  id: string;
  name: string;
};

const CURRENT_USER_SYNC_RETRIES = 8;
const CURRENT_USER_SYNC_DELAY_MS = 250;

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
    data.verification.institutionEmail ||
    data.verification.email ||
    sessionEmail ||
    "";
  const { firstName, lastName } = splitName(onboardingName);

  if (!firstName || !data.account.username || !email) {
    return null;
  }

  return {
    first_name: firstName,
    last_name: lastName || firstName,
    email,
    username: data.account.username,
    profile_image_key: data.role.profileImageKey,
    account_type: data.plan.selectedPlan === "premium",
    verified_athlete_status: data.role.role === "athlete" ? "pending" : "none",
  };
}

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

function delay(ms: number) {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

async function waitForCurrentUser(headers: Record<string, string>) {
  for (let attempt = 0; attempt < CURRENT_USER_SYNC_RETRIES; attempt += 1) {
    const currentUserResponse = await fetch("/api/v1/user/current", {
      headers,
    });

    if (currentUserResponse.ok) {
      return;
    }

    if (currentUserResponse.status !== 404) {
      throw new Error("Unable to verify current user.");
    }

    if (attempt < CURRENT_USER_SYNC_RETRIES - 1) {
      await delay(CURRENT_USER_SYNC_DELAY_MS);
    }
  }

  throw new Error("User creation did not finish syncing. Please try again.");
}

async function getOrCreateTagId(
  tagName: string,
  headers: Record<string, string>,
) {
  const encodedName = encodeURIComponent(tagName);
  const existingTagResponse = await fetch(`/api/v1/tag/name/${encodedName}`, {
    headers,
  });

  if (existingTagResponse.ok) {
    const existingTag = unwrapBody<TagPayload>(
      await existingTagResponse.json(),
    );

    if (existingTag?.id) {
      return existingTag.id;
    }
  } else if (existingTagResponse.status !== 404) {
    throw new Error(`Unable to look up tag "${tagName}".`);
  }

  const createTagResponse = await fetch("/api/v1/tag/", {
    method: "POST",
    headers,
    body: JSON.stringify({
      name: tagName,
    }),
  });

  if (createTagResponse.status === 409) {
    return getOrCreateTagId(tagName, headers);
  }

  if (!createTagResponse.ok) {
    throw new Error(`Unable to create tag "${tagName}".`);
  }

  const createdTag = unwrapBody<TagPayload>(await createTagResponse.json());

  if (!createdTag?.id) {
    throw new Error(`Unable to read tag "${tagName}" from response.`);
  }

  return createdTag.id;
}

async function syncSelectedTagFollows(
  data: OnboardingData,
  headers: Record<string, string>,
) {
  const selectedTags = [
    ...new Set(data.topicTags.selectedTags.map((tag) => tag.trim())),
  ].filter(Boolean);

  if (!selectedTags.length) {
    return;
  }

  await Promise.all(
    selectedTags.map(async (tagName) => {
      const tagId = await getOrCreateTagId(tagName, headers);
      const followResponse = await fetch("/api/v1/user/tag/", {
        method: "POST",
        headers,
        body: JSON.stringify({
          tag_id: tagId,
        }),
      });

      if (followResponse.ok || followResponse.status === 409) {
        return;
      }

      throw new Error(`Unable to follow tag "${tagName}".`);
    }),
  );
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

  if (!currentUserResponse.ok) {
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

    await waitForCurrentUser(headers);
  }

  await syncSelectedTagFollows(data, headers);
}
