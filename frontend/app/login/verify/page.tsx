"use client";

import { Suspense, useMemo, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { createSupabaseBrowserClient } from "@/utils/supabase/client";

const CODE_LENGTH = 8;

function normalizeDigits(value: string) {
  return value.replace(/\D/g, "").slice(0, CODE_LENGTH);
}

function LoginVerifyContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const email = searchParams.get("email")?.trim() ?? "";
  const inputRefs = useRef<Array<HTMLInputElement | null>>([]);
  const [digits, setDigits] = useState<string[]>(() =>
    Array.from({ length: CODE_LENGTH }, () => ""),
  );
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isResending, setIsResending] = useState(false);

  const code = useMemo(() => digits.join(""), [digits]);
  const canContinue = code.length === CODE_LENGTH && Boolean(email);

  const focusInput = (index: number) => {
    inputRefs.current[index]?.focus();
    inputRefs.current[index]?.select();
  };

  const handleDigitChange = (index: number, value: string) => {
    const normalizedValue = normalizeDigits(value);

    if (!normalizedValue) {
      setDigits((current) => {
        const next = [...current];
        next[index] = "";
        return next;
      });
      return;
    }

    setDigits((current) => {
      const next = [...current];

      if (normalizedValue.length > 1) {
        normalizedValue.split("").forEach((digit, offset) => {
          const targetIndex = index + offset;
          if (targetIndex < CODE_LENGTH) {
            next[targetIndex] = digit;
          }
        });
      } else {
        next[index] = normalizedValue;
      }

      return next;
    });

    const nextIndex = Math.min(index + normalizedValue.length, CODE_LENGTH - 1);
    focusInput(nextIndex);
  };

  const handleKeyDown = (
    index: number,
    event: React.KeyboardEvent<HTMLInputElement>,
  ) => {
    if (event.key === "Backspace" && !digits[index] && index > 0) {
      event.preventDefault();
      setDigits((current) => {
        const next = [...current];
        next[index - 1] = "";
        return next;
      });
      focusInput(index - 1);
    }

    if (event.key === "ArrowLeft" && index > 0) {
      event.preventDefault();
      focusInput(index - 1);
    }

    if (event.key === "ArrowRight" && index < CODE_LENGTH - 1) {
      event.preventDefault();
      focusInput(index + 1);
    }
  };

  const handlePaste = (
    index: number,
    event: React.ClipboardEvent<HTMLInputElement>,
  ) => {
    event.preventDefault();
    handleDigitChange(index, event.clipboardData.getData("text"));
  };

  const handleVerify = async () => {
    if (!canContinue) {
      return;
    }

    setIsSubmitting(true);
    setError("");
    setNotice("");

    try {
      const supabase = createSupabaseBrowserClient();
      const { error: verifyError } = await supabase.auth.verifyOtp({
        email,
        token: code,
        type: "email",
      });

      if (verifyError) {
        setError(verifyError.message || "Unable to verify your code.");
        return;
      }

      router.push("/");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleResend = async () => {
    if (!email) {
      setError("Missing email address. Start the login flow again.");
      return;
    }

    setIsResending(true);
    setError("");
    setNotice("");

    try {
      const supabase = createSupabaseBrowserClient();
      const { error: resendError } = await supabase.auth.signInWithOtp({
        email,
      });

      if (resendError) {
        setError(resendError.message || "Unable to send a new login code.");
        return;
      }

      setNotice(`A new 8-digit code was sent to ${email}.`);
    } finally {
      setIsResending(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
      <div className="w-full max-w-lg space-y-6 bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Verify Login</h1>
          <p className="text-sm text-gray-600">
            Enter the 8-digit code sent to {email || "your email"}.
          </p>
        </div>

        <div className="space-y-3">
          <label
            htmlFor="login-code-0"
            className="block text-sm font-medium text-black"
          >
            Code
          </label>
          <div className="flex justify-center gap-3">
            {digits.map((digit, index) => (
              <Input
                key={index}
                id={`login-code-${index}`}
                ref={(element) => {
                  inputRefs.current[index] = element;
                }}
                type="text"
                inputMode="numeric"
                autoComplete={index === 0 ? "one-time-code" : undefined}
                maxLength={CODE_LENGTH}
                value={digit}
                onChange={(event) => {
                  handleDigitChange(index, event.target.value);
                }}
                onKeyDown={(event) => {
                  handleKeyDown(index, event);
                }}
                onPaste={(event) => {
                  handlePaste(index, event);
                }}
                className="h-14 w-12 rounded-xl border-gray-300 text-center text-xl font-semibold"
              />
            ))}
          </div>
          {error ? (
            <p className="text-center text-sm text-red-600" role="alert">
              {error}
            </p>
          ) : null}
          {notice ? (
            <p className="text-center text-sm text-green-700" role="status">
              {notice}
            </p>
          ) : null}
        </div>

        <div className="flex flex-col gap-3">
          <Button
            type="button"
            onClick={handleVerify}
            disabled={!canContinue || isSubmitting}
            variant="secondary"
          >
            {isSubmitting ? "VERIFYING..." : "VERIFY CODE"}
          </Button>
          <Button
            type="button"
            onClick={handleResend}
            disabled={!email || isResending}
            variant="outline"
          >
            {isResending ? "SENDING..." : "RESEND CODE"}
          </Button>
          <Button
            type="button"
            onClick={() => router.push("/login")}
            disabled={isSubmitting || isResending}
            variant="ghost"
          >
            Use a different email
          </Button>
        </div>
      </div>
    </div>
  );
}

export default function LoginVerifyPage() {
  return (
    <Suspense
      fallback={
        <div className="flex min-h-screen items-center justify-center bg-stone px-6 py-12">
          <div className="w-full max-w-lg bg-white p-8 text-center shadow-sm">
            <p className="text-sm text-gray-600">Loading verification form...</p>
          </div>
        </div>
      }
    >
      <LoginVerifyContent />
    </Suspense>
  );
}
