"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

export default function OnboardingRolePage() {
  const router = useRouter();
  const [role, setRole] = useState("");

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Choose Your Role</h1>
          <p className="text-sm text-gray-600">
            Select the role that best describes how you&apos;ll use Inside
            Athletics.
          </p>
        </div>

        <div className="space-y-3">
          <label htmlFor="role" className="block text-sm font-medium text-black">
            Role
          </label>
          <Select value={role} onValueChange={setRole}>
            <SelectTrigger id="role" className="h-10 w-full text-sm">
              <SelectValue placeholder="Select a role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="athlete">Athlete</SelectItem>
              <SelectItem value="coach">Coach</SelectItem>
              <SelectItem value="parent">Parent</SelectItem>
              <SelectItem value="recruiter">Recruiter</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <Button
          type="button"
          className="h-10 w-full text-sm font-semibold"
          style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
          onClick={() => {
            router.push("/");
          }}
          disabled={!role}
        >
          Continue
        </Button>
      </div>
    </div>
  );
}
