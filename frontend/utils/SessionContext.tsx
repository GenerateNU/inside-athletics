"use client";

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";
import { Session } from "@supabase/supabase-js";
import { createBrowserClient } from "@supabase/ssr";
import { getApiV1UtilityAccessCheck } from "@/api/clients/getApiV1UtilityAccessCheck";

type SessionContextValue = {
  session: Session | null;
  hasPremium: boolean;
  isAdmin: boolean;
};

const SessionContext = createContext<SessionContextValue>({
  session: null,
  hasPremium: false,
  isAdmin: false,
});

function createSupabaseClient() {
  return createBrowserClient(
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_URL!
      : process.env.NEXT_PUBLIC_DEV_SUPABASE_URL!,
    process.env.NODE_ENV === "production"
      ? process.env.NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY!
      : process.env.NEXT_PUBLIC_SUPABASE_DEV_PUBLISHABLE_KEY!,
  );
}

export function SessionProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<Session | null>(null);
  const [hasPremium, setHasPremium] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);

  useEffect(() => {
    const supabase = createSupabaseClient();
    supabase.auth.getSession().then(({ data }) => setSession(data.session));

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange((_, session) => {
      setSession(session);
    });

    return () => subscription.unsubscribe();
  }, []);

  useEffect(() => {
    if (!session?.access_token) {
      setHasPremium(false);
      setIsAdmin(false);
      return;
    }

    getApiV1UtilityAccessCheck({
      headers: { Authorization: `Bearer ${session.access_token}` },
    })
      .then((data) => {
        setHasPremium(data.has_premium ?? false);
        setIsAdmin(data.is_admin ?? false);
      })
      .catch(() => {
        setHasPremium(false);
        setIsAdmin(false);
      });
  }, [session?.access_token]);

  return (
    <SessionContext value={{ session, hasPremium, isAdmin }}>
      {children}
    </SessionContext>
  );
}

export function useSession() {
  return useContext(SessionContext).session;
}

export function usePermissions() {
  const { hasPremium, isAdmin } = useContext(SessionContext);
  return { hasPremium, isAdmin };
}
