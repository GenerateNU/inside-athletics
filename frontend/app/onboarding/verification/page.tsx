"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOnboarding } from "@/utils/onboarding";

export default function OnboardingVerificationPage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [fullName, setFullName] = useState("");
  const [institutionEmail, setInstitutionEmail] = useState("");
  const [school, setSchool] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    setName(data.verification.name);
    setEmail(data.verification.email);
    setFullName(data.verification.fullName);
    setInstitutionEmail(data.verification.institutionEmail);
    setSchool(data.verification.school);
  }, [data.verification, hydrated]);

  const isAthlete = role === "athlete";

  const canContinue = isAthlete
    ? Boolean(fullName.trim() && institutionEmail.trim() && school.trim())
    : Boolean(name.trim() && email.trim());

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Verification</h1>
          <p className="text-sm text-gray-600">
            {isAthlete
              ? "Confirm your athlete details so we can verify your profile."
              : "Confirm your details so we can verify your profile."}
          </p>
        </div>

        {isAthlete ? (
          <>
            <div className="space-y-3">
              <label
                htmlFor="full-name"
                className="block text-sm font-medium text-black"
              >
                Full Name
              </label>
              <Input
                id="full-name"
                type="text"
                value={fullName}
                placeholder="As it appears on your team roster"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setFullName(event.target.value);
                }}
              />
            </div>

            <div className="space-y-3">
              <label
                htmlFor="institution-email"
                className="block text-sm font-medium text-black"
              >
                Institution Email
              </label>
              <Input
                id="institution-email"
                type="email"
                value={institutionEmail}
                placeholder="name@school.edu"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setInstitutionEmail(event.target.value);
                }}
              />
            </div>

            <div className="space-y-3">
              <label
                htmlFor="school"
                className="block text-sm font-medium text-black"
              >
                School
              </label>
              <Input
                id="school"
                type="text"
                value={school}
                placeholder="Enter your school"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setSchool(event.target.value);
                }}
              />
            </div>
          </>
        ) : (
          <>
            <div className="space-y-3">
              <label
                htmlFor="name"
                className="block text-sm font-medium text-black"
              >
                Name
              </label>
              <Input
                id="name"
                type="text"
                value={name}
                placeholder="Enter your full name"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setName(event.target.value);
                }}
              />
            </div>

            <div className="space-y-3">
              <label
                htmlFor="email"
                className="block text-sm font-medium text-black"
              >
                Email
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                placeholder="name@email.com"
                className="h-10 rounded-xl px-3 text-sm"
                onChange={(event) => {
                  setEmail(event.target.value);
                }}
              />
            </div>
          </>
        )}

        <Button
          type="button"
          className="h-10 w-full rounded-xl text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            updateSection("verification", {
              name,
              email,
              fullName,
              institutionEmail,
              school,
            });
            router.push(
              `/onboarding/verification/code?role=${encodeURIComponent(role)}`,
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
