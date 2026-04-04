import { createServerClient } from "@supabase/ssr";
import { NextResponse, type NextRequest } from "next/server";

const BACKEND_ORIGIN =
  process.env.NEXT_PUBLIC_BACKEND_URL ?? "http://127.0.0.1:8080";

const ONBOARDING_PATH_PREFIX = "/onboarding";

async function getAppUserStatus(accessToken: string) {
  try {
    const response = await fetch(`${BACKEND_ORIGIN}/api/v1/user/current`, {
      headers: {
        Authorization: `Bearer ${accessToken}`,
      },
    });

    if (response.ok) {
      return "exists" as const;
    }

    if (response.status === 404) {
      return "missing" as const;
    }
  } catch {}

  return "unknown" as const;
}

export async function updateSession(request: NextRequest) {
  let supabaseResponse = NextResponse.next({
    request,
  });

  const supabase = createServerClient(
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_URL!
      : process.env.NEXT_PUBLIC_DEV_SUPABASE_URL!,
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY!
      : process.env.NEXT_PUBLIC_SUPABASE_DEV_PUBLISHABLE_KEY!,
    {
      cookies: {
        getAll() {
          return request.cookies.getAll();
        },
        setAll(cookiesToSet) {
          cookiesToSet.forEach(({ name, value }) =>
            request.cookies.set(name, value),
          );
          supabaseResponse = NextResponse.next({
            request: {
              headers: request.headers,
            },
          });
          cookiesToSet.forEach(({ name, value, options }) =>
            supabaseResponse.cookies.set(name, value, options),
          );
        },
      },
    },
  );
  const {
    data: { user },
  } = await supabase.auth.getUser();
  const {
    data: { session },
  } = await supabase.auth.getSession();
  const pathname = request.nextUrl.pathname;
  const isOnboardingRoute = pathname.startsWith(ONBOARDING_PATH_PREFIX);
  const isAuthRoute =
    pathname.startsWith("/login") || pathname.startsWith("/signup");

  if (
    !user &&
    !isAuthRoute &&
    !request.nextUrl.pathname.startsWith("/error")
  ) {
    const url = request.nextUrl.clone();
    url.pathname = "/login";
    return NextResponse.redirect(url);
  }

  if (user && session?.access_token) {
    const appUserStatus = await getAppUserStatus(session.access_token);

    if (appUserStatus === "missing" && !isOnboardingRoute) {
      const url = request.nextUrl.clone();
      url.pathname = "/onboarding";
      return NextResponse.redirect(url);
    }

    if (appUserStatus === "exists" && isOnboardingRoute) {
      const url = request.nextUrl.clone();
      url.pathname = "/";
      return NextResponse.redirect(url);
    }
  }

  return supabaseResponse;
}
