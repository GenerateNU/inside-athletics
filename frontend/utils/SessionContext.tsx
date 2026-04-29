"use client";

import {
  createContext,
  ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";
import { Session } from "@supabase/supabase-js";
import { getApiV1UtilityAccessCheck } from "@/api/clients/getApiV1UtilityAccessCheck";
import { getApiV1UserCurrent } from "@/api/clients/getApiV1UserCurrent";
import { getApiV1UserCurrentQueryKey } from "@/api/hooks";
import { createSupabaseBrowserClient } from "@/utils/supabase/client";
import { useQueryClient } from "@tanstack/react-query";
import type { GetUserResponse } from "@/api";

type SessionContextValue = {
  session: Session | null;
  hasPremium: boolean;
  isAdmin: boolean;
  currentUser: GetUserResponse | null;
};

const SessionContext = createContext<SessionContextValue>({
  session: null,
  hasPremium: false,
  isAdmin: false,
  currentUser: null,
});

export function SessionProvider({
  children,
  initialSession,
}: {
  children: ReactNode;
  initialSession: Session | null;
}) {
  const queryClient = useQueryClient();
  const [session, setSession] = useState<Session | null>(initialSession);
  const [hasPremium, setHasPremium] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const [currentUser, setCurrentUser] = useState<GetUserResponse | null>(null);

  useEffect(() => {
    const supabase = createSupabaseBrowserClient();
    // Keep session in sync with Supabase auth state changes (login, logout, token refresh)
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
      setCurrentUser(null);
      queryClient.removeQueries({ queryKey: getApiV1UserCurrentQueryKey() });
      return;
    }

    const headers = { Authorization: `Bearer ${session.access_token}` };

    // Fetch permissions and current user in parallel so both are ready immediately
    Promise.all([
      getApiV1UtilityAccessCheck({ headers }),
      getApiV1UserCurrent({ headers }),
    ])
      .then(([access, user]) => {
        setHasPremium(access.has_premium ?? false);
        setIsAdmin(access.is_admin ?? false);
        setCurrentUser(user);
        // Seed React Query's cache so useGetApiV1UserCurrent resolves instantly everywhere
        queryClient.setQueryData(getApiV1UserCurrentQueryKey(), user);
      })
      .catch(() => {
        setHasPremium(false);
        setIsAdmin(false);
      });
  }, [session?.access_token, queryClient]);

  return (
    <SessionContext value={{ session, hasPremium, isAdmin, currentUser }}>
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

export function useCurrentUser() {
  return useContext(SessionContext).currentUser;
}
