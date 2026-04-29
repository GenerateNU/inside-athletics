"use client";

import Link from "next/link";
import { use, useState, useEffect, useRef } from "react";
import { useSession } from "@/utils/SessionContext";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetApiV1CollegeById,
  useGetApiV1CollegesSearch,
  useGetApiV1UserCollegeFollows,
  usePostApiV1UserCollege,
  useDeleteApiV1UserCollegeById,
  useGetApiV1PostsFilter,
  getApiV1UserCollegeFollowsQueryKey,
} from "@/api/hooks";
import { Navbar } from "@/components/ui/navbar";
import { SearchBar } from "@/components/post/SearchBar";
import SmallPost from "@/components/post/SmallPost";
import { RatingPanel } from "@/components/ui/rating-panel";
import Image from "next/image";
import { Plus, Check } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import type { PostResponse } from "@/api/models/PostResponse";

const PAGE_SIZE = 20;


export default function CollegePage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const [search, setSearch] = useState("");
  const [description, setDescription] = useState<string | null>(null);

  // Posts infinite scroll
  const [offset, setOffset] = useState(0);
  const [accPosts, setAccPosts] = useState<PostResponse[]>([]);
  const [total, setTotal] = useState(0);
  const sentinelRef = useRef<HTMLDivElement>(null);

  const queryClient = useQueryClient();

  const { data: college, isLoading, error } = useGetApiV1CollegeById(id, {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { data: collegeFollows } = useGetApiV1UserCollegeFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const isFollowing = collegeFollows?.college_ids?.includes(id) ?? false;

  const { mutate: followCollege } = usePostApiV1UserCollege({
    client: { headers: authHeaders },
  });

  const { mutate: unfollowCollege } = useDeleteApiV1UserCollegeById({
    client: { headers: authHeaders },
  });

  const { data: postsData, isFetching: fetchingPosts } = useGetApiV1PostsFilter(
    { college_ids: id, limit: PAGE_SIZE, offset },
    {
      query: { enabled },
      client: { headers: authHeaders },
    }
  );

  // Reset posts when college changes
  useEffect(() => {
    setOffset(0);
    setAccPosts([]);
    setTotal(0);
  }, [id]);

  // Accumulate pages
  useEffect(() => {
    if (!postsData?.posts) return;
    setTotal(postsData.total);
    setAccPosts((prev) => {
      if (offset === 0) return postsData.posts ?? [];
      const existingIds = new Set(prev.map((p) => p.id));
      return [...prev, ...(postsData.posts ?? []).filter((p) => !existingIds.has(p.id))];
    });
  }, [postsData, offset]);

  // Infinite scroll sentinel
  useEffect(() => {
    const sentinel = sentinelRef.current;
    if (!sentinel) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !fetchingPosts && accPosts.length < total) {
          setOffset((prev) => prev + PAGE_SIZE);
        }
      },
      { rootMargin: "200px" }
    );
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [fetchingPosts, accPosts.length, total]);

  function handleFollowToggle() {
    const key = getApiV1UserCollegeFollowsQueryKey();
    const previous = queryClient.getQueryData(key);

    if (isFollowing) {
      queryClient.setQueryData<typeof collegeFollows>(key, (old) =>
        old ? { ...old, college_ids: (old.college_ids ?? []).filter((cid) => cid !== id) } : old
      );
      unfollowCollege(
        { id },
        {
          onSuccess: () => queryClient.invalidateQueries({ queryKey: key }),
          onError: () => queryClient.setQueryData(key, previous),
        }
      );
    } else {
      queryClient.setQueryData<typeof collegeFollows>(key, (old) =>
        old ? { ...old, college_ids: [...(old.college_ids ?? []), id] } : old
      );
      followCollege(
        { data: { college_id: id } },
        {
          onSuccess: () => queryClient.invalidateQueries({ queryKey: key }),
          onError: () => queryClient.setQueryData(key, previous),
        }
      );
    }
  }

  useEffect(() => {
    if (!college?.name) return;
    const slug = college.name.replace(/\s+/g, "_");
    fetch(`https://en.wikipedia.org/api/rest_v1/page/summary/${slug}`)
      .then((res) => res.json())
      .then((data) => setDescription(data.extract ?? null))
      .catch(() => setDescription(null));
  }, [college?.name]);

  const { data: searchResults } = useGetApiV1CollegesSearch(
    { search_str: search },
    {
      query: { enabled: enabled && search.length > 0 },
      client: { headers: authHeaders },
    }
  );

  if (isLoading) return <div className="flex items-center justify-center min-h-screen"><Spinner className="size-6" /></div>;
  if (error) return <div>Error loading college</div>;

  return (
    <div className="flex min-h-screen bg-white bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60">
      <Navbar className="sticky top-0 h-screen shrink-0" />

      <div className="flex min-w-0 w-full pt-10 px-10 flex-col gap-4 m-10 rounded-4xl">
        {/* Search */}
        <div className="relative">
          <SearchBar value={search} onChange={setSearch} placeholder="Search colleges..." />
          {search.length > 0 && searchResults && (
            <ul className="absolute z-50 top-full left-0 right-0 mt-1 bg-white rounded-2xl shadow-lg border border-zinc-100 max-h-64 overflow-y-auto">
              {searchResults.results?.map((c) => (
                <li key={c.id} className="border-b border-zinc-100 last:border-0">
                  <Link href={`/college/${c.id}`} className="block py-2 px-4 text-black hover:bg-zinc-50">
                    {c.name}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Main college module */}
        <main className="flex min-w-0 w-full p-10 flex-col bg-white rounded-4xl">
          {/* Banner + logo */}
          <div className="relative w-full rounded-2xl overflow-visible mb-6">
            <div className="w-full h-20 rounded-2xl bg-linear-to-r from-blue-500 to-teal-400" />
            <div className="absolute -bottom-8 left-6">
              {college?.logo ? (
                <Image
                  src={college.logo}
                  alt={college.name}
                  width={80}
                  height={80}
                  className="rounded-md shadow-lg border-2 p-2 border-white bg-gray-100 object-contain"
                />
              ) : (
                <div className="w-20 h-20 rounded-md shadow-lg border-2 border-white bg-gray-100 flex items-center justify-center text-zinc-400 text-xs font-medium">
                  No logo
                </div>
              )}
            </div>
          </div>

          {/* Name + follow + metadata */}
          <div className="mt-5 flex items-start justify-between">
            <div className="flex-1">
              <h1 className="text-xl font-semibold text-black">{college?.name}</h1>
              <div className="flex gap-4 mt-1 text-xs text-zinc-400">
                {college?.city && college?.state && (
                  <span>{college.city}, {college.state}</span>
                )}
                {college?.website && (
                  <a href={college.website} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                    Website
                  </a>
                )}
              </div>
              <p className="mt-2 text-sm text-zinc-500 leading-relaxed">
                {description ?? "No description available."}
              </p>
            </div>

            {/* Follow button */}
            <button
              onClick={handleFollowToggle}
              className={`ml-6 px-4 py-1.5 cursor-pointer rounded-full text-sm font-medium transition-colors flex items-center gap-1 ${
                isFollowing
                  ? "bg-zinc-100 text-zinc-600 hover:bg-zinc-200"
                  : "bg-blue-500 text-white hover:bg-blue-600"
              }`}
            >
              {isFollowing ? (
                <><Check className="size-4" /> Following</>
              ) : (
                <><Plus className="size-4" /> Follow</>
              )}
            </button>
          </div>
        </main>

        {/* Ratings */}
        <h2 className="text-lg font-semibold text-black">Rating</h2>
        <RatingPanel collegeId={id} />

        {/* Posts */}
        <div className="flex flex-col gap-4 pb-10">
          <h2 className="text-lg font-semibold text-black">Posts</h2>
          {accPosts.length === 0 && !fetchingPosts ? (
            <p className="text-sm text-zinc-400">No posts yet.</p>
          ) : (
            accPosts.map((post) => (
              <SmallPost key={post.id} post={post} />
            ))
          )}
          {fetchingPosts && (
            <div className="flex justify-center py-4">
              <Spinner className="size-5" />
            </div>
          )}
          <div ref={sentinelRef} className="h-1" />
        </div>
      </div>
    </div>
  );
}
