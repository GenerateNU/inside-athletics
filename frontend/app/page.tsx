"use client";

import { useState, useEffect, Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { useSession } from "@/utils/SessionContext";

import { SearchBar } from "@/components/post/SearchBar";
import SmallPost from "@/components/post/SmallPost";
import { Navbar } from "@/components/ui/navbar";

import { useQueries } from "@tanstack/react-query";
import {
  useGetApiV1PostsFilter,
  useGetApiV1PostsPopular,
  useGetApiV1PostsSearch,
} from "@/api/hooks";
import { CancellableTag } from "@/components/filtering/CancellableTag";
import { GetTagResponse } from "@/api";
import SearchPopup from "@/components/ui/search-popup";
import { Button } from "@/components/ui/button";
import { ChevronDown } from "lucide-react";


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
  const [showFilterPopup, setShowFilterPopup] = useState(false);

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

  const { data: allPostsData, isLoading: loadingAllPosts } = useGetApiV1PostsPopular(
    {},
    {
      query: { enabled },
      client: { headers: authHeaders }
    },
  );

  const { data: searchedPosts, isLoading: loadignSearchedPosts } = useGetApiV1PostsSearch(
    { search_str: debouncedQuery },
    {
      query: { enabled },
      client: { headers: authHeaders }
    },
  )
  console.log(debouncedQuery)
  console.log(searchedPosts)

  const { data: filteredPostsData, isLoading: loadingFilteredPosts } = useGetApiV1PostsFilter(
    { tag_ids: activeTags.map((t) => t.id).join(",") },
    {
      query: { enabled: activeTags.length > 0 },
      client: { headers: authHeaders },
    },
  );


  // post are EITHER : all posts (default), filteredPosts (if active tags), or searchedPosts (if query is searched)
  const posts = activeTags.length > 0
    ? (filteredPostsData?.posts ?? [])
    : debouncedQuery !== ""
      ? (searchedPosts?.posts ?? [])
      : (allPostsData?.posts ?? []);
  const isLoading = activeTags.length > 0 ? loadingFilteredPosts : loadingAllPosts;


  return (
    <div className="min-h-screen bg-zinc-50">

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

      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
          <div className="flex w-full max-w-5xl flex-col gap-6">
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
                <p className="text-sm text-zinc-400">Loading posts...</p>
              ) : posts.length > 0 ? (
                posts.map((post) => (
                  <SmallPost key={post.id} post={post} />
                ))
              ) : (
                <p className="text-sm text-zinc-400">No posts found.</p>
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
