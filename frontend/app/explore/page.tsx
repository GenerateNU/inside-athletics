"use client";

// Currently the popular tags are just replaced with user tag follows!

import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { useSession } from "@/utils/SessionContext";

import { SearchBar } from "@/components/post/SearchBar";
import SmallPost from "@/components/post/SmallPost";
import { Navbar } from "@/components/ui/navbar";
import { Tag } from "@/components/post/Tag";

import { useQueries } from "@tanstack/react-query";
import {
    useGetApiV1PostsFilter,
    useGetApiV1Posts,
    useGetApiV1UserTagFollows,
    useGetApiV1UserCollegeFollows,
    useGetApiV1UserSportFollows,
    getApiV1TagByIdQueryOptions,
    getApiV1CollegeByIdQueryOptions,
    getApiV1SportByIdQueryOptions,
} from "@/api/hooks";
import { CancellableTag } from "@/components/filtering/CancellableTag";
import { GetCollegeResponse, GetTagResponse, PostResponse, SportResponse } from "@/api";
import { Spinner } from "@/components/ui/spinner";

const PAGE_SIZE = 20;

export default function ExplorePage() {
    const session = useSession();
    const enabled = !!session?.access_token;

    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const router = useRouter();
    const [query, setQuery] = useState("");
    const [activeTags, setActiveTags] = useState<GetTagResponse[]>([]);
    const [activeColleges, setActiveCollege] = useState<GetCollegeResponse[]>([]);
    const [activeSports, setActiveSports] = useState<SportResponse[]>([]);

    const [offset, setOffset] = useState(0);
    const [accPosts, setAccPosts] = useState<PostResponse[]>([]);
    const [total, setTotal] = useState(0);
    const sentinelRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (!query.trim()) return;
        const t = setTimeout(() => {
            router.push(`/?q=${encodeURIComponent(query.trim())}`);
        }, 300);
        return () => clearTimeout(t);
    }, [query, router]);

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

    const { data: tagsfollowsData } = useGetApiV1UserTagFollows(
        {
            query: { enabled },
            client: { headers: authHeaders }
        },
    );
    const { data: collegefollowsData } = useGetApiV1UserCollegeFollows(
        {
            query: { enabled },
            client: { headers: authHeaders }
        },
    );
    const { data: sportfollowsData } = useGetApiV1UserSportFollows(
        {
            query: { enabled },
            client: { headers: authHeaders }
        },
    );

    const tagIds = tagsfollowsData?.tag_ids ?? [];
    const tagQueries = useQueries({
        queries: tagIds.map((id) =>
            getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
        ),
    });
    const collegeIds = collegefollowsData?.college_ids ?? [];
    const collegeQueries = useQueries({
        queries: collegeIds.map((id) =>
            getApiV1CollegeByIdQueryOptions(id, { headers: authHeaders }),
        ),
    });
    const sportIds = sportfollowsData?.sport_ids ?? [];
    const sportQueries = useQueries({
        queries: sportIds.map((id) =>
            getApiV1SportByIdQueryOptions(id, { headers: authHeaders }),
        ),
    });

    const followedTags = tagQueries
        .map((q) => q.data)
        .filter((t) => t !== undefined);
    const followedColleges = collegeQueries
        .map((q) => q.data)
        .filter((t) => t !== undefined);
    const followedSports = sportQueries
        .map((q) => q.data)
        .filter((t) => t !== undefined);

    const hasActiveFilters = activeTags.length > 0 || activeColleges.length > 0 || activeSports.length > 0;

    const filterKey = `${hasActiveFilters ? "filtered" : "all"}|${[...activeTags].map(t => t.id).sort().join(",")}|${[...activeColleges].map(c => c.id).sort().join(",")}|${[...activeSports].map(s => s.id).sort().join(",")}`;

    // Reset pagination when filters change
    useEffect(() => { setOffset(0); setAccPosts([]); setTotal(0); }, [filterKey]);

    const { data: allPostsData, isFetching: fetchingAll } = useGetApiV1Posts(
        { limit: PAGE_SIZE, offset },
        {
            query: { enabled: enabled && !hasActiveFilters },
            client: { headers: authHeaders },
        },
    );

    const { data: filteredPostsData, isFetching: fetchingFiltered } = useGetApiV1PostsFilter(
        {
            sport_ids: activeSports.map((t) => t.id).join(","),
            tag_ids: activeTags.map((t) => t.id).join(","),
            college_ids: activeColleges.map((t) => t.id).join(","),
            limit: PAGE_SIZE,
            offset,
        },
        {
            query: { enabled: enabled && hasActiveFilters },
            client: { headers: authHeaders },
        },
    );

    const activePosts = hasActiveFilters ? filteredPostsData?.posts : allPostsData?.posts;
    const activeTotal = hasActiveFilters ? (filteredPostsData?.total ?? 0) : (allPostsData?.total ?? 0);
    const isFetching = hasActiveFilters ? fetchingFiltered : fetchingAll;
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

    return (
        <div className="min-h-screen bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60 w-full ">
            <div className="flex h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10 overflow-scroll max-h-screen">
                    <div className="flex w-full flex-col gap-6 h-full">
                        <SearchBar
                            value={query}
                            onChange={setQuery}
                            placeholder="Search posts..."
                            className="w-full"
                        />

                        <div className="flex flex-col gap-3 w-full">
                            <h2 className="font-semibold text-lg">Popular Tags</h2>
                            <div className="flex flex-wrap gap-2">
                                {followedColleges.map((college) => (
                                    <button
                                        key={college.id}
                                        onClick={() => toggleCollege(college)}
                                    >
                                        <Tag
                                            label={college.name}
                                            className={activeColleges.some((c) => c.id === college.id) ? "border-[#A8C96A] bg-[#D4E896]" : undefined}
                                        />
                                    </button>
                                ))}
                                {followedSports.map((sport) => (
                                    <button
                                        key={sport.id}
                                        onClick={() => toggleSport(sport)}
                                    >
                                        <Tag
                                            label={sport.name}
                                            className={activeSports.some((s) => s.id === sport.id) ? "border-[#A8C96A] bg-[#D4E896]" : undefined}
                                        />
                                    </button>
                                ))}
                                {followedTags.map((tag) => (
                                    <button
                                        key={tag.id}
                                        onClick={() => toggleTag(tag)}
                                    >
                                        <Tag
                                            label={tag.name}
                                            className={activeTags.some((t) => t.id === tag.id) ? "border-[#A8C96A] bg-[#D4E896]" : undefined}
                                        />
                                    </button>
                                ))}
                            </div>
                        </div>

                        <div className="flex items-center gap-2 w-full flex-wrap">
                            <span className="font-semibold text-2xl">Explore</span>
                            {activeSports.map((sport) => (
                                <CancellableTag
                                    key={sport.id}
                                    label={sport.name}
                                    onRemove={() => toggleSport(sport)}
                                />
                            ))}
                            {activeColleges.map((college) => (
                                <CancellableTag
                                    key={college.id}
                                    label={college.name}
                                    onRemove={() => toggleCollege(college)}
                                />
                            ))}
                            {activeTags.map((tag) => (
                                <CancellableTag
                                    key={tag.id}
                                    label={tag.name}
                                    onRemove={() => toggleTag(tag)}
                                />
                            ))}
                        </div>

                        <div className="flex flex-col gap-4 w-full flex-1 min-h-0">
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
