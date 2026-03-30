"use client";

import * as React from "react";
import {
  Heart,
  Home,
  Plus,
  Search,
  Settings,
  Tag,
  X,
  Bookmark,
  PenSquare,
  UserRound,
} from "lucide-react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

type FeedView = "posts" | "comments" | "likes";
type ProfileType = "athlete" | "regular";

const onboardingTags = [
  "Recruiting & NIL",
  "Campus & Lifestyle",
  "Intensity & Competition",
];

const communities = ["Basketball", "University of Washington"];

const postRows = Array.from({ length: 5 }, (_, i) => ({
  id: `post-${i}`,
  text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
}));

const commentRows = Array.from({ length: 4 }, (_, i) => ({
  id: `comment-${i}`,
  text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
  reply: i % 2 === 0 ? "So interesting!" : "Can you explain more?",
}));

function FeedPost({
  text,
  filledHeart,
}: {
  text: string;
  filledHeart: boolean;
}) {
  return (
    <article className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
      <Avatar className="h-9 w-9 bg-slate-300">
        <AvatarFallback />
      </Avatar>
      <div>
        <div className="mb-1 flex items-baseline gap-3">
          <p className="text-sm font-semibold">username</p>
          <p className="text-sm leading-6 text-slate-700">{text}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
            <Bookmark className="mr-1 h-3 w-3" />
            Northeastern
          </Badge>
          <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
            Swim
          </Badge>
        </div>
      </div>
      <button
        type="button"
        className="mt-1 text-[#2f6aa0] transition-opacity hover:opacity-70"
      >
        <Heart className={cn("h-5 w-5", filledHeart && "fill-current")} />
      </button>
    </article>
  );
}

function FeedComment({ text, reply }: { text: string; reply: string }) {
  return (
    <article className="space-y-2">
      <div className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
        <Avatar className="h-9 w-9 bg-slate-300">
          <AvatarFallback />
        </Avatar>
        <div className="mb-1 flex items-baseline gap-3">
          <p className="text-sm font-semibold">otherperson</p>
          <p className="text-sm leading-6 text-slate-700">{text}</p>
        </div>
        <button
          type="button"
          className="mt-1 text-slate-800 transition-opacity hover:opacity-70"
        >
          <Heart className="h-5 w-5" />
        </button>
      </div>
      <div className="ml-12 border-l border-slate-300 bg-slate-100/75 px-3 py-2">
        <div className="grid grid-cols-[auto_1fr_auto] items-center gap-3">
          <Avatar className="h-7 w-7 bg-slate-300">
            <AvatarFallback />
          </Avatar>
          <div className="flex items-baseline gap-2">
            <p className="text-sm font-semibold">username</p>
            <p className="text-sm text-slate-700">{reply}</p>
          </div>
          <button
            type="button"
            className="text-slate-800 transition-opacity hover:opacity-70"
          >
            <Heart className="h-4 w-4" />
          </button>
        </div>
      </div>
    </article>
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
      <p className={cn("text-3xl font-semibold", multiline && "text-base font-normal")}>
        {value}
      </p>
    </div>
  );
}

