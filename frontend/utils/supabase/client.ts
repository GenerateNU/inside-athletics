"use client";

import { createBrowserClient } from "@supabase/ssr";

export function createSupabaseBrowserClient() {
  return createBrowserClient(
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_URL!
      : process.env.NEXT_PUBLIC_DEV_SUPABASE_URL!,
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY!
      : process.env.NEXT_PUBLIC_SUPABASE_DEV_PUBLISHABLE_KEY!,
  );
}
