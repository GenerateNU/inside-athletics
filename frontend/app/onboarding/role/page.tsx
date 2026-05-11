"use client";

import { useEffect, useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";
import { PlusIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";
import { useGetApiV1Sports } from "@/api/hooks";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const roleOptions = [
  { label: "Current college athlete", value: "current-college-athlete" },
  { label: "Former college athlete", value: "former-college-athlete" },
  { label: "Youth athlete", value: "youth-athlete" },
  { label: "Parent of an athlete", value: "parent-of-athlete" },
  { label: "Other", value: "other" },
] as const;

const ROLES_NEEDING_UNIVERSITY = new Set([
  "current-college-athlete",
  "former-college-athlete",
  "youth-athlete",
  "parent-of-athlete",
]);

function roleNeedsUniversity(role: string) {
  return ROLES_NEEDING_UNIVERSITY.has(role);
}

const programOptions = [
  { label: "Women's", value: "womens" },
  { label: "Men's", value: "mens" },
] as const;

type UploadUrlPayload = {
  document_id: string;
  upload_url: string;
  key: string;
};

type ConfirmUploadPayload = {
  key: string;
  download_url: string;
};

type CollegeResponse = {
  id: string;
  name: string;
};

type CollegeListPayload = {
  colleges: CollegeResponse[];
  total: number;
};

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

function sanitizeFileName(fileName: string) {
  return fileName.toLowerCase().replace(/[^a-z0-9._-]+/g, "-");
}

export default function OnboardingRolePage() {
  const router = useRouter();
  const session = useSession();
  const { data, hydrated, updateSection } = useOnboarding();
  const [role, setRole] = useState("");
  const [profileImage, setProfileImage] = useState<string | null>(null);
  const [profileImageKey, setProfileImageKey] = useState<string | null>(null);
  const [primarySport, setPrimarySport] = useState("");
  const [primarySportId, setPrimarySportId] = useState("");
  const [program, setProgram] = useState("");
  const [university, setUniversity] = useState("");
  const [universityId, setUniversityId] = useState("");
  const [collegeOptions, setCollegeOptions] = useState<
    Array<{ id: string; name: string }>
  >([]);
  const [isLoadingColleges, setIsLoadingColleges] = useState(false);
  const [collegeError, setCollegeError] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [isUploadingProfileImage, setIsUploadingProfileImage] = useState(false);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setRole(data.role.role);
    setProfileImage(data.role.profileImage);
    setProfileImageKey(data.role.profileImageKey);
    setPrimarySport(data.preferences.primarySport);
    setPrimarySportId(data.preferences.primarySportId);
    setProgram(data.preferences.program);
    setUniversity(data.preferences.university);
    setUniversityId(data.preferences.universityId);
  }, [data.preferences, data.role, hydrated]);

  const sportsQuery = useGetApiV1Sports(
    { limit: 200, offset: 0 },
    {
      query: { enabled: Boolean(session?.access_token) },
      client: { headers: session?.access_token ? { Authorization: `Bearer ${session.access_token}` } : undefined },
    },
  );
  const sportOptions = (sportsQuery.data?.sports ?? [])
    .map((sport) => ({ id: sport.id, name: sport.name }))
    .filter((sport) => sport.id && sport.name);

  useEffect(() => {
    const accessToken = session?.access_token;

    if (!roleNeedsUniversity(role) || !accessToken) {
      setIsLoadingColleges(false);
      setCollegeError("");
      if (!roleNeedsUniversity(role)) {
        setUniversity("");
        setUniversityId("");
      }
      return;
    }

    let cancelled = false;

    async function loadColleges() {
      try {
        setCollegeError("");
        setIsLoadingColleges(true);

        const response = await fetch("/api/v1/college/?limit=500", {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        });

        if (!response.ok) {
          const errorText = await response.text();
          const trimmedErrorText = errorText.trim();
          throw new Error(
            trimmedErrorText
              ? `Unable to load colleges (${response.status}): ${trimmedErrorText}`
              : `Unable to load colleges (${response.status} ${response.statusText}).`,
          );
        }

        const payload = unwrapBody<CollegeListPayload>(await response.json());
        const colleges = (payload?.colleges ?? [])
          .map((college) => ({ id: college.id, name: college.name.trim() }))
          .filter((college) => college.id && college.name);

        if (cancelled) {
          return;
        }

        setCollegeOptions(colleges);
      } catch (error) {
        if (cancelled) {
          return;
        }

        setCollegeError(
          error instanceof Error ? error.message : "Unable to load colleges.",
        );
      } finally {
        if (!cancelled) {
          setIsLoadingColleges(false);
        }
      }
    }

    loadColleges();

    return () => {
      cancelled = true;
    };
  }, [role, session?.access_token]);

  const handleRoleChange = (value: string | null) => {
    const nextRole = value ?? "";
    setRole(nextRole);

    if (!roleNeedsUniversity(nextRole)) {
      setUniversity("");
      setUniversityId("");
    }
  };

  const handleProfileImageChange = async (
    event: ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];

    if (!file) {
      setProfileImage(null);
      setProfileImageKey(null);
      setUploadError("");
      return;
    }

    const accessToken = session?.access_token;

    if (!accessToken) {
      setUploadError("You need an active session to upload a profile image.");
      event.target.value = "";
      return;
    }

    setIsUploadingProfileImage(true);
    setUploadError("");

    try {
      const extension = file.name.includes(".")
        ? file.name.slice(file.name.lastIndexOf("."))
        : "";
      const key = `profiles/onboarding/${Date.now()}-${sanitizeFileName(
        `${crypto.randomUUID()}${extension}`,
      )}`;

      const uploadUrlResponse = await fetch("/api/v1/content/upload-url", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${accessToken}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          key,
          fileName: file.name,
          fileType: file.type || "application/octet-stream",
        }),
      });

      if (!uploadUrlResponse.ok) {
        throw new Error("Unable to prepare profile image upload.");
      }

      const uploadUrlPayload = unwrapBody<UploadUrlPayload>(
        await uploadUrlResponse.json(),
      );

      if (!uploadUrlPayload?.upload_url || !uploadUrlPayload.key) {
        throw new Error("Upload URL response was incomplete.");
      }

      const directUploadResponse = await fetch(uploadUrlPayload.upload_url, {
        method: "PUT",
        headers: {
          "Content-Type": file.type || "application/octet-stream",
        },
        body: file,
      });

      if (!directUploadResponse.ok) {
        throw new Error("Unable to upload profile image.");
      }

      const confirmUploadResponse = await fetch(
        "/api/v1/content/confirm-upload",
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${accessToken}`,
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            key: uploadUrlPayload.key,
          }),
        },
      );

      if (!confirmUploadResponse.ok) {
        throw new Error("Unable to confirm profile image upload.");
      }

      const confirmUploadPayload = unwrapBody<ConfirmUploadPayload>(
        await confirmUploadResponse.json(),
      );

      if (!confirmUploadPayload?.download_url || !confirmUploadPayload.key) {
        throw new Error("Confirmed upload response was incomplete.");
      }

      setProfileImage(confirmUploadPayload.download_url);
      setProfileImageKey(confirmUploadPayload.key);
      updateSection("role", {
        profileImage: confirmUploadPayload.download_url,
        profileImageKey: confirmUploadPayload.key,
      });
    } catch (error) {
      setUploadError(
        error instanceof Error
          ? error.message
          : "Unable to upload profile image.",
      );
      event.target.value = "";
    } finally {
      setIsUploadingProfileImage(false);
    }
  };

  const selectedRoleLabel =
    roleOptions.find((option) => option.value === role)?.label ?? "";
  const selectedPrimarySportLabel =
    sportOptions.find((sport) => sport.id === primarySportId)?.name ?? primarySport;

  const canContinue = Boolean(
    role && primarySport && program && (roleNeedsUniversity(role) ? university : true),
  );

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">About you</h1>
          <p className="text-sm text-gray-600">Tell us about yourself!</p>
        </div>

        <div className="flex flex-col items-center gap-3">
          <label htmlFor="profile-image" className="cursor-pointer">
            <div className="relative pb-3 pr-3">
              <div className="flex size-28 items-center justify-center overflow-hidden rounded-full bg-gray-200">
                {profileImage ? (
                  <img
                    src={profileImage}
                    alt="Profile preview"
                    className="h-full w-full object-cover"
                  />
                ) : (
                  <div className="h-full w-full bg-gray-200" />
                )}
              </div>
              <div
                className="absolute bottom-0 right-0 flex size-6 items-center justify-center rounded-full bg-[#3E7DBB] text-white shadow-sm"
              >
                <PlusIcon className="size-3" />
              </div>
            </div>
          </label>
          {isUploadingProfileImage ? (
            <p className="text-sm text-gray-600">Uploading profile image...</p>
          ) : null}
          {uploadError ? (
            <p className="text-sm text-red-600" role="alert">
              {uploadError}
            </p>
          ) : null}
          <input
            id="profile-image"
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleProfileImageChange}
          />
        </div>

        <div className="space-y-3">
          <label
            htmlFor="role"
            className="block text-sm font-medium text-black"
          >
            Role
          </label>
          <Select value={role} onValueChange={handleRoleChange}>
            <SelectTrigger
              id="role"
              className="h-10 w-full border-[#3E7DBB] text-sm"
            >
              <SelectValue placeholder="Select a role">
                {selectedRoleLabel}
              </SelectValue>
            </SelectTrigger>
            <SelectContent>
              {roleOptions.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {role ? (
          <div className="space-y-3">
            <label
              htmlFor="primary-sport"
              className="block text-sm font-medium text-black"
            >
              Primary Sport
            </label>
            <Select
              value={primarySportId}
              onValueChange={(value) => {
                const nextId = value ?? "";
                setPrimarySportId(nextId);
                const match = sportOptions.find((sport) => sport.id === nextId);
                setPrimarySport(match?.name ?? "");
              }}
            >
              <SelectTrigger
                id="primary-sport"
                className="h-10 w-full border-[#3E7DBB] text-sm"
              >
                <SelectValue
                  placeholder={
                    sportsQuery.isLoading
                      ? "Loading sports..."
                      : "Select a primary sport"
                  }
                >
                  {selectedPrimarySportLabel}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {sportOptions.map((sport) => (
                  <SelectItem key={sport.id} value={sport.id}>
                    {sport.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        ) : null}

        {role && primarySport ? (
          <div className="space-y-3">
            <p className="block text-sm font-medium text-black">
              {role === "current-college-athlete"
                ? "Which team do you belong to?"
                : role === "former-college-athlete"
                ? "Which team did you play for?"
                : role === "parent-of-athlete"
                ? "Which program does your athlete belong to?"
                : "Which program are you interested in?"}
            </p>
            <div className="grid grid-cols-2 gap-3">
              {programOptions.map((option) => {
                const isSelected = program === option.value;

                return (
                  <Button
                    key={option.value}
                    type="button"
                    variant="outline"
                    className={`h-12 rounded-xl border text-sm font-semibold text-black ${
                      isSelected
                        ? "border-[#7F8C2D] bg-[#D4E94B80]"
                        : "border-[#D4E94B] bg-[#FCFDF1]"
                    }`}
                    onClick={() => {
                      setProgram(option.value);
                    }}
                  >
                    {option.label}
                  </Button>
                );
              })}
            </div>
          </div>
        ) : null}

        {roleNeedsUniversity(role) && primarySport && program ? (
          <div className="space-y-3">
            <label
              htmlFor="university"
              className="block text-sm font-medium text-black"
            >
              University
            </label>
            <Select
              value={universityId}
              onValueChange={(value) => {
                const nextId = value ?? "";
                setUniversityId(nextId);
                const match = collegeOptions.find(
                  (school) => school.id === nextId,
                );
                setUniversity(match?.name ?? "");
              }}
            >
              <SelectTrigger
                id="university"
                className="h-10 w-full border-[#3E7DBB] text-sm"
              >
                <SelectValue
                  placeholder={
                    isLoadingColleges
                      ? "Loading universities..."
                      : "Select a university"
                  }
                >
                  {university}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {collegeOptions.map((school) => (
                  <SelectItem key={school.id} value={school.id}>
                    {school.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
            {collegeError ? (
              <p className="text-sm text-red-600" role="alert">
                {collegeError}
              </p>
            ) : null}
          </div>
        ) : null}

        <Button
          type="button"
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={() => {
            updateSection("role", {
              role,
              profileImage,
              profileImageKey,
            });
            updateSection("preferences", {
              primarySport,
              primarySportId,
              program,
              university: roleNeedsUniversity(role) ? university : "",
              universityId: roleNeedsUniversity(role) ? universityId : "",
            });
            router.push(`/onboarding/legal?role=${encodeURIComponent(role)}`);
          }}
          disabled={!canContinue || isUploadingProfileImage}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
