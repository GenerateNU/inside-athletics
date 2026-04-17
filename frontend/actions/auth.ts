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
  const password = (formData.get("password") as string | null) ?? "";

  if (!email || !password) {
    return {
      success: false,
      message: "Email and password are required.",
    };
  }

  const { error } = await supabase.auth.signInWithPassword({ email, password });

  if (error) {
    return {
      success: false,
      message: error.message || "Invalid email or password.",
    };
  }

  revalidatePath("/", "layout");
  redirect("/");
}

export async function signup(
  prevState: signupInitialState,
  formData: FormData,
) {
  const supabase = await createSupabaseServerClient();
  const email = (formData.get("email") as string | null)?.trim() ?? "";
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
  redirect(
    `/onboarding/verification/code?source=signup&email=${encodeURIComponent(email)}`,
  );
}