export default function ProfilePage() {
  const [profileType, setProfileType] = React.useState<ProfileType>("athlete");
  const [activeView, setActiveView] = React.useState<FeedView>("posts");
  const [showEditModal, setShowEditModal] = React.useState(false);

  const isAthlete = profileType === "athlete";
  const showSurveyPrompt = isAthlete;

  return (
    <div className="min-h-screen bg-[#eff2f5] text-slate-900">
      <div className="mx-auto flex w-full max-w-[1400px]">
        <aside className="sticky top-0 hidden h-screen w-[72px] border-r border-slate-300/80 bg-white/65 px-3 py-4 md:flex md:flex-col md:items-center md:justify-between">
          <div className="flex w-full flex-col items-center gap-5">
            <div className="h-7 w-7 rounded bg-slate-200" />
            <Home className="h-5 w-5 text-slate-700" />
            <Search className="h-5 w-5 text-slate-500" />
            <Plus className="h-5 w-5 text-slate-500" />
            <Tag className="h-5 w-5 text-slate-500" />
          </div>
          <Settings className="h-5 w-5 text-slate-400" />
        </aside>

        <main className="flex-1 border-r border-slate-300/80 px-6 py-8 md:px-10">
          <div className="mb-5 flex items-start justify-between gap-4">
            <div className="flex min-w-0 flex-1 gap-5">
              <Avatar className="h-[120px] w-[120px] border-slate-300 bg-slate-300 text-slate-500">
                <AvatarFallback />
              </Avatar>
              <div className="min-w-0">
                <h1 className="text-5xl font-black tracking-tight text-[#0f2f58]">
                  @username
                </h1>
                <p className="mt-2 text-[34px] font-semibold leading-none">Name</p>
                <p className="text-2xl text-slate-500">pro/nouns</p>
                {isAthlete ? (
                  <div className="mt-3 flex flex-wrap gap-2">
                    <Badge className="rounded-md bg-[#067b78] px-3 py-1 text-xs text-white">
                      D1
                    </Badge>
                    <Badge className="rounded-md bg-[#0f965b] px-3 py-1 text-xs text-white">
                      Football
                    </Badge>
                    <Badge className="rounded-md bg-[#19558f] px-3 py-1 text-xs text-white">
                      UMASS Amherst
                    </Badge>
                  </div>
                ) : null}
                {isAthlete ? (
                  <p className="mt-2 text-xs text-slate-600 underline">
                    username@gmail.com
                  </p>
                ) : null}
              </div>
            </div>
            <div className="flex flex-col items-end gap-2">
              <Button
                className="rounded-lg bg-[#2d6ca6] px-4 py-2 text-white hover:bg-[#235a8a]"
                onClick={() => setShowEditModal(true)}
              >
                Edit profile
              </Button>
              <Button
                variant="outline"
                className="text-[11px]"
                onClick={() =>
                  setProfileType((value) =>
                    value === "athlete" ? "regular" : "athlete",
                  )
                }
              >
                Toggle {isAthlete ? "regular" : "athlete"} view
              </Button>
            </div>
          </div>

          <section>
            <h2 className="text-sm font-bold">About</h2>
            <p className="max-w-[70ch] text-sm leading-6 text-slate-700">
              I&apos;m currently a Division I football player at University of
              Massachusetts Amherst with first-hand experience navigating the
              college recruiting process. Coming out of a competitive high school
              program in New Jersey, I learned how to position myself, communicate
              with coaches, and find the right fit both athletically and
              academically. I&apos;m passionate about helping prospective
              student-athletes simplify that process.
            </p>
          </section>

          <div className="mt-5 border-b border-slate-300">
            <div className="flex gap-8 text-sm font-semibold text-slate-500">
              {(["posts", "comments", "likes"] as FeedView[]).map((tab) => (
                <button
                  key={tab}
                  type="button"
                  className={cn(
                    "cursor-pointer border-b-2 border-transparent pb-2 capitalize transition-colors",
                    activeView === tab && "border-slate-900 text-slate-900",
                  )}
                  onClick={() => setActiveView(tab)}
                >
                  {tab}
                </button>
              ))}
            </div>
          </div>

          <div className="mt-4 space-y-6">
            {activeView === "posts"
              ? postRows.map((row) => (
                  <FeedPost key={row.id} text={row.text} filledHeart={false} />
                ))
              : null}
            {activeView === "likes"
              ? postRows.map((row) => (
                  <FeedPost key={row.id} text={row.text} filledHeart />
                ))
              : null}
            {activeView === "comments"
              ? commentRows.map((row) => (
                  <FeedComment key={row.id} text={row.text} reply={row.reply} />
                ))
              : null}
          </div>
        </main>

        <aside className="w-[310px] space-y-4 px-4 py-8">
          {showSurveyPrompt ? (
            <section className="rounded-xl border border-[#27537f]/60 bg-white/70 p-4">
              <h3 className="text-3xl font-black">To Do</h3>
              <Button className="mt-3 rounded-lg bg-[#2d6ca6] px-4 py-2 text-white shadow-md shadow-[#2d6ca6]/35 hover:bg-[#235a8a]">
                Complete School Survey
              </Button>
            </section>
          ) : null}

          <section className="rounded-xl border border-[#27537f]/60 bg-white/70 p-4">
            <h3 className="text-3xl font-black">Communities</h3>
            <div className="mt-3 flex flex-wrap gap-2">
              {communities.map((community) => (
                <Badge
                  key={community}
                  className="rounded-xl border-none bg-gradient-to-r from-[#8eb86c] via-[#86b97a] to-[#6da9c8] px-3 py-1 text-xs font-semibold text-slate-900"
                >
                  {community}
                </Badge>
              ))}
            </div>
          </section>

          <section className="rounded-xl border border-[#27537f]/60 bg-white/70 p-4">
            <h3 className="text-3xl font-black">Interests</h3>
            <div className="mt-3 flex flex-wrap gap-2">
              {onboardingTags.map((tag) => (
                <Badge
                  key={tag}
                  className="rounded-xl border border-[#87b76f] bg-[#d9edcd] px-3 py-1 text-xs font-semibold text-slate-800"
                >
                  {tag}
                </Badge>
              ))}
            </div>
          </section>
        </aside>
      </div>

      {showEditModal ? (
        <div className="fixed inset-0 z-40 bg-black/35 p-6">
          <div className="mx-auto max-h-[90vh] w-full max-w-[980px] overflow-y-auto rounded-3xl bg-white p-8 shadow-2xl md:p-10">
            <div className="mb-6 flex items-center justify-between">
              <div className="flex items-center gap-3">
                <button
                  type="button"
                  className="rounded-full p-1 text-slate-600 hover:bg-slate-100"
                  onClick={() => setShowEditModal(false)}
                >
                  <X className="h-5 w-5" />
                </button>
                <h2 className="text-3xl font-bold">Edit Profile</h2>
              </div>
              <Button
                className="rounded-lg bg-[#2d6ca6] px-4 py-2 text-white hover:bg-[#235a8a]"
                onClick={() => setShowEditModal(false)}
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

              <div>
                <h4 className="text-sm font-semibold">General</h4>
                <div className="mt-2 flex flex-wrap gap-2">
                  {[
                    "Recruiting & NIL",
                    "Finances",
                    "Injury & Recovery",
                    "DEI",
                    "Transfer Portal",
                    "Campus & Lifestyle",
                  ].map((tag) => (
                    <Badge
                      key={tag}
                      className="rounded-md border border-slate-300 bg-slate-100 px-3 py-1 text-sm text-slate-800"
                    >
                      {tag}
                    </Badge>
                  ))}
                </div>
              </div>

              <div>
                <h4 className="text-sm font-semibold">Priorities</h4>
                <div className="mt-2 flex flex-wrap gap-2">
                  {[
                    "Skill & Development",
                    "Academics & Career",
                    "Mental Health",
                    "Intensity & Competition",
                    "Team Dynamics",
                    "Coaching Style",
                  ].map((tag) => (
                    <Badge
                      key={tag}
                      className="rounded-md border border-slate-300 bg-slate-100 px-3 py-1 text-sm text-slate-800"
                    >
                      {tag}
                    </Badge>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}