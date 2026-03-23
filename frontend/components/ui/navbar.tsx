"use client";

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
  return (
    <aside
      data-slot="navbar"
      className={cn(
        "flex h-full w-full max-w-[clamp(16rem,24vw,22rem)] shrink-0 flex-col overflow-y-auto border-r border-black/5 bg-white px-4 py-5 sm:px-5 sm:py-6",
        className,
      )}
      {...props}
    >
      <div className="flex items-center gap-3">
        <div
          aria-hidden="true"
          className="h-10 w-10 shrink-0 rounded-sm bg-zinc-300"
        />
        <span className="text-lg font-bold tracking-tight text-black">
          Inside Athletics
        </span>
      </div>

      <div className="pt-4">
        <Input
          type="search"
          placeholder="Search"
          aria-label="Search"
          className="h-10 rounded-lg border-zinc-200 bg-white px-3 text-sm text-zinc-700 placeholder:text-zinc-400"
        />
      </div>

      <Separator className="my-4 bg-zinc-200/80" />

      <nav aria-label="Primary" className="flex flex-col gap-1">
        {navItems.map(({ label, icon: Icon }) => (
          <Button
            key={label}
            variant="ghost"
            size="lg"
            className="h-11 justify-start gap-3 rounded-lg px-3 text-sm font-medium text-zinc-700 hover:bg-zinc-100 hover:text-zinc-900"
          >
            <Icon className="size-4 text-zinc-700" />
            <span>{label}</span>
          </Button>
        ))}
      </nav>

      <div className="mt-6 space-y-3">
        <div className="flex items-center gap-3 px-3">
          <Briefcase className="size-4 text-zinc-700" />
          <span className="text-sm font-medium text-zinc-800">
            Schools/Tags Following
          </span>
        </div>

        <div className="flex flex-col gap-1">
          {followingItems.map(({ label, type }) => (
            <button
              key={label}
              type="button"
              className="flex items-center gap-3 rounded-lg px-3 py-2 text-left text-sm text-zinc-700 transition-colors hover:bg-zinc-100 hover:text-zinc-900"
            >
              {type === "school" ? (
                <BookOpen className="size-4 text-zinc-700" />
              ) : (
                <span
                  aria-hidden="true"
                  className="h-1.5 w-1.5 rounded-full bg-black"
                />
              )}
              <span>{label}</span>
            </button>
          ))}
        </div>
      </div>
    </aside>
  );
}
