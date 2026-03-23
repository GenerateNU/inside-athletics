"use client";

import { useEffect, useRef, useState } from "react";
import { BookOpen, Briefcase, Home, Plus, Search } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

const navItems = [
  { label: "Home", icon: Home },
  { label: "Explore", icon: Search },
  { label: "Post", icon: Plus },
];

const followingItems = [
  { label: "Swim", type: "tag" as const },
  { label: "Airbnb", type: "tag" as const },
  { label: "Coaching", type: "tag" as const },
  { label: "Northwestern", type: "school" as const },
  { label: "Georgia Tech", type: "school" as const },
  { label: "University of Michigan", type: "school" as const },
];

type NavbarProps = React.ComponentProps<"aside">;

export function Navbar({ className, ...props }: NavbarProps) {
  const navRef = useRef<HTMLElement>(null);
  const [isCollapsed, setIsCollapsed] = useState(false);

  useEffect(() => {
    const element = navRef.current;
    if (!element) return;

    const observer = new ResizeObserver(([entry]) => {
      setIsCollapsed(entry.contentRect.width < 160);
    });

    observer.observe(element);

    return () => observer.disconnect();
  }, []);

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
              Schools/Tags Following
            </span>
          )}
        </div>

        <div className="flex flex-col gap-1">
          {followingItems.map(({ label, type }) => (
            <button
              key={label}
              type="button"
              className="flex min-w-0 items-center gap-[clamp(0.5rem,1vw,0.75rem)] rounded-lg px-[clamp(0.625rem,1vw,0.75rem)] py-[clamp(0.45rem,0.9vw,0.55rem)] text-left text-[clamp(0.78rem,1.05vw,0.88rem)] text-zinc-700 transition-colors hover:bg-zinc-100 hover:text-zinc-900"
              aria-label={label}
              title={label}
            >
              {type === "school" ? (
                <BookOpen className="size-[clamp(0.9rem,1.2vw,1rem)] shrink-0 text-zinc-700" />
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
