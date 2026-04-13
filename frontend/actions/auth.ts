"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { createSupabaseServerClient } from "@/utils/auth";

type loginInitialState = {
  success: boolean;
  message: string;
};

type signupInitialState = {
  success: boolean;
  message: string;
  email?: string;
};

export async function login(prevState: loginInitialState, formData: FormData) {
  const supabase = await createSupabaseServerClient();
  const email = (formData.get("email") as string | null)?.trim() ?? "";

  if (!email) {
    return {
      success: false,
      message: "Email is required.",
    };
  }

  const { error } = await supabase.auth.signInWithOtp({
    email,
  });

  if (error) {
    return {
      success: false,
      message: error.message || "Unable to send login code.",
    };
  }

  revalidatePath("/", "layout");
  redirect(`/login/verify?email=${encodeURIComponent(email)}`);
}

export async function signup(
  prevState: signupInitialState,
  formData: FormData,
) {
  const supabase = await createSupabaseServerClient();
  const email = formData.get("email") as string;
  const payload = {
    email,
    password: formData.get("password") as string,
    options: {
      data: {},
    },
  };
  const { error } = await supabase.auth.signUp(payload);
  if (error) {
    return {
      success: false,
      message: error.message || "Login failed",
    };
  }

  revalidatePath("/", "layout");
  redirect("/onboarding");
}
