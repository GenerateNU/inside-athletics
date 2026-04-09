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

type UploadUrlPayload = {
  upload_url: string;
  key: string;
};

type ConfirmUploadPayload = {
  key: string;
  download_url: string;
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
  const [uploadError, setUploadError] = useState("");
  const [isUploadingProfileImage, setIsUploadingProfileImage] = useState(false);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setRole(data.role.role);
    setProfileImage(data.role.profileImage);
    setProfileImageKey(data.role.profileImageKey);
  }, [
    data.role.profileImage,
    data.role.profileImageKey,
    data.role.role,
    hydrated,
  ]);

  const canContinue = Boolean(role);

  const handleRoleChange = (value: string | null) => {
    setRole(value ?? "");
  };

  const selectedRoleLabel =
    roleOptions.find((option) => option.value === role)?.label ?? "";

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

    if (!session?.access_token) {
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
          Authorization: `Bearer ${session.access_token}`,
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
            Authorization: `Bearer ${session.access_token}`,
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
                className="absolute bottom-0 right-0 flex size-6 items-center justify-center rounded-full shadow-sm"
                style={{ backgroundColor: "#3E7DBB", color: "#FFFFFF" }}
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

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("role", {
              role,
              profileImage,
              profileImageKey,
            });
            router.push(
              `/onboarding/preferences?role=${encodeURIComponent(role)}`,
            );
          }}
          disabled={!canContinue || isUploadingProfileImage}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
