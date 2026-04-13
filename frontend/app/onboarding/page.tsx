"use client";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOnboarding } from "@/utils/onboarding";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

export default function SignUpPage() {
  const router = useRouter();
  const { data, hydrated, updateSection } = useOnboarding();
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  useEffect(() => {
    if (!hydrated) return;

    setName(data.account.name);
    setEmail(data.verification.email);
    setUsername(data.account.username);
  }, [
    data.account.name,
    data.account.username,
    data.verification.email,
    hydrated,
  ]);

  const canContinue = Boolean(
    name.trim() && email.trim() && username.trim() && password,
  );

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg rounded-[1rem] bg-white px-8 py-10 shadow-[0_18px_45px_rgba(44,100,154,0.16)]">
        <div className="space-y-6">
          <div className="space-y-4 text-center">
            <h1 className="text-4xl font-bold text-[#001F3E]">
              Join Inside Athletics
            </h1>
            <div className="flex w-full items-center justify-center rounded-md px-6 py-8">
              <Image
                src="/ia mark.svg"
                alt="Inside Athletics"
                width={240}
                height={240}
                priority
                className="h-auto w-full max-w-[14rem]"
              />
            </div>
          </div>
          <div className="space-y-6 rounded-md">
            <div className="flex w-full flex-col space-y-4">
              <Input
                id="name"
                name="name"
                type="text"
                className="border-[#3E7DBB] bg-[#F0F4F8]"
                value={name}
                placeholder="Name"
                required
                onChange={(e) => setName(e.target.value)}
              />
              <Input
                id="email"
                name="email"
                type="email"
                className="border-[#3E7DBB] bg-[#F0F4F8]"
                value={email}
                placeholder="Email"
                required
                onChange={(e) => setEmail(e.target.value)}
              />
              <Input
                id="username"
                name="username"
                type="text"
                className="border-[#3E7DBB] bg-[#F0F4F8]"
                value={username}
                placeholder="Username"
                required
                onChange={(e) => setUsername(e.target.value)}
              />
              <Input
                id="password"
                name="password"
                type="password"
                className="border-[#3E7DBB] bg-[#F0F4F8]"
                value={password}
                placeholder="Password"
                required
                onChange={(e) => setPassword(e.target.value)}
              />
            </div>

            <div className="flex w-full flex-col items-center gap-2">
              <Button
                type="button"
                variant="default"
                className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
                onClick={() => {
                  updateSection("account", { name, username });
                  updateSection("verification", { name, email });
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
    </div>
  );
}
