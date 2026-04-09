"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { PlusIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { useOnboarding } from "@/utils/onboarding";
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

export default function OnboardingRolePage() {
  const router = useRouter();
  const { data, hydrated, updateSection } = useOnboarding();
  const [role, setRole] = useState("");
  const [profileImage, setProfileImage] = useState<string | null>(null);

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setRole(data.role.role);
    setProfileImage(data.role.profileImage);
  }, [data.role.profileImage, data.role.role, hydrated]);

  const canContinue = Boolean(role);

  const handleRoleChange = (value: string | null) => {
    setRole(value ?? "");
  };

  const selectedRoleLabel =
    roleOptions.find((option) => option.value === role)?.label ?? "";

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
          <input
            id="profile-image"
            type="file"
            accept="image/*"
            className="hidden"
            onChange={(event) => {
              const file = event.target.files?.[0];
              if (!file) {
                setProfileImage(null);
                return;
              }
              setProfileImage(URL.createObjectURL(file));
            }}
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
            });
            router.push(
              `/onboarding/preferences?role=${encodeURIComponent(role)}`,
            );
          }}
          disabled={!canContinue}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
