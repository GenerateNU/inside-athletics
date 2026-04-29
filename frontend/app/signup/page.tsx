"use client";
import { signup } from "@/actions/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useActionState } from "react";
import { useFormStatus } from "react-dom";

type signupInitialState = {
  success: boolean;
  message: string;
  email?: string;
};

const initialState: signupInitialState = {
  success: false,
  message: "",
  email: "",
};

export default function SignUpPage() {
  const router = useRouter();
  const [state, signupAction] = useActionState(signup, initialState);
  const status = useFormStatus();

  return (
    <div className="flex min-h-screen items-center justify-center bg-[linear-gradient(180deg,#A8C8E8_0%,#E8F1FA_100%)] px-6 py-12">
      <div className="w-full max-w-lg space-y-6 rounded-md bg-white p-8 shadow-sm">
        <div className="space-y-2 text-center">
          <h1 className="text-4xl font-bold text-[#001F3E]">Sign Up</h1>
          <div className="flex w-full items-center justify-center px-6 py-4">
            <Image
              src="/logo_image.svg"
              alt="Inside Athletics"
              width={240}
              height={240}
              priority
              className="h-auto w-full max-w-[10rem]"
            />
          </div>
        </div>

        <form className="space-y-6">
          <div className="flex w-full flex-col space-y-4">
            <Input
              id="email"
              name="email"
              type="email"
              placeholder="Email"
              className="border-[#3E7DBB] bg-[#F0F4F8]"
              required
            />
            <Input
              id="password"
              name="password"
              type="password"
              placeholder="Password"
              className="border-[#3E7DBB] bg-[#F0F4F8]"
              required
            />
            {!state?.success && state.message ? (
              <p className="text-center text-sm text-red-600" role="alert">
                {state.message}
              </p>
            ) : null}
          </div>

          <div className="flex w-full flex-col gap-2">
            <Button
              formAction={signupAction}
              type="submit"
              disabled={status.pending}
              className="h-10 w-full rounded-xl bg-[#2C649A] text-sm font-semibold text-white"
            >
              {status.pending ? "Signing Up..." : "Sign Up"}
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={() => router.push("/login")}
              disabled={status.pending}
              className="h-10 w-full rounded-xl border-[#2C649A] text-sm font-semibold text-[#2C649A]"
            >
              Log In
            </Button>
            <button
              type="button"
              className="text-sm text-[#2C649A] underline-offset-2 hover:underline"
            >
              Forgot Password?
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
