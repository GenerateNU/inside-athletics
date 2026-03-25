"use client";

import { useEffect, useRef, useState } from "react";
import { BookOpen, Briefcase, Home, Plus, Search } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { useSession } from "@/utils/SessionContext";

const navItems = [
  { label: "Home", icon: Home },
  { label: "Explore", icon: Search },
  { label: "Post", icon: Plus },
];

type NavbarProps = React.ComponentProps<"aside">;
type FollowingItem = {
  label: string;
  type: "tag" | "sport" | "school";
};

type TagFollowResponse = {
  tag_ids: string[] | null;
};

type SportFollowResponse = {
  sport_ids: string[] | null;
};

type CollegeFollowResponse = {
  college_ids: string[] | null;
};

type TagResponse = {
  id: string;
  name: string;
};

type SportResponse = {
  id: string;
  name: string;
};

type CollegeResponse = {
  id: string;
  name: string;
};

const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL?.replace(/\/$/, "") ?? "";

function getApiUrl(path: string) {
  return `${apiBaseUrl}${path}`;
}

export function Navbar({ className, ...props }: NavbarProps) {
  const navRef = useRef<HTMLElement>(null);
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [followingItems, setFollowingItems] = useState<FollowingItem[]>([]);
  const [isLoadingFollowing, setIsLoadingFollowing] = useState(true);
  const session = useSession();

  useEffect(() => {
    const element = navRef.current;
    if (!element) return;

    const observer = new ResizeObserver(([entry]) => {
      setIsCollapsed(entry.contentRect.width < 160);
    });

    observer.observe(element);

    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    const accessToken = session?.access_token;

    if (!accessToken) {
      setFollowingItems([]);
      setIsLoadingFollowing(false);
      return;
    }

    const controller = new AbortController();

    async function fetchJson<T>(path: string): Promise<T> {
      const response = await fetch(getApiUrl(path), {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
        signal: controller.signal,
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch ${path}: ${response.status}`);
      }

      return response.json() as Promise<T>;
    }

    async function loadFollowing() {
      setIsLoadingFollowing(true);

      try {
        const [tagFollows, sportFollows, collegeFollows] = await Promise.all([
          fetchJson<TagFollowResponse>("/api/v1/user/tag/follows"),
          fetchJson<SportFollowResponse>("/api/v1/user/sport/follows"),
          fetchJson<CollegeFollowResponse>("/api/v1/user/college/follows"),
        ]);

        const [tags, sports, colleges] = await Promise.all([
          Promise.all(
            (tagFollows.tag_ids ?? []).map((id) =>
              fetchJson<TagResponse>(`/api/v1/tag/${id}`),
            ),
          ),
          Promise.all(
            (sportFollows.sport_ids ?? []).map((id) =>
              fetchJson<SportResponse>(`/api/v1/sport/${id}`),
            ),
          ),
          Promise.all(
            (collegeFollows.college_ids ?? []).map((id) =>
              fetchJson<CollegeResponse>(`/api/v1/college/${id}`),
            ),
          ),
        ]);

        setFollowingItems([
          ...sports.map((sport) => ({
            label: sport.name,
            type: "sport" as const,
          })),
          ...tags.map((tag) => ({
            label: tag.name,
            type: "tag" as const,
          })),
          ...colleges.map((college) => ({
            label: college.name,
            type: "school" as const,
          })),
        ]);
      } catch (error) {
        if (controller.signal.aborted) {
          return;
        }

        console.error("Unable to load followed items", error);
        setFollowingItems([]);
      } finally {
        if (!controller.signal.aborted) {
          setIsLoadingFollowing(false);
        }
      }
    }

    loadFollowing();

    return () => controller.abort();
  }, [session?.access_token]);

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
      <div
        className={cn(
          "flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)]",
          isCollapsed && "w-full justify-center",
        )}
      >
        <div
          aria-hidden="true"
          className="h-[clamp(2rem,3vw,2.5rem)] w-[clamp(2rem,3vw,2.5rem)] shrink-0 rounded-sm bg-zinc-300"
        />
        {!isCollapsed && (
          <span className="truncate text-[clamp(0.95rem,1.4vw,1.125rem)] font-bold tracking-tight text-black">
            Inside Athletics
          </span>
        )}
      </div>

      <div
        className={cn(
          "pt-[clamp(0.75rem,1.2vw,1rem)]",
          isCollapsed && "w-full",
        )}
      >
        <Input
          type="search"
          placeholder={isCollapsed ? "" : "Search"}
          aria-label="Search"
          className="h-[clamp(2.25rem,3.2vw,2.5rem)] rounded-lg border-zinc-200 bg-white px-[clamp(0.625rem,1vw,0.75rem)] text-[clamp(0.8rem,1.1vw,0.9rem)] text-zinc-700 placeholder:text-zinc-400"
        />
      </div>

      <Separator className="my-[clamp(0.875rem,1.4vw,1rem)] bg-zinc-200/80" />

      <nav
        aria-label="Primary"
        className={cn("flex flex-col gap-1", isCollapsed && "w-full")}
      >
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

      <div
        className={cn(
          "mt-[clamp(1rem,2vw,1.5rem)] space-y-[clamp(0.5rem,1vw,0.75rem)]",
          isCollapsed && "w-full",
        )}
      >
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
                <span
                  aria-hidden="true"
                  className="h-[clamp(0.3rem,0.5vw,0.375rem)] w-[clamp(0.3rem,0.5vw,0.375rem)] shrink-0 rounded-full bg-black"
                />
              )}
              {!isCollapsed && <span className="truncate">{label}</span>}
            </button>
          ))}
        </div>
      </div>
    </aside>
  );
}
