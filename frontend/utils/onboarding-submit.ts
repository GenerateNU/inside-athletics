"use client";

import type { OnboardingData } from "@/utils/onboarding";
import { getApiV1UserCurrent } from "@/api/clients/getApiV1UserCurrent";
import { postApiV1User } from "@/api/clients/postApiV1User";
import { getApiV1Roles } from "@/api/clients/getApiV1Roles";
import { postApiV1UserByIdRoles } from "@/api/clients/postApiV1UserByIdRoles";
import { getApiV1TagNameByName } from "@/api/clients/getApiV1TagNameByName";
import { postApiV1Tag } from "@/api/clients/postApiV1Tag";
import { postApiV1UserTag } from "@/api/clients/postApiV1UserTag";

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
    verified_athlete_status: data.role.role === "athlete" ? "pending" : "none",
  };
}

function delay(ms: number) {
  return new Promise((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

async function waitForCurrentUser(headers: Record<string, string>) {
  for (let attempt = 0; attempt < CURRENT_USER_SYNC_RETRIES; attempt += 1) {
    try {
      await getApiV1UserCurrent({ headers });
      return;
    } catch (error: any) {
      if (error?.response?.status !== 404) {
        throw new Error("Unable to verify current user.");
      }

      if (attempt < CURRENT_USER_SYNC_RETRIES - 1) {
        await delay(CURRENT_USER_SYNC_DELAY_MS);
      }
    }
  }

  throw new Error("User creation did not finish syncing. Please try again.");
}

async function getOrCreateTagId(
  tagName: string,
  headers: Record<string, string>,
): Promise<string | null> {
  try {
    const tag = await getApiV1TagNameByName(tagName, { headers });
    if (tag?.id) return tag.id;
  } catch (error: any) {
    if (error?.response?.status !== 404) {
      throw new Error(`Unable to look up tag "${tagName}".`);
    }
  }

  try {
    const created = await postApiV1Tag({ name: tagName, type: "topic" }, { headers });
    if (created?.id) return created.id;
    return null;
  } catch (error: any) {
    if (error?.response?.status === 409) {
      return getOrCreateTagId(tagName, headers);
    }
    throw new Error(`Unable to create tag "${tagName}".`);
  }
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
      if (!tagId) return;

      try {
        await postApiV1UserTag({ tag_id: tagId }, { headers });
      } catch (error: any) {
        if (error?.response?.status !== 409) {
          throw new Error(`Unable to follow tag "${tagName}".`);
        }
      }
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

  try {
    await getApiV1UserCurrent({ headers });
  } catch (error: any) {
    if (error?.response?.status !== 404) {
      throw new Error("Unable to verify current user.");
    }

    let userId: string;
    try {
      const created = await postApiV1User(payload, { headers });
      userId = created.id;
    } catch {
      throw new Error("Unable to create user.");
    }

    await waitForCurrentUser(headers);

    const rolesData = await getApiV1Roles(undefined, { headers });
    const roleName = data.plan.selectedPlan === "premium" ? "premium_user" : "user";
    const role = (rolesData?.roles ?? []).find((r) => r.name === roleName);
    if (role) {
      try {
        await postApiV1UserByIdRoles(userId, { role_id: role.id }, { headers });
      } catch (e: any) {
        if (e?.response?.status !== 409) {
          throw new Error("Unable to assign user role.");
        }
      }
    }
  }

  await syncSelectedTagFollows(data, headers);
}
