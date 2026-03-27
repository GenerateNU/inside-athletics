"use client";

import { useRef, useState, useEffect } from "react";
import { BookOpen, Briefcase, Home, Plus, Search } from "lucide-react";
import { useQueries } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { useSession } from "@/utils/SessionContext";

// Generated hooks from Kubb
import {
  getApiV1CollegeByIdQueryOptions,
  getApiV1SportByIdQueryOptions,
  getApiV1TagByIdQueryOptions,
  useGetApiV1UserCollegeByUserIdFollows,
  useGetApiV1UserSportByUserIdFollows,
  useGetApiV1UserTagByUserIdFollows,
} from "@/api/hooks";

const navItems = [
  { label: "Home", icon: Home },
  { label: "Explore", icon: Search },
  { label: "Post", icon: Plus },
];

type NavbarProps = React.ComponentProps<"aside">;

export function Navbar({ className, ...props }: NavbarProps) {
  const navRef = useRef<HTMLElement>(null);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const session = useSession();
  const userId = session?.user?.id;
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  // Resize observer — unchanged
  useEffect(() => {
    const el = navRef.current;
    if (!el) return;
    const observer = new ResizeObserver(([entry]) => {
      setIsCollapsed(entry.contentRect.width < 160);
    });
    observer.observe(el);
    return () => observer.disconnect();
  }, []);

  // Step 1: Fetch the followed IDs for all three types in parallel
  const { data: tagFollows } = useGetApiV1UserTagByUserIdFollows(userId ?? "", {
    query: { enabled: enabled && !!userId },
    client: { headers: authHeaders },
  });
  const { data: sportFollows } = useGetApiV1UserSportByUserIdFollows(
    userId ?? "",
    {
      query: { enabled: enabled && !!userId },
      client: { headers: authHeaders },
    },
  );
  const { data: collegeFollows } = useGetApiV1UserCollegeByUserIdFollows(
    userId ?? "",
    {
      query: { enabled: enabled && !!userId },
      client: { headers: authHeaders },
    },
  );

  const tagIds = tagFollows?.tag_ids ?? [];
  const sportIds = sportFollows?.sport_ids ?? [];
  const collegeIds = collegeFollows?.college_ids ?? [];

  // Step 2: Fetch each individual item using useQueries (parallel, no waterfalls)
  const tagResults = useQueries({
    queries: tagIds.map((id) =>
      getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  const sportResults = useQueries({
    queries: sportIds.map((id) =>
      getApiV1SportByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  const collegeResults = useQueries({
    queries: collegeIds.map((id) =>
      getApiV1CollegeByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  // Step 3: Derive loading state and following items from query results
  const isLoadingFollowing =
    tagResults.some((r) => r.isLoading) ||
    sportResults.some((r) => r.isLoading) ||
    collegeResults.some((r) => r.isLoading);

  const followingItems = [
    ...sportResults.flatMap((r) =>
      r.data ? [{ label: r.data.name, type: "sport" as const }] : []
    ),
    ...tagResults.flatMap((r) =>
      r.data ? [{ label: r.data.name, type: "tag" as const }] : []
    ),
    ...collegeResults.flatMap((r) =>
      r.data ? [{ label: r.data.name, type: "school" as const }] : []
    ),
  ];

  return (
    <aside
      ref={navRef}
      data-slot="navbar"
      className={cn(
        "flex h-full w-[clamp(0rem,24vw,22rem)] min-w-0 max-w-[22rem] flex-col overflow-x-hidden overflow-y-auto border-r border-black/5 bg-white px-[clamp(0.5rem,1.5vw,1.25rem)] py-[clamp(0.75rem,1.75vw,1.5rem)]",
        isCollapsed && "items-center px-2 py-4",
        className,
      )}
      {...props}
    >
      {/* Logo */}
      <div className={cn("flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)]", isCollapsed && "w-full justify-center")}>
        <div aria-hidden="true" className="h-[clamp(2rem,3vw,2.5rem)] w-[clamp(2rem,3vw,2.5rem)] shrink-0 rounded-sm bg-zinc-300" />
        {!isCollapsed && (
          <span className="truncate text-[clamp(0.95rem,1.4vw,1.125rem)] font-bold tracking-tight text-black">
            Inside Athletics
          </span>
        )}
      </div>

      {/* Search */}
      <div className={cn("pt-[clamp(0.75rem,1.2vw,1rem)]", isCollapsed && "w-full")}>
        <Input
          type="search"
          placeholder={isCollapsed ? "" : "Search"}
          aria-label="Search"
          className="h-[clamp(2.25rem,3.2vw,2.5rem)] rounded-lg border-zinc-200 bg-white px-[clamp(0.625rem,1vw,0.75rem)] text-[clamp(0.8rem,1.1vw,0.9rem)] text-zinc-700 placeholder:text-zinc-400"
        />
      </div>

      <Separator className="my-[clamp(0.875rem,1.4vw,1rem)] bg-zinc-200/80" />

      {/* Nav items — unchanged */}
      <nav aria-label="Primary" className={cn("flex flex-col gap-1", isCollapsed && "w-full")}>
        {navItems.map(({ label, icon: Icon }) => (
          <Button
            key={label}
            variant="ghost"
            size="lg"
            className="h-[clamp(2.5rem,3.5vw,2.75rem)] min-w-0 justify-start gap-[clamp(0.5rem,1vw,0.75rem)] rounded-lg px-[clamp(0.625rem,1vw,0.75rem)] text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium text-zinc-700 hover:bg-zinc-100 hover:text-zinc-900"
            aria-label={label}
            title={label}
          >
            <Icon className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
            {!isCollapsed && <span className="truncate">{label}</span>}
          </Button>
        ))}
      </nav>

      {/* Following section — same JSX, driven by new data */}
      <div className={cn("mt-[clamp(1rem,2vw,1.5rem)] space-y-[clamp(0.5rem,1vw,0.75rem)]", isCollapsed && "w-full")}>
        <div className="flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]">
          <Briefcase className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
          {!isCollapsed && (
            <span className="truncate text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium text-zinc-800">
              Schools/Sports/Tags Following
            </span>
          )}
        </div>

        <div className="flex flex-col gap-1">
          {!isCollapsed && isLoadingFollowing && (
            <span className="px-[clamp(0.625rem,1vw,0.75rem)] py-[clamp(0.45rem,0.9vw,0.55rem)] text-[clamp(0.78rem,1.05vw,0.88rem)] text-zinc-500">
              Loading...
            </span>
          )}
          {!isLoadingFollowing && !followingItems.length && !isCollapsed && (
            <span className="px-[clamp(0.625rem,1vw,0.75rem)] py-[clamp(0.45rem,0.9vw,0.55rem)] text-[clamp(0.78rem,1.05vw,0.88rem)] text-zinc-500">
              No follows yet
            </span>
          )}
          {followingItems.map(({ label, type }) => (
            <button
              key={`${type}-${label}`}
              type="button"
              className="flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)] rounded-lg px-[clamp(0.625rem,1vw,0.75rem)] py-[clamp(0.45rem,0.9vw,0.55rem)] text-left text-[clamp(0.78rem,1.05vw,0.88rem)] text-zinc-700 transition-colors hover:bg-zinc-100 hover:text-zinc-900"
              aria-label={label}
              title={label}
            >
              {type === "school" ? (
                <BookOpen className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
              ) : type === "sport" ? (
                <Briefcase className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
              ) : (
                <span aria-hidden="true" className="h-[clamp(0.3rem,0.5vw,0.375rem)] w-[clamp(0.3rem,0.5vw,0.375rem)] shrink-0 rounded-full bg-black" />
              )}
              {!isCollapsed && <span className="truncate">{label}</span>}
            </button>
          ))}
        </div>
      </div>
    </aside>
  );
}
