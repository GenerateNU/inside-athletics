"use client";

import { useEffect, useState, type ChangeEvent } from "react";
import { useRouter } from "next/navigation";
import { PlusIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";
import { useSession } from "@/utils/SessionContext";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const roleOptions = [
  { label: "Athlete", value: "athlete" },
  { label: "Prospective Athlete", value: "prospective-athlete" },
] as const;

const primarySports = [
  { label: "Basketball", value: "basketball" },
  { label: "Soccer", value: "soccer" },
  { label: "Track & Field", value: "track-and-field" },
  { label: "Volleyball", value: "volleyball" },
  { label: "Tennis", value: "tennis" },
  { label: "Swimming", value: "swimming" },
  { label: "Softball", value: "softball" },
  { label: "Baseball", value: "baseball" },
] as const;

const programOptions = [
  { label: "Women's", value: "womens" },
  { label: "Men's", value: "mens" },
] as const;

type UploadUrlPayload = {
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
  const [program, setProgram] = useState("");
  const [university, setUniversity] = useState("");
  const [collegeOptions, setCollegeOptions] = useState<string[]>([]);
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
    setProgram(data.preferences.program);
    setUniversity(data.preferences.university);
  }, [data.preferences, data.role, hydrated]);

  useEffect(() => {
    const accessToken = session?.access_token;

    if (role !== "athlete" || !accessToken) {
      setIsLoadingColleges(false);
      setCollegeError("");
      if (role !== "athlete") {
        setUniversity("");
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
          .map((college) => college.name.trim())
          .filter(Boolean);

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

    if (nextRole !== "athlete") {
      setUniversity("");
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
    primarySports.find((sport) => sport.value === primarySport)?.label ?? "";

  const canContinue = Boolean(
    role && primarySport && program && (role === "athlete" ? university : true),
  );

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">About you</h1>
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
            <SelectTrigger id="role" className="h-10 w-full text-sm">
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
              value={primarySport}
              onValueChange={(value) => {
                setPrimarySport(value ?? "");
              }}
            >
              <SelectTrigger id="primary-sport" className="h-10 w-full text-sm">
                <SelectValue placeholder="Select a primary sport">
                  {selectedPrimarySportLabel}
                </SelectValue>
              </SelectTrigger>
              <SelectContent>
                {primarySports.map((sport) => (
                  <SelectItem key={sport.value} value={sport.value}>
                    {sport.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        ) : null}

        {role && primarySport ? (
          <div className="space-y-3">
            <p className="block text-sm font-medium text-black">
              {role === "athlete"
                ? "Which team do you belong to?"
                : "Which program would you join?"}
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

        {role === "athlete" && primarySport && program ? (
          <div className="space-y-3">
            <label
              htmlFor="university"
              className="block text-sm font-medium text-black"
            >
              University
            </label>
            <Select
              value={university}
              onValueChange={(value) => {
                setUniversity(value ?? "");
              }}
            >
              <SelectTrigger id="university" className="h-10 w-full text-sm">
                <SelectValue
                  placeholder={
                    isLoadingColleges
                      ? "Loading universities..."
                      : "Select a university"
                  }
                />
              </SelectTrigger>
              <SelectContent>
                {collegeOptions.map((school) => (
                  <SelectItem key={school} value={school}>
                    {school}
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
              program,
              university: role === "athlete" ? university : "",
            });
            if (
              role === "athlete" &&
              data.verification.email.trim() &&
              !session?.user.email_confirmed_at
            ) {
              router.push(
                `/onboarding/verification/code?source=signup&role=${encodeURIComponent(role)}&email=${encodeURIComponent(data.verification.email)}`,
              );
              return;
            }

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
