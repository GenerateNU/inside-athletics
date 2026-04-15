"use client";
import { login } from "@/actions/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useRouter } from "next/navigation";
import { useActionState } from "react";
import { useFormStatus } from "react-dom";

type loginInitialState = {
  success: boolean;
  message: string;
};

const initialState: loginInitialState = {
  success: false,
  message: "",
};

export default function LoginPage() {
  const router = useRouter();
  const [state, loginAction] = useActionState(login, initialState);
  const status = useFormStatus();

  return (
    <div className="min-h-screen flex items-center justify-center bg-stone">
      <div className="max-w-lg w-full space-y-8">
        <div className="flex justify-center">
          <label className="block text-4xl text-black font-bold">Log In</label>
        </div>
        <form className="mt-8 space-y-6 bg-white p-8">
          <div className="w-full flex flex-col items-center space-y-4">
            <p className="text-center text-sm text-gray-600">
              Enter your email and we&apos;ll send you an 8-digit login code.
            </p>
            <Input
              id="email"
              name="email"
              type="email"
              placeholder="Email"
              required
            />
            {!state?.success && (
              <p className="text-red-500 text-sm"> {state.message}</p>
            )}
          </div>

          <div className="w-full flex flex-col gap-2 items-center">
            <Button
              formAction={loginAction}
              variant="secondary"
              disabled={status.pending}
              type="submit"
            >
              SEND CODE
            </Button>
            <Button
              type="button"
              variant="default"
              onClick={() => router.push("/signup")}
              disabled={status.pending}
            >
              Sign Up
            </Button>
            <p> Forgot Password?</p>
          </div>
        </form>
      </div>
    </div>
  );
}
