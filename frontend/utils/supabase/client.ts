import { AppState, Platform } from 'react-native'
import { createClient, processLock } from '@supabase/supabase-js'
import { secureStorage } from '../secure-storage'

const supabaseUrl = process.env.SUPABASE_URL
const supabasePublishableKey = process.env.SUPABASE_PUBLISHABLE_KEY

/**
 * Creates a supabase client that writes tokens to keychain and can be used
 * to handle any authorization requests
 */
export const supabase = createClient(supabaseUrl, supabasePublishableKey, {
  auth: {
    ...(Platform.OS !== "web" ? { storage: secureStorage } : {}),
    autoRefreshToken: true,
    persistSession: true,
    detectSessionInUrl: false,
    lock: processLock,
  },
})
// Tells Supabase Auth to continuously refresh the session automatically
// if the app is in the foreground. When this is added, you will continue
// to receive `onAuthStateChange` events with the `TOKEN_REFRESHED` or
// `SIGNED_OUT` event if the user's session is terminated. This should
// only be registered once.
if (Platform.OS !== "web") {
  AppState.addEventListener('change', (state) => {
    if (state === 'active') {
      supabase.auth.startAutoRefresh()
    } else {
      supabase.auth.stopAutoRefresh()
    }
  })
}


