"use client";
import { signup } from "@/actions/auth";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { redirect, useRouter } from "next/navigation";
import { useActionState, useEffect } from "react";
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
  const [state, signupAction] = useActionState(signup, initialState);
  const status = useFormStatus();
  const router = useRouter();

  useEffect(() => {
    if (state.success) {
      redirect("/");
    }
  }, [state]);

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
        <form className="space-y-6 rounded-md bg-white p-8 shadow-sm">
          <div className="flex w-full flex-col space-y-4">
            <Input
              id="name"
              name="name"
              type="name"
              placeholder="Name"
              required
            />
            <Input
              id="username"
              name="username"
              type="username"
              placeholder="Username"
              required
            />
            <Input
              id="password"
              name="password"
              type="password"
              placeholder="Password"
              required
            />
            {!state?.success && (
              <p className="text-red-500 text-sm"> {state.message}</p>
            )}
          </div>

          <div className="flex w-full flex-col items-center gap-2">
            <Button
              type="button"
              variant="default"
              className="h-10 w-full text-sm font-semibold"
              style={{ backgroundColor: "#2C649A", color: "#FFFFFF" }}
              onClick={() => {
                router.push("/onboarding/role");
              }}
              disabled={status.pending}
            >
              Continue
            </Button>
            <p> Forgot Password?</p>
          </div>
        </form>
      </div>
    </div>
  );
}
