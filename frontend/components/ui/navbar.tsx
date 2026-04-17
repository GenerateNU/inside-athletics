"use client";
import Image from "next/image";
import { useState, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { BookOpen, Briefcase, Crown, Home, Plus, Search, Settings } from "lucide-react";
import { useQueries } from "@tanstack/react-query";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { useSession } from "@/utils/SessionContext";
import { useRouter } from "next/navigation";

import {
  getApiV1CollegeByIdQueryOptions,
  getApiV1SportByIdQueryOptions,
  getApiV1TagByIdQueryOptions,
  useGetApiV1UserCollegeFollows,
  useGetApiV1UserSportFollows,
  useGetApiV1UserTagFollows,
} from "@/api/hooks";
import type { GetCollegeFollowsByUserResponse } from "@/api/models/GetCollegeFollowsByUserResponse";
import type { GetCollegeResponse } from "@/api/models/GetCollegeResponse";
import type { GetSportFollowsByUserResponse } from "@/api/models/GetSportFollowsByUserResponse";
import type { GetTagFollowsByUserResponse } from "@/api/models/GetTagFollowsByUserResponse";
import type { GetTagResponse } from "@/api/models/GetTagResponse";
import type { SportResponse } from "@/api/models/SportResponse";

const navItems = [
  { label: "Home", icon: Home, href: "/" },
  { label: "Explore", icon: Search, href: "/explore" },
];

function unwrapBody<T>(value: unknown): T | undefined {
  let current = value;

  for (let depth = 0; depth < 5; depth += 1) {
    if (!current || typeof current !== "object") {
      return current as T | undefined;
    }

    if ("body" in current && current.body !== undefined) {
      current = current.body;
      continue;
    }

    if ("Body" in current && current.Body !== undefined) {
      current = current.Body;
      continue;
    }

    return current as T | undefined;
  }

  return current as T | undefined;
}

type NavbarProps = React.ComponentProps<"aside">;

export function Navbar({ className, ...props }: NavbarProps) {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const [isAdmin, setIsAdmin] = useState(false);
  const pathname = usePathname();
  const router = useRouter();
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  useEffect(() => {
    const updateCollapsed = () => {
      setIsCollapsed(window.innerWidth < 900);
    };

    updateCollapsed();
    window.addEventListener("resize", updateCollapsed);
    return () => window.removeEventListener("resize", updateCollapsed);
  }, []);

  useEffect(() => {
    if (!session?.access_token) return;
    fetch("/api/v1/role/roles", {
      headers: { Authorization: `Bearer ${session.access_token}` },
    })
      .then((r) => r.json())
      .then((data) => {
        setIsAdmin(data.roles.some((r: { name: string }) => r.name === "admin"));
      });
  }, [session?.access_token]);

  const { data: tagFollows } = useGetApiV1UserTagFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { data: sportFollows } = useGetApiV1UserSportFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { data: collegeFollows } = useGetApiV1UserCollegeFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const tagIds =
    unwrapBody<GetTagFollowsByUserResponse>(tagFollows)?.tag_ids ?? [];
  const sportIds =
    unwrapBody<GetSportFollowsByUserResponse>(sportFollows)?.sport_ids ?? [];
  const collegeIds =
    unwrapBody<GetCollegeFollowsByUserResponse>(collegeFollows)?.college_ids ?? [];

  const tagResults = useQueries({
    queries: tagIds.map((id: string) =>
      getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  const sportResults = useQueries({
    queries: sportIds.map((id: string) =>
      getApiV1SportByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  const collegeResults = useQueries({
    queries: collegeIds.map((id: string) =>
      getApiV1CollegeByIdQueryOptions(id, { headers: authHeaders }),
    ),
  });

  const isLoadingFollowing =
    tagResults.some((r) => r.isLoading) ||
    sportResults.some((r) => r.isLoading) ||
    collegeResults.some((r) => r.isLoading);

  const followingItems = [
    ...sportResults.flatMap((r) =>
      unwrapBody<SportResponse>(r.data)
        ? [{ label: unwrapBody<SportResponse>(r.data)!.name, type: "sport" as const }]
        : [],
    ),
    ...tagResults.flatMap((r) =>
      unwrapBody<GetTagResponse>(r.data)
        ? [{ label: unwrapBody<GetTagResponse>(r.data)!.name, type: "tag" as const }]
        : [],
    ),
    ...collegeResults.flatMap((r) =>
      unwrapBody<GetCollegeResponse>(r.data)
        ? [{ label: unwrapBody<GetCollegeResponse>(r.data)!.name, type: "school" as const }]
        : [],
    ),
  ];

  return (
    <aside
      data-slot="navbar"
      className={cn(
        "flex h-full max-w-[22rem] flex-col overflow-x-hidden overflow-y-auto border-r border-black/5 bg-white py-[clamp(0.75rem,1.75vw,1.5rem)] transition-[width,padding] duration-200",
        isCollapsed
          ? "w-20 min-w-20 items-center px-2"
          : "w-[clamp(16rem,24vw,22rem)] min-w-[16rem] px-[clamp(0.5rem,1.5vw,1.25rem)]",
        className,
      )}
      {...props}
    >
      <div
        className={cn(
          "flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)]",
          isCollapsed ? "w-full justify-center" : "justify-between",
        )}
      >
        <div
          className={cn(
            "flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)]",
            isCollapsed && "justify-center",
          )}
        >
          <Image
            src={"/logo_image.svg"}
            width={45}
            height={45}
            alt="Picture of the author"
          />
          {!isCollapsed && (
            <span className="truncate text-[clamp(0.95rem,1.4vw,1.125rem)] font-bold tracking-tight text-black">
              Inside Athletics
            </span>
          )}
        </div>
      </div>

      {!isCollapsed && (
        <div className="pt-[clamp(0.75rem,1.2vw,1rem)]">
          <Input
            type="search"
            placeholder="Search"
            aria-label="Search"
            className="h-[clamp(2.25rem,3.2vw,2.5rem)] rounded-lg border-zinc-200 bg-white px-[clamp(0.625rem,1vw,0.75rem)] text-[clamp(0.8rem,1.1vw,0.9rem)] text-zinc-700 placeholder:text-zinc-400"
          />
        </div>
      )}

      <Separator className="my-[clamp(0.875rem,1.4vw,1rem)] bg-zinc-200/80" />

      <nav
        aria-label="Primary"
        className={cn(
          "flex flex-col gap-1",
          isCollapsed && "w-full items-center",
        )}
      >
        {navItems.map(({ label, icon: Icon, href }) => {
          const isActive = pathname === href;
          return (
            <Button
              key={label}
              variant="ghost"
              size="lg"
              className={cn(
                "h-[clamp(2.5rem,3.5vw,2.75rem)] min-w-0 rounded-lg text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium hover:bg-zinc-100 hover:text-zinc-900",
                isActive ? "bg-zinc-100 text-zinc-900" : "text-zinc-700",
                isCollapsed
                  ? "w-12 justify-center px-0"
                  : "justify-start gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
              )}
              aria-label={label}
              aria-current={isActive ? "page" : undefined}
              title={label}
              nativeButton={false}
              render={<Link href={href} />}
            >
              <Icon
                className={cn(
                  "size-[clamp(0.9rem,1.2vw,1rem)] shrink-0",
                  isActive ? "text-zinc-900" : "text-zinc-700",
                )}
              />
              {!isCollapsed && <span className="truncate">{label}</span>}
            </Button>
          );
        })}

        <Button
          variant="ghost"
          size="lg"
          onClick={() => router.push("/?createPost=true")}
          className={cn(
            "h-[clamp(2.5rem,3.5vw,2.75rem)] min-w-0 rounded-lg text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium text-zinc-700 hover:bg-zinc-100 hover:text-zinc-900",
            isCollapsed
              ? "w-12 justify-center px-0"
              : "justify-start gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
          )}
          aria-label="Post"
          title="Post"
        >
          <Plus className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
          {!isCollapsed && <span className="truncate">Post</span>}
        </Button>
      </nav>

      <div
        className={cn(
          "mt-[clamp(1rem,2vw,1.5rem)] space-y-[clamp(0.5rem,1vw,0.75rem)]",
          isCollapsed && "w-full",
        )}
      >
        <div
          className={cn(
            "flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
            isCollapsed && "justify-center px-0",
          )}
        >
          <Briefcase className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
          {!isCollapsed && (
            <span className="truncate text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium text-zinc-800">
              Followed Tags & Schools
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
              className={cn(
                "flex min-w-0 items-center rounded-lg py-[clamp(0.45rem,0.9vw,0.55rem)] text-left text-[clamp(0.78rem,1.05vw,0.88rem)] text-zinc-700 transition-colors hover:bg-zinc-100 hover:text-zinc-900",
                isCollapsed
                  ? "w-12 justify-center px-0"
                  : "gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
              )}
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

      <Separator className="my-[clamp(0.875rem,1.4vw,1rem)] bg-zinc-200/80" />

      <Button
        variant="ghost"
        size="lg"
        className={cn(
          "h-[clamp(2.5rem,3.5vw,2.75rem)] min-w-0 rounded-lg text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium hover:bg-zinc-100 hover:text-zinc-900",
          pathname === "/insidercontaent" ? "bg-zinc-100 text-zinc-900" : "text-zinc-700",
          isCollapsed
            ? "w-12 justify-center px-0"
            : "justify-start gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
        )}
        aria-label="Insider Content"
        title="Insider Content"
        nativeButton={false}
        render={<Link href="/insidercontent" />}
      >
        <Crown className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0" />
        {!isCollapsed && <span className="truncate">Insider Content</span>}
      </Button>

      {isAdmin && (
        <Button
          variant="ghost"
          size="lg"
          className={cn(
            "h-[clamp(2.5rem,3.5vw,2.75rem)] min-w-0 rounded-lg text-[clamp(0.8rem,1.1vw,0.9rem)] font-medium hover:bg-zinc-100 hover:text-zinc-900",
            pathname === "/settings" ? "bg-zinc-100 text-zinc-900" : "text-zinc-700",
            isCollapsed
              ? "w-12 justify-center px-0"
              : "justify-start gap-[clamp(0.5rem,1vw,0.75rem)] px-[clamp(0.625rem,1vw,0.75rem)]",
          )}
          aria-label="Settings"
          title="Settings"
          nativeButton={false}
          render={<Link href="/settings" />}
        >
          <Settings className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0" />
          {!isCollapsed && <span className="truncate">Settings</span>}
        </Button>
      )}
    </aside>
  );
}