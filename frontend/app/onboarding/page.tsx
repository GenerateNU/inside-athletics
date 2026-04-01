"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRouter } from "next/navigation";
import { useState } from "react";

export default function SignUpPage() {
  const router = useRouter();
  const [name, setName] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const canContinue = Boolean(name.trim() && username.trim() && password);

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6">
        <div className="space-y-4 text-center">
          <h1 className="text-4xl font-bold text-black">
            Join Inside Athletics
          </h1>
          <div className="flex h-56 w-full items-center justify-center rounded-md bg-gray-200 text-sm font-medium text-gray-500">
            Placeholder Image
          </div>
        </div>
        <div className="space-y-6 rounded-md bg-white p-8 shadow-sm">
          <div className="flex w-full flex-col space-y-4">
            <Input
              id="name"
              name="name"
              type="text"
              value={name}
              placeholder="Name"
              required
              onChange={(event) => {
                setName(event.target.value);
              }}
            />
            <Input
              id="username"
              name="username"
              type="text"
              value={username}
              placeholder="Username"
              required
              onChange={(event) => {
                setUsername(event.target.value);
              }}
            />
            <Input
              id="password"
              name="password"
              type="password"
              value={password}
              placeholder="Password"
              required
              onChange={(event) => {
                setPassword(event.target.value);
              }}
            />
          </div>

          <div className="flex w-full flex-col items-center gap-2">
            <Button
              type="button"
              variant="default"
              className="h-10 w-full rounded-xl text-sm font-semibold"
              style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
              onClick={() => {
                router.push("/onboarding/role");
              }}
              disabled={!canContinue}
            >
              Continue
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
