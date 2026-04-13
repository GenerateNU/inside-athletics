"use client";

import { useEffect, useMemo, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useOnboarding } from "@/utils/onboarding";
import { createSupabaseBrowserClient } from "@/utils/supabase/client";

const CODE_LENGTH = 6;

function normalizeDigits(value: string) {
  return value.replace(/\D/g, "").slice(0, CODE_LENGTH);
}

export default function OnboardingVerificationCodePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { data, hydrated, updateSection } = useOnboarding();
  const role = searchParams.get("role") ?? "";
  const source = searchParams.get("source") ?? "";
  const signupEmail = searchParams.get("email") ?? "";
  const inputRefs = useRef<Array<HTMLInputElement | null>>([]);
  const [digits, setDigits] = useState<string[]>(() =>
    Array.from({ length: CODE_LENGTH }, () => ""),
  );
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isResending, setIsResending] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  useEffect(() => {
    if (!hydrated) {
      return;
    }

    const normalizedCode = normalizeDigits(data.verification.code);
    const nextDigits = Array.from({ length: CODE_LENGTH }, (_, index) => {
      return normalizedCode[index] ?? "";
    });

    setDigits(nextDigits);
  }, [data.verification.code, hydrated]);

  const code = useMemo(() => digits.join(""), [digits]);
  const canContinue = code.length === CODE_LENGTH;
  const isSignupVerification = source === "signup" && Boolean(signupEmail);

  useEffect(() => {
    if (!isSignupVerification) {
      return;
    }

    let cancelled = false;

    async function resendSignupCode() {
      setIsResending(true);
      setError("");
      setNotice("");

      try {
        const supabase = createSupabaseBrowserClient();
        const { error: resendError } = await supabase.auth.resend({
          type: "signup",
          email: signupEmail,
        });

        if (cancelled) {
          return;
        }

        if (resendError) {
          setError(resendError.message || "Unable to send verification code.");
          return;
        }

        setNotice(`A 6-digit code was sent to ${signupEmail}.`);
      } finally {
        if (!cancelled) {
          setIsResending(false);
        }
      }
    }

    resendSignupCode();

    return () => {
      cancelled = true;
    };
  }, [isSignupVerification, signupEmail]);

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

  const handleContinue = async () => {
    if (!canContinue) {
      return;
    }

    updateSection("verification", {
      code,
    });

    if (!isSignupVerification) {
      router.push(
        role === "prospective-athlete"
          ? `/onboarding/teams-of-interest?role=${encodeURIComponent(role)}`
          : "/onboarding/plan",
      );
      return;
    }

    setIsSubmitting(true);
    setError("");

    try {
      const supabase = createSupabaseBrowserClient();
      const { error: verifyError } = await supabase.auth.verifyOtp({
        email: signupEmail,
        token: code,
        type: "email",
      });

      if (verifyError) {
        setError(verifyError.message || "Unable to verify your code.");
        return;
      }

      router.push(
        role
          ? `/onboarding/legal?role=${encodeURIComponent(role)}`
          : "/onboarding",
      );
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleResend = async () => {
    if (!isSignupVerification) {
      return;
    }

    setIsResending(true);
    setError("");
    setNotice("");

    try {
      const supabase = createSupabaseBrowserClient();
      const { error: resendError } = await supabase.auth.resend({
        type: "signup",
        email: signupEmail,
      });

      if (resendError) {
        setError(resendError.message || "Unable to send verification code.");
        return;
      }

      setNotice(`A new 6-digit code was sent to ${signupEmail}.`);
    } finally {
      setIsResending(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-black">Verification Code</h1>
          <p className="text-sm text-gray-600">
            Enter the verification code sent to your email.
          </p>
        </div>

        <div className="space-y-3">
          <label
            htmlFor="verification-code-0"
            className="block text-sm font-medium text-black"
          >
            Code
          </label>
          <div className="flex justify-center gap-3">
            {digits.map((digit, index) => (
              <Input
                key={index}
                id={`verification-code-${index}`}
                ref={(element) => {
                  inputRefs.current[index] = element;
                }}
                type="text"
                inputMode="numeric"
                autoComplete={index === 0 ? "one-time-code" : "off"}
                maxLength={CODE_LENGTH}
                value={digit}
                className={`h-12 w-12 rounded-xl border px-0 text-center text-lg font-semibold ${
                  digit
                    ? "border-[#7F8C2D] bg-[#D4E94B80]"
                    : "border-[#D4E94B] bg-[#FCFDF1]"
                }`}
                onChange={(event) => {
                  handleDigitChange(index, event.target.value);
                }}
                onKeyDown={(event) => {
                  handleKeyDown(index, event);
                }}
                onPaste={(event) => {
                  handlePaste(index, event);
                }}
              />
            ))}
          </div>
        </div>

        <Button
          type="button"
          className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
          onClick={handleContinue}
          disabled={!canContinue || isSubmitting}
        >
          Continue
        </Button>
        {isSignupVerification ? (
          <Button
            type="button"
            variant="outline"
            className="h-10 w-full rounded-xl border-[#2C649A] text-sm font-semibold text-[#2C649A]"
            onClick={handleResend}
            disabled={isResending || isSubmitting}
          >
            {isResending ? "Sending code..." : "Resend Code"}
          </Button>
        ) : null}
        {notice ? (
          <p className="text-sm text-green-700" role="status">
            {notice}
          </p>
        ) : null}
        {error ? (
          <p className="text-sm text-red-600" role="alert">
            {error}
          </p>
        ) : null}
      </div>
    </div>
  );
}
