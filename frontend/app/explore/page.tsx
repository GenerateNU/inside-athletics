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
    <div className="flex min-h-screen bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60">
      <Navbar className="sticky top-0 h-screen shrink-0" />

      <main className="flex min-w-0 flex-1 flex-col">
        {/* Search bar */}
        <div className="sticky top-0 z-10 px-6 py-3">
          <SearchBar
            value={searchQuery}
            onChange={setSearchQuery}
            placeholder="Search"
            className="max-w-sm"
          />
        </div>

        

      </main>
    </div>
  );
}
