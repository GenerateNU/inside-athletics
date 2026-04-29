"use client";

import { useState, useEffect, Suspense, useRef } from "react";
import { useSearchParams } from "next/navigation";
import { useSession } from "@/utils/SessionContext";

import CreatePostPopup from "@/components/ui/create-post-popup";
import { SearchBar } from "@/components/post/SearchBar";
import SmallPost from "@/components/post/SmallPost";
import { Navbar } from "@/components/ui/navbar";

import {
  useGetApiV1PostsFilter,
  useGetApiV1PostsPopular,
  useGetApiV1PostsSearch,
} from "@/api/hooks";
import { CancellableTag } from "@/components/filtering/CancellableTag";
import { GetCollegeResponse, GetTagResponse, PostResponse, SportResponse } from "@/api";
import SearchPopup from "@/components/ui/search-popup";
import { Button } from "@/components/ui/button";
import { ChevronDown } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";

const PAGE_SIZE = 20;

function HomePageContent() {
  const searchParams = useSearchParams();
  const initialQuery = searchParams.get("q") ?? "";

  const session = useSession();
  const enabled = !!session?.access_token;

  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const [searchQuery, setSearchQuery] = useState(initialQuery);
  const [debouncedQuery, setDebouncedQuery] = useState(initialQuery);
  const [activeTags, setActiveTags] = useState<GetTagResponse[]>([]);
  const [activeColleges, setActiveCollege] = useState<GetCollegeResponse[]>([]);
  const [activeSports, setActiveSports] = useState<SportResponse[]>([]);
  const [showFilterPopup, setShowFilterPopup] = useState(false);

  const [offset, setOffset] = useState(0);
  const [accPosts, setAccPosts] = useState<PostResponse[]>([]);
  const [total, setTotal] = useState(0);
  const sentinelRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const t = setTimeout(() => setDebouncedQuery(searchQuery.trim()), 300);
    return () => clearTimeout(t);
  }, [searchQuery]);

  function toggleTag(tag: GetTagResponse) {
    setActiveTags((prev) =>
      prev.some((t) => t.id === tag.id)
        ? prev.filter((t) => t.id !== tag.id)
        : [...prev, tag],
    );
  }
  function toggleCollege(tag: GetCollegeResponse) {
    setActiveCollege((prev) =>
      prev.some((c) => c.id === tag.id)
        ? prev.filter((t) => t.id !== tag.id)
        : [...prev, tag],
    );
  }
  function toggleSport(tag: SportResponse) {
    setActiveSports((prev) =>
      prev.some((t) => t.id === tag.id)
        ? prev.filter((t) => t.id !== tag.id)
        : [...prev, tag],
    );
  }

  const hasActiveFilters = activeTags.length > 0 || activeColleges.length > 0 || activeSports.length > 0;
  const mode = hasActiveFilters ? "filtered" : debouncedQuery !== "" ? "search" : "popular";

  const filterKey = `${mode}|${debouncedQuery}|${[...activeTags].map(t => t.id).sort().join(",")}|${[...activeColleges].map(c => c.id).sort().join(",")}|${[...activeSports].map(s => s.id).sort().join(",")}`;

  // Reset pagination when filters or mode change
  useEffect(() => { setOffset(0); setAccPosts([]); setTotal(0); }, [filterKey]);

  const { data: allPostsData, isFetching: fetchingAll } = useGetApiV1PostsPopular(
    { limit: PAGE_SIZE, offset },
    {
      query: { enabled: enabled && mode === "popular" },
      client: { headers: authHeaders },
    },
  );

  const { data: searchedPosts, isFetching: fetchingSearch } = useGetApiV1PostsSearch(
    { search_str: debouncedQuery, limit: PAGE_SIZE, offset },
    {
      query: { enabled: enabled && mode === "search" },
      client: { headers: authHeaders },
    },
  );

  const { data: filteredPostsData, isFetching: fetchingFiltered } = useGetApiV1PostsFilter(
    {
      sport_ids: activeTags.filter((t) => t.type === "sports").map((t) => t.id).join(","),
      tag_ids: activeTags.filter((t) => t.type !== "sports").map((t) => t.id).join(","),
      college_ids: activeColleges.map((t) => t.id).join(","),
      limit: PAGE_SIZE,
      offset,
    },
    {
      query: { enabled: enabled && hasActiveFilters },
      client: { headers: authHeaders },
    },
  );

  const activePosts = mode === "popular" ? allPostsData?.posts
    : mode === "search" ? searchedPosts?.posts
    : filteredPostsData?.posts;

  const activeTotal = mode === "popular" ? (allPostsData?.total ?? 0)
    : mode === "search" ? (searchedPosts?.count ?? 0)
    : (filteredPostsData?.total ?? 0);

  const isFetching = mode === "popular" ? fetchingAll
    : mode === "search" ? fetchingSearch
    : fetchingFiltered;

  const isLoading = isFetching && accPosts.length === 0;

  // Accumulate posts as pages arrive
  useEffect(() => {
    if (!activePosts) return;
    setTotal(activeTotal);
    setAccPosts(prev => {
      if (offset === 0) return activePosts;
      const existingIds = new Set(prev.map(p => p.id));
      return [...prev, ...activePosts.filter(p => !existingIds.has(p.id))];
    });
  }, [activePosts, activeTotal, offset]);

  const hasMore = accPosts.length < total;

  // Trigger next page when sentinel scrolls into view
  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting && !isFetching && hasMore) {
          setOffset(prev => prev + PAGE_SIZE);
        }
      },
      { rootMargin: "200px" },
    );
    observer.observe(el);
    return () => observer.disconnect();
  }, [isFetching, hasMore]);

  const showCreatePost = searchParams.get("createPost") === "true";

  return (
    <div className="min-h-screen bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60 w-full">

      {showFilterPopup && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          <div
            className="absolute inset-0 bg-black/40"
            onClick={() => setShowFilterPopup(false)}
          />
          <div className="relative z-10">
            <SearchPopup
              activeTags={activeTags}
              setActiveTagsAction={setActiveTags}
              onBackAction={() => setShowFilterPopup(false)}
            />
          </div>
        </div>
      )}

      {showCreatePost && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <CreatePostPopup />
        </div>
      )}
      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10 overflow-scroll max-h-screen">
          <div className="flex w-full w-full flex-col gap-6">
            <SearchBar
              value={searchQuery}
              onChange={setSearchQuery}
              placeholder="Search posts..."
              className="w-full"
            />

            {/* show either the search result label OR the filter button + active tags */}
            <div className="flex items-center gap-2 w-full flex-wrap">
              {debouncedQuery ? (
                <p className="text-sm text-zinc-500">
                  Showing Search Results for &ldquo;<span className="font-medium text-zinc-800">{debouncedQuery}</span>&rdquo;
                </p>
              ) : (
                <>
                  <Button
                    variant="ghost"
                    onClick={() => setShowFilterPopup(true)}
                    className="inline-flex items-center rounded-lg border-1 border-[#D4E94B] bg-[#FCFDF1] px-3 py-1 text-xs text-zinc-800"
                  >
                    <ChevronDown size={16} />
                    Filter
                  </Button>

                  {activeTags.map((tag) => (
                    <CancellableTag
                      key={tag.id}
                      label={tag.name}
                      onRemove={() => toggleTag(tag)}
                    />
                  ))}
                </>
              )}
            </div>

            <div className="flex flex-col gap-4 w-full">
              {isLoading ? (
                <Spinner className="size-6 mx-auto text-zinc-400" />
              ) : accPosts.length > 0 ? (
                accPosts.map((post) => (
                  <SmallPost key={post.id} post={post} />
                ))
              ) : (
                <p className="text-sm text-zinc-400">No posts found.</p>
              )}
              <div ref={sentinelRef} className="h-1" />
              {isFetching && accPosts.length > 0 && (
                <Spinner className="size-5 mx-auto text-zinc-400" />
              )}
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}

export default function HomePage() {
  return (
    <Suspense>
      <HomePageContent />
    </Suspense>
  );
}
