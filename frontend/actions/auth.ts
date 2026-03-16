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
    const payload = {
        email: formData.get("email") as string,
        password: formData.get("password") as string,
    };
    const { error } = await supabase.auth.signInWithPassword(payload);
    console.log(error)
    if (error) {
        return {
            success: false,
            message: error.message || "Login failed",
        };
    }
    revalidatePath("/", "layout");
    redirect("/");
}

export async function signup(prevState: signupInitialState, formData: FormData) {
    const supabase = await createSupabaseServerClient();
    const payload = {
        email: formData.get("email") as string,
        password: formData.get("password") as string,
         options: {
            data: {
            }
        }
    };
    const { error } = await supabase.auth.signUp(payload);
    if (error) {
        return {
            success: false,
            message: error.message || "Login failed",
        };
    }

    return { success: true, message: "Form submitted successfully!", email: payload.email };
}