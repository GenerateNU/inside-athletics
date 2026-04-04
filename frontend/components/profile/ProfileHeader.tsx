"use client";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import type { ProfilePageData } from "@/components/profile/types";

type Props = {
  user: ProfilePageData["user"];
  isAthlete: boolean;
  onEdit: () => void;
};

export function ProfileHeader({ user, isAthlete, onEdit }: Props) {
  return (
    <>
      <div className="mb-5 flex items-start justify-between gap-4">
        <div className="flex min-w-0 flex-1 gap-5">
          <Avatar className="h-[120px] w-[120px] border-slate-300 bg-slate-300 text-slate-500">
            <AvatarFallback />
          </Avatar>
          <div className="min-w-0">
            <h1 className="text-5xl font-black tracking-tight text-[#0f2f58]">
              @{user.username || "username"}
            </h1>
            <p className="mt-2 text-[34px] font-semibold leading-none">
              {user.firstName} {user.lastName}
            </p>
            <p className="text-2xl text-slate-500">{user.pronouns}</p>
            {isAthlete ? (
              <div className="mt-3 flex flex-wrap gap-2">
                {user.divisionTag ? (
                  <Badge className="rounded-md bg-[#067b78] px-3 py-1 text-xs text-white">
                    {user.divisionTag}
                  </Badge>
                ) : null}
                {user.sportTag ? (
                  <Badge className="rounded-md bg-[#0f965b] px-3 py-1 text-xs text-white">
                    {user.sportTag}
                  </Badge>
                ) : null}
                {user.collegeTag ? (
                  <Badge className="rounded-md bg-[#19558f] px-3 py-1 text-xs text-white">
                    {user.collegeTag}
                  </Badge>
                ) : null}
              </div>
            ) : null}
            {isAthlete && user.email ? (
              <p className="mt-2 text-xs text-slate-600 underline">
                {user.email}
              </p>
            ) : null}
          </div>
        </div>
        <Button
          className="rounded-lg bg-[#2d6ca6] px-4 py-2 text-white hover:bg-[#235a8a]"
          onClick={onEdit}
        >
          Edit profile
        </Button>
      </div>

      <section>
        <h2 className="text-sm font-bold">About</h2>
        <p className="max-w-[70ch] text-sm leading-6 text-slate-700">
          {user.about}
        </p>
      </section>
    </>
  );
}
