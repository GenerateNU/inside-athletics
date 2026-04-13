"use client";

import { useState, useMemo, useEffect } from "react";
import { Navbar } from "@/components/ui/navbar";
import { Badge } from "@/components/ui/badge";
import { PostCard } from "@/components/explore/PostCard";
import { SearchBar } from "@/components/explore/SearchBar";
import { useSession } from "@/utils/SessionContext";
import {
  useGetApiV1PostsPopular,
  useGetApiV1PostsSearch,
  useGetApiV1PostsFilter,
} from "@/api/hooks";
import type { Tag } from "@/api/models/Tag";
import { cn } from "@/lib/utils";

export default function ExplorePage() {
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const [searchQuery, setSearchQuery] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [selectedTag, setSelectedTag] = useState<Tag | null>(null);

  useEffect(() => {
    const timer = setTimeout(() => setDebouncedSearch(searchQuery), 400);
    return () => clearTimeout(timer);
  }, [searchQuery]);

  const isSearching = debouncedSearch.trim().length > 0;
  const isFiltering = !isSearching && !!selectedTag;


  const { data: popularData, isLoading: loadingPopular } =
    useGetApiV1PostsPopular(undefined, {
      query: { enabled: enabled && !isSearching && !isFiltering },
      client: { headers: authHeaders },
    });

  const { data: searchData, isLoading: loadingSearch } =
    useGetApiV1PostsSearch(
      { search_str: debouncedSearch },
      {
        query: { enabled: enabled && isSearching },
        client: { headers: authHeaders },
      },
    );

  const { data: filterData, isLoading: loadingFilter } =
    useGetApiV1PostsFilter(
      { tag_ids: selectedTag?.id },
      {
        query: { enabled: enabled && isFiltering },
        client: { headers: authHeaders },
      },
    );

  const allPopularPosts = popularData?.posts ?? [];

  const popularTags = useMemo<Tag[]>(() => {
    const seen = new Set<string>();
    const tags: Tag[] = [];
    for (const post of allPopularPosts) {
      for (const tag of post.tags ?? []) {
        if (!seen.has(tag.id)) {
          seen.add(tag.id);
          tags.push(tag);
        }
      }
    }
    return tags.slice(0, 8);
  }, [allPopularPosts]);

  const displayedPosts = isSearching
    ? (searchData?.posts ?? [])
    : isFiltering
      ? (filterData?.posts ?? [])
      : allPopularPosts;

  const isLoading = isSearching
    ? loadingSearch
    : isFiltering
      ? loadingFilter
      : loadingPopular;

  function handleTagClick(tag: Tag) {
    setSelectedTag((prev) => (prev?.id === tag.id ? null : tag));
  }

  return (
    <div className="flex min-h-screen bg-zinc-50">
      <Navbar className="sticky top-0 h-screen shrink-0" />

      <main className="flex min-w-0 flex-1 flex-col">
        {/* Search bar */}
        <div className="sticky top-0 z-10 border-b border-zinc-100 bg-white px-6 py-3">
          <SearchBar
            value={searchQuery}
            onChange={setSearchQuery}
            placeholder="Search"
            className="max-w-sm"
          />
        </div>

        <div className="mx-auto w-full max-w-2xl">
          {/* Popular Tags */}
          <div className="border-b border-zinc-100 bg-white px-6 py-4">
            <p className="mb-3 text-xs font-semibold uppercase tracking-wide text-zinc-400">
              Popular Tags
            </p>
            <div className="flex flex-wrap gap-2">
              {popularTags.length === 0 && !loadingPopular && (
                <span className="text-xs text-zinc-400">No tags found</span>
              )}
              {popularTags.map((tag) => (
                <button
                  key={tag.id}
                  type="button"
                  onClick={() => handleTagClick(tag)}
                >
                  <Badge
                    variant={selectedTag?.id === tag.id ? "default" : "outline"}
                    className="cursor-pointer rounded-full px-3 py-1 text-xs"
                  >
                    {tag.name}
                  </Badge>
                </button>
              ))}
            </div>
          </div>

          {/* Tabs */}
          <div className="flex gap-1 border-b border-zinc-100 bg-white px-6">
            <button
              type="button"
              onClick={() => {
                setSelectedTag(null);
                setSearchQuery("");
              }}
              className={cn(
                "pb-3 pr-1 pt-3 text-sm font-medium transition-colors",
                !selectedTag && !isSearching
                  ? "border-b-2 border-zinc-900 text-zinc-900"
                  : "text-zinc-400 hover:text-zinc-700",
              )}
            >
              Explore
            </button>
            {isSearching && (
              <button
                type="button"
                className="ml-5 border-b-2 border-zinc-900 pb-3 pt-3 text-sm font-medium text-zinc-900"
              >
                &ldquo;{debouncedSearch}&rdquo;
              </button>
            )}
            {!isSearching && selectedTag && (
              <button
                type="button"
                className="ml-5 border-b-2 border-zinc-900 pb-3 pt-3 text-sm font-medium text-zinc-900"
              >
                {selectedTag.name}
              </button>
            )}
          </div>

          {/* Post feed */}
          <div className="bg-white">
            {isLoading ? (
              <div className="flex items-center justify-center py-16 text-sm text-zinc-400">
                Loading posts...
              </div>
            ) : displayedPosts.length === 0 ? (
              <div className="flex items-center justify-center py-16 text-sm text-zinc-400">
                {isSearching
                  ? `No results for "${debouncedSearch}".`
                  : "No posts found."}
              </div>
            ) : (
              displayedPosts.map((post) => (
                <PostCard key={post.id} post={post} />
              ))
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
