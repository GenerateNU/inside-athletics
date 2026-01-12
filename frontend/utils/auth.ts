import { supabase } from "./supabase/client";

/**
 * Uses supabase client to get the current auth token
 * @returns the auth token
 */
export async function retrieveToken(): Promise<string> {
    const { data } = await supabase.auth.getSession();
    if (!data.session?.access_token) {
        throw new Error("Authorization token is missing.");
    }
    return data.session.access_token;
}