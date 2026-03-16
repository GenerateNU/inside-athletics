# Passing Authorization Headers with Kubb Hooks

All generated Kubb hooks accept a `client` option where you can pass headers.

## Client Components

Use `useSession()` to get the token and pass it via the `client` option:

```tsx
'use client'

import { useSession } from '@/components/session-provider'
import { useDeleteApiV1CollegeById } from '@/gen/hooks/useDeleteApiV1CollegeById'

export function DeleteCollegeButton({ id }: { id: string }) {
  const session = useSession()

  const { mutate } = useDeleteApiV1CollegeById({
    client: {
      headers: { Authorization: `Bearer ${session?.access_token}` },
    },
  })

  return <button onClick={() => mutate({ id })}>Delete College</button>
}
```

## Server Side

Hooks are client-only. Use the underlying client function directly instead:

```ts
'use server'

import { deleteApiV1CollegeById } from '@/gen/clients/deleteApiV1CollegeById'
import { getServerAuthorizationHeader } from '@/lib/supabase'

export async function deleteCollegeAction(id: string) {
  const headers = await getServerAuthorizationHeader()
  return deleteApiV1CollegeById(id, { headers })
}
```