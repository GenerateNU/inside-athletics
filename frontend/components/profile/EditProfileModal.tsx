"use client";

import { PenSquare, UserRound, X } from "lucide-react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type Props = {
  open: boolean;
  onClose: () => void;
};

function TagSection({ title, tags }: { title: string; tags: string[] }) {
  return (
    <div>
      <h4 className="text-sm font-semibold">{title}</h4>
      <div className="mt-2 flex flex-wrap gap-2">
        {tags.map((tag) => (
          <Badge
            key={tag}
            className="rounded-md border border-slate-300 bg-slate-100 px-3 py-1 text-sm text-slate-800"
          >
            {tag}
          </Badge>
        ))}
      </div>
    </div>
  );
}

function FieldCard({
  title,
  value,
  multiline,
}: {
  title: string;
  value: string;
  multiline?: boolean;
}) {
  return (
    <div className="rounded-3xl border-2 border-slate-500/65 px-4 py-3">
      <p className="text-2xl font-semibold">{title}</p>
      <p
        className={cn(
          "text-3xl font-semibold",
          multiline && "text-base font-normal",
        )}
      >
        {value}
      </p>
    </div>
  );
}

export function EditProfileModal({ open, onClose }: Props) {
  if (!open) return null;

  return (
    <div className="fixed inset-0 z-40 bg-black/35 p-6">
      <div className="mx-auto max-h-[90vh] w-full max-w-[980px] overflow-y-auto rounded-3xl bg-white p-8 shadow-2xl md:p-10">
        <div className="mb-6 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <button
              type="button"
              className="rounded-full p-1 text-slate-600 hover:bg-slate-100"
              onClick={onClose}
            >
              <X className="h-5 w-5" />
            </button>
            <h2 className="text-3xl font-bold">Edit Profile</h2>
          </div>
          <Button
            className="rounded-lg bg-[#2d6ca6] px-4 py-2 text-white hover:bg-[#235a8a]"
            onClick={onClose}
          >
            Save
          </Button>
        </div>

        <div className="space-y-6">
          <div className="flex justify-center">
            <div className="relative">
              <Avatar className="h-24 w-24 border-slate-500 text-slate-600">
                <AvatarFallback>
                  <UserRound className="h-12 w-12" />
                </AvatarFallback>
              </Avatar>
              <button
                type="button"
                className="absolute right-0 bottom-0 rounded-full bg-slate-200 p-1"
              >
                <PenSquare className="h-4 w-4 text-slate-700" />
              </button>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <FieldCard title="Name" value="Name Lastname" />
            <FieldCard title="Pronouns" value="Pro/nouns" />
            <div className="md:col-span-2">
              <FieldCard
                title="About"
                value="I'm a high school softball player with a goal of continuing my career at the college level. I've been actively working on my skills year-round and learning what it takes to stand out in recruiting."
                multiline
              />
            </div>
          </div>

          <TagSection
            title="General"
            tags={[
              "Recruiting & NIL",
              "Finances",
              "Injury & Recovery",
              "DEI",
              "Transfer Portal",
              "Campus & Lifestyle",
            ]}
          />
          <TagSection
            title="Priorities"
            tags={[
              "Skill & Development",
              "Academics & Career",
              "Mental Health",
              "Intensity & Competition",
              "Team Dynamics",
              "Coaching Style",
            ]}
          />
        </div>
      </div>
    </div>
  );
}
