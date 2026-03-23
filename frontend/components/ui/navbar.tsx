"use client";

import { Home, Plus, Search } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";

const navItems = [
  { label: "Home", icon: Home },
  { label: "Explore", icon: Search },
  { label: "Post", icon: Plus },
];

type NavbarProps = React.ComponentProps<"aside">;

export function Navbar({ className, ...props }: NavbarProps) {
  return (
    <aside
      data-slot="navbar"
      className={cn(
        "flex min-h-screen w-72 flex-col border-r border-black/5 bg-white px-5 py-6",
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
    </aside>
  );
}
