"use client";

import { useQueries } from "@tanstack/react-query";
import * as React from "react";

import { EditProfileModal } from "@/components/profile/EditProfileModal";
import { ProfileFeed } from "@/components/profile/ProfileFeed";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import { ProfileSidebar } from "@/components/profile/ProfileSidebar";
import type { FeedView } from "@/components/profile/types";
import { Navbar } from "@/components/ui/navbar";
import {
  getApiV1TagByIdQueryOptions,
  listApiV1PostByPostIdCommentsQueryOptions,
  useGetApiV1PostsByAuthorByAuthorId,
  useGetApiV1UserCurrent,
  useGetApiV1UserTagFollows,
} from "@/api/hooks";
import type { CommentResponse } from "@/api/models/CommentResponse";
import { cn } from "@/lib/utils";
import { useSession } from "@/utils/SessionContext";

export function ProfilePageClient() {
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const [activeView, setActiveView] = React.useState<FeedView>("posts");
  const [showEditModal, setShowEditModal] = React.useState(false);

  const userQuery = useGetApiV1UserCurrent({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const user = userQuery.data;
  const userId = user?.id;

  const postsQuery = useGetApiV1PostsByAuthorByAuthorId(
    userId ?? "",
    { limit: 20, offset: 0 },
    {
      query: { enabled: enabled && !!userId },
      client: { headers: authHeaders },
    },
  );

  const tagFollowsQuery = useGetApiV1UserTagFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const tagIds = tagFollowsQuery.data?.tag_ids ?? [];

  const tagResults = useQueries({
    queries: tagIds.slice(0, 8).map((id: string) => ({
      ...getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
      enabled,
    })),
  });

  const posts = postsQuery.data?.posts ?? [];

  const commentPostIds = posts.slice(0, 6).map((p) => p.id);

  const commentResults = useQueries({
    queries: commentPostIds.map((postId) => ({
      ...listApiV1PostByPostIdCommentsQueryOptions(postId, {
        headers: authHeaders,
      }),
      enabled: enabled && !!postId,
    })),
  });

  const isLoadingProfile =
    enabled &&
    (userQuery.isLoading ||
      (!!userId && postsQuery.isLoading) ||
      tagFollowsQuery.isLoading);

  if (!enabled) {
    return (
      <div className="min-h-screen bg-zinc-50">
        <div className="flex min-h-screen">
          <Navbar className="h-screen shrink-0" />
          <main className="flex min-w-0 flex-1 items-center justify-center p-6">
            <p className="text-center text-muted-foreground">
              Sign in to view your profile.
            </p>
          </main>
        </div>
      </div>
    );
  }

  if (isLoadingProfile) {
    return (
      <div className="min-h-screen bg-zinc-50">
        <div className="flex min-h-screen">
          <Navbar className="h-screen shrink-0" />
          <main className="flex min-w-0 flex-1 items-center justify-center p-6">
            <p className="text-muted-foreground">Loading profile…</p>
          </main>
        </div>
      </div>
    );
  }

  if (userQuery.isError || !user) {
    return (
      <div className="min-h-screen bg-zinc-50">
        <div className="flex min-h-screen">
          <Navbar className="h-screen shrink-0" />
          <main className="flex min-w-0 flex-1 items-center justify-center p-6">
            <p className="text-center text-muted-foreground">
              Unable to load profile data.
            </p>
          </main>
        </div>
      </div>
    );
  }

  const likedPosts = posts.filter((post) => Boolean(post.is_liked));

  const comments = commentResults
    .flatMap((r) => r.data ?? [])
    .filter(Boolean)
    .slice(0, 8) as CommentResponse[];

  const interestNames = tagResults
    .map((r) => r.data?.name)
    .filter((name): name is string => Boolean(name));

  const roleNames = (user.roles ?? []).map((role) => role.name.toLowerCase());
  const isAthlete =
    roleNames.includes("athlete") ||
    user.verified_athlete_status.toLowerCase().includes("verified");

  const communities = Array.from(
    new Set(
      [
        user.college?.name,
        ...posts.map((post) => post.college?.name).filter(Boolean),
      ].filter(Boolean),
    ),
  ) as string[];

  const showSurveyPrompt = isAthlete;

  const headerUser = {
    id: user.id,
    username: user.username,
    firstName: user.first_name,
    lastName: user.last_name,
    pronouns: "pro/nouns",
    email: isAthlete ? user.email : undefined,
    about: user.bio || "No bio yet.",
    divisionTag: user.division ? `D${user.division}` : undefined,
    sportTag: user.sport?.name,
    collegeTag: user.college?.name,
  };

  return (
    <div className="min-h-screen bg-zinc-50">
      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        <div className="mx-auto flex min-w-0 flex-1 max-w-[1400px] border-x border-black/5 bg-[#eff2f5] text-slate-900">
          <main className="min-w-0 flex-1 border-r border-slate-300/80 px-6 py-8 md:px-10">
            <ProfileHeader
              user={headerUser}
              isAthlete={isAthlete}
              onEdit={() => setShowEditModal(true)}
            />

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

            <ProfileFeed
              posts={posts}
              likedPosts={likedPosts}
              comments={comments}
              activeView={activeView}
            />
          </main>

          <ProfileSidebar
            showSurveyPrompt={showSurveyPrompt}
            communities={communities}
            interests={interestNames}
          />
        </div>
      </div>

      <EditProfileModal
        open={showEditModal}
        onClose={() => setShowEditModal(false)}
      />
    </div>
  );
}
