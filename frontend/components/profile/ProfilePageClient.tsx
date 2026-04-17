"use client";

import { useQueries } from "@tanstack/react-query";
import * as React from "react";

import { EditProfileModal } from "@/components/profile/EditProfileModal";
import { ProfileFeed } from "@/components/profile/ProfileFeed";
import { ProfileHeader } from "@/components/profile/ProfileHeader";
import { ProfileSidebar } from "@/components/profile/ProfileSidebar";
import type { FeedView } from "@/components/profile/types";
import { Navbar } from "@/components/ui/navbar";
import Loading from "@/components/ui/loading";
import {
  getApiV1TagByIdQueryOptions,
  listApiV1PostByPostIdCommentsQueryOptions,
  usePatchApiV1User,
  usePostApiV1UserTag,
  useGetApiV1PostsByAuthorByAuthorId,
  useGetApiV1UserById,
  useGetApiV1UserCurrent,
  useGetApiV1UserTagFollows,
} from "@/api/hooks";
import type { CommentResponse } from "@/api/models/CommentResponse";
import { cn } from "@/lib/utils";
import { useSession } from "@/utils/SessionContext";

type Props = {
  profileUserId?: string;
};

export function ProfilePageClient({ profileUserId }: Props) {
  const session = useSession();
  const enabled = !!session?.access_token;
  const isViewingOwnProfile = !profileUserId;

  // Memoize authHeaders so they don't trigger infinite loops in useEffect
  const authHeaders = React.useMemo(() => {
    return session?.access_token
      ? { Authorization: `Bearer ${session.access_token}` }
      : undefined;
  }, [session?.access_token]);

  // Shared configuration to keep the page stable
  const commonQueryConfig = {
    enabled,
    staleTime: 1000 * 60 * 5, // 5 minutes: data won't be "stale" immediately
    refetchOnWindowFocus: false, // Prevents refreshing when you click back onto the tab
  };

  const [activeView, setActiveView] = React.useState<FeedView>("posts");
  const [showEditModal, setShowEditModal] = React.useState(false);
  const [pronouns, setPronouns] = React.useState("Pro/nouns");
  const [availableTags, setAvailableTags] = React.useState<
    Array<{ id: string; name: string }>
  >([]);

  const currentUserQuery = useGetApiV1UserCurrent({
    query: { ...commonQueryConfig, enabled: enabled && isViewingOwnProfile },
    client: { headers: authHeaders },
  });

  const viewedUserQuery = useGetApiV1UserById(profileUserId ?? "", {
    query: { ...commonQueryConfig, enabled: enabled && !isViewingOwnProfile },
    client: { headers: authHeaders },
  });

  const userQuery = isViewingOwnProfile ? currentUserQuery : viewedUserQuery;
  const user = userQuery.data;
  const userId = user?.id;

  const patchUserMutation = usePatchApiV1User({
    client: { headers: authHeaders },
  });
  const createTagFollowMutation = usePostApiV1UserTag({
    client: { headers: authHeaders },
  });

  React.useEffect(() => {
    const saved = window.localStorage.getItem("profile_pronouns");
    if (saved?.trim()) {
      setPronouns(saved);
    }
  }, []);

  // THIS KEPT LOOPING BUT NOW IT DOESNT BECAUSE authHeaders is now memoized
  React.useEffect(() => {
    if (!enabled || !authHeaders || !isViewingOwnProfile) return;
    const controller = new AbortController();

    const loadTags = async () => {
      try {
        const response = await fetch("/api/v1/tag?limit=200&offset=0", {
          method: "GET",
          headers: authHeaders,
          signal: controller.signal,
        });
        if (!response.ok) return;
        const data = (await response.json()) as {
          tags?: Array<{ id: string; name: string }>;
        };
        setAvailableTags(data.tags ?? []);
      } catch (err) {
        if (err instanceof Error && err.name !== "AbortError") {
          console.error("Failed to load tags:", err);
        }
      }
    };

    void loadTags();
    return () => controller.abort();
  }, [enabled, authHeaders, isViewingOwnProfile]);

  const postsQuery = useGetApiV1PostsByAuthorByAuthorId(
    userId ?? "",
    { limit: 20, offset: 0 },
    {
      query: { ...commonQueryConfig, enabled: enabled && !!userId },
      client: { headers: authHeaders },
    },
  );

  const tagFollowsQuery = useGetApiV1UserTagFollows({
    query: { ...commonQueryConfig, enabled: enabled && isViewingOwnProfile },
    client: { headers: authHeaders },
  });

  const tagIds = isViewingOwnProfile ? (tagFollowsQuery.data?.tag_ids ?? []) : [];

  const tagResults = useQueries({
    queries: tagIds.slice(0, 8).map((id: string) => ({
      ...getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
      ...commonQueryConfig,
    })),
  });

  const posts = postsQuery.data?.posts ?? [];
  const commentPostIds = posts.slice(0, 6).map((p) => p.id);

  const commentResults = useQueries({
    queries: commentPostIds.map((postId) => ({
      ...listApiV1PostByPostIdCommentsQueryOptions(postId, {
        headers: authHeaders,
      }),
      ...commonQueryConfig,
      enabled: enabled && !!postId,
    })),
  });

  const isLoadingProfile =
    enabled &&
    (userQuery.isLoading ||
      (!!userId && postsQuery.isLoading) ||
      (isViewingOwnProfile && tagFollowsQuery.isLoading));

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
          <div className="mx-auto flex min-w-0 flex-1 max-w-[1400px] border-x border-black/5 bg-[#eff2f5]">
            <main className="min-w-0 flex-1 border-r border-slate-300/80 px-6 py-8 md:px-10">
              <div className="mb-8 flex gap-5">
                <div className="h-[120px] w-[120px] shrink-0 rounded-full bg-gray-200 animate-pulse" />
                <div className="min-w-0 flex-1 space-y-3 pt-2">
                  <Loading lines={5} />
                </div>
              </div>
              <div className="mb-5 max-w-[70ch]">
                <Loading lines={3} />
              </div>
              <div className="mb-4 border-b border-slate-300 pb-2">
                <div className="h-5 w-48 rounded-full bg-gray-200 animate-pulse" />
              </div>
              <Loading lines={8} />
            </main>
            <aside className="hidden w-[310px] shrink-0 space-y-6 px-4 py-8 md:block">
              <div className="rounded-xl border border-[#27537f]/40 bg-white/70 p-4">
                <Loading lines={4} />
              </div>
              <div className="rounded-xl border border-[#27537f]/40 bg-white/70 p-4">
                <Loading lines={3} />
              </div>
            </aside>
          </div>
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

  // Data preparation logic
  const likedPosts = posts.filter((post) => Boolean(post.is_liked));

  const comments = commentResults
    .flatMap((r) => r.data ?? [])
    .filter(Boolean)
    .slice(0, 8) as CommentResponse[];

  const interestNames = isViewingOwnProfile
    ? tagResults
        .map((r) => r.data?.name)
        .filter((name): name is string => Boolean(name))
    : Array.from(
        new Set(
          posts
            .flatMap((post) => post.tags ?? [])
            .map((tag) => tag.name)
            .filter(Boolean),
        ),
      );

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

  const showSurveyPrompt = isViewingOwnProfile && isAthlete;

  const headerUser = {
    id: user.id,
    username: user.username,
    firstName: user.first_name,
    lastName: user.last_name,
    pronouns,
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
              showEditButton={isViewingOwnProfile}
              onEdit={
                isViewingOwnProfile ? () => setShowEditModal(true) : undefined
              }
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

      {isViewingOwnProfile ? (
        <EditProfileModal
          open={showEditModal}
          onClose={() => setShowEditModal(false)}
          isSaving={patchUserMutation.isPending}
          user={{
            firstName: user.first_name,
            lastName: user.last_name,
            pronouns,
            about: user.bio || "",
          }}
          onSave={async (values) => {
            await patchUserMutation.mutateAsync({
              data: {
                first_name: values.firstName,
                last_name: values.lastName,
                bio: values.about,
              },
            });

            const existingTagIds = new Set(tagIds);
            const selectedTagIds = new Set(values.selectedTagIds);
            const tagsToAdd = values.selectedTagIds.filter(
              (id) => !existingTagIds.has(id),
            );
            const tagsToRemove = tagIds.filter((id) => !selectedTagIds.has(id));

            await Promise.all(
              tagsToAdd.map((tagId) =>
                createTagFollowMutation.mutateAsync({ data: { tag_id: tagId } }),
              ),
            );

            await Promise.all(
              tagsToRemove.map((tagId) =>
                fetch(`/api/v1/user/tag/tag/${tagId}`, {
                  method: "DELETE",
                  headers: authHeaders,
                }),
              ),
            );

            setPronouns(values.pronouns);
            window.localStorage.setItem("profile_pronouns", values.pronouns);
            setShowEditModal(false);
            await Promise.all([userQuery.refetch(), tagFollowsQuery.refetch()]);
          }}
          availableTags={availableTags}
          selectedTagIds={tagIds}
        />
      ) : null}
    </div>
  );
}