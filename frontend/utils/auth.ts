import { createServerClient } from "@supabase/ssr";
import { cookies } from "next/headers";

export async function createSupabaseServerClient() {
  const cookieStore = await cookies();
  return createServerClient(
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_URL!
      : process.env.NEXT_PUBLIC_DEV_SUPABASE_URL!,
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY!
      : process.env.NEXT_PUBLIC_SUPABASE_DEV_PUBLISHABLE_KEY!,
    {
      cookies: {
        getAll() {
          return cookieStore.getAll();
        },
        setAll(cookiesToSet) {
          try {
            cookiesToSet.forEach(({ name, value, options }) =>
              cookieStore.set(name, value, options),
            );
          } catch {}
        },
      },
    },
  );
}

/**
 * Server-side: Creates the Authorization Header for server-side api calls
 */
export async function getServerAuthorizationHeader(): Promise<HeadersInit> {
  const supabase = await createSupabaseServerClient();
  const { data } = await supabase.auth.getSession();
  if (!data.session?.access_token) {
    throw new Error("Authorization token is missing.");
  }
  return {
    Authorization: `Bearer ${data.session?.access_token}`,
  };
}
