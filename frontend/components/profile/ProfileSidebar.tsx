"use client";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

type Props = {
  showSurveyPrompt: boolean;
  communities: string[];
  interests: string[];
};

/** Figma: Communities / Interests tags */
const sidebarTagClassName =
  "h-[30px] rounded-[12px] border border-[#D4E94B] bg-[#FCFDF1] px-[8px] py-[5px] text-xs font-semibold text-slate-900 gap-[5px]";

export function ProfileSidebar({
  showSurveyPrompt,
  communities,
  interests,
}: Props) {
  return (
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
            <Badge key={community} className={sidebarTagClassName}>
              {community}
            </Badge>
          ))}
        </div>
      </section>

      <section className="rounded-xl border border-[#27537f]/60 bg-white/70 p-4">
        <h3 className="text-3xl font-black">Interests</h3>
        <div className="mt-3 flex flex-wrap gap-2">
          {interests.map((tag) => (
            <Badge key={tag} className={sidebarTagClassName}>
              {tag}
            </Badge>
          ))}
        </div>
      </section>
    </aside>
  );
}
