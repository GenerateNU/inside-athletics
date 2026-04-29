"use client";

import {
    useGetApiV1PostsPremium,
    useGetApiV1PostsPremiumSearch,
    useGetApiV1PostsPremiumFilter,
} from "@/api/hooks";
import PremiumSmallPost from "@/components/post/PremiumSmallPost";
import { Navbar } from "@/components/ui/navbar";
import { useSession, usePermissions } from "@/utils/SessionContext";
import { useState, useEffect, useRef } from "react";
import { useRouter } from "next/navigation";
import { Plus, ChevronDown } from "lucide-react";
import { Spinner } from "@/components/ui/spinner";
import CreatePremiumPostPopup from "@/components/ui/create-premium-post-popup";
import PremiumPaymentPopup from "@/components/ui/premium-payment-popup";
import { SearchBar } from "@/components/post/SearchBar";
import { CancellableTag } from "@/components/filtering/CancellableTag";
import SearchPopup from "@/components/ui/search-popup";
import { Button } from "@/components/ui/button";
import { GetCollegeResponse, GetTagResponse, PremiumPostResponse, SportResponse } from "@/api";

const PAGE_SIZE = 20;

export default function InsiderContentPage() {
    const router = useRouter();
    const session = useSession();
    const { isAdmin, hasPremium } = usePermissions();
    const enabled = !!session?.access_token && hasPremium;
    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const [showCreatePopup, setShowCreatePopup] = useState(false);
    const [showFilterPopup, setShowFilterPopup] = useState(false);
    const [searchQuery, setSearchQuery] = useState("");
    const [debouncedQuery, setDebouncedQuery] = useState("");
    const [activeTags, setActiveTags] = useState<GetTagResponse[]>([]);
    const [activeColleges, setActiveColleges] = useState<GetCollegeResponse[]>([]);
    const [activeSports, setActiveSports] = useState<SportResponse[]>([]);

    const [offset, setOffset] = useState(0);
    const [accPosts, setAccPosts] = useState<PremiumPostResponse[]>([]);
    const [total, setTotal] = useState(0);
    const sentinelRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        const t = setTimeout(() => setDebouncedQuery(searchQuery.trim()), 300);
        return () => clearTimeout(t);
    }, [searchQuery]);

    const hasActiveFilters = activeTags.length > 0 || activeColleges.length > 0 || activeSports.length > 0;
    const mode = debouncedQuery !== "" ? "search" : hasActiveFilters ? "filtered" : "all";

    const filterKey = `${mode}|${debouncedQuery}|${[...activeTags].map(t => t.id).sort().join(",")}|${[...activeColleges].map(c => c.id).sort().join(",")}|${[...activeSports].map(s => s.id).sort().join(",")}`;

    // Reset pagination when filters or mode change
    useEffect(() => { setOffset(0); setAccPosts([]); setTotal(0); }, [filterKey]);

    const { data, isError, refetch, isFetching: fetchingAll } = useGetApiV1PostsPremium(
        { limit: PAGE_SIZE, offset },
        { query: { enabled: enabled && mode === "all" }, client: { headers: authHeaders } }
    );

    const { data: searchData, isFetching: fetchingSearch } = useGetApiV1PostsPremiumSearch(
        { search_str: debouncedQuery, limit: PAGE_SIZE, offset },
        { query: { enabled: enabled && mode === "search" }, client: { headers: authHeaders } }
    );

    const { data: filteredData, isFetching: fetchingFiltered } = useGetApiV1PostsPremiumFilter(
        {
            sport_ids: activeSports.map((s) => s.id).join(","),
            tag_ids: activeTags.map((t) => t.id).join(","),
            college_ids: activeColleges.map((c) => c.id).join(","),
            limit: PAGE_SIZE,
            offset,
        },
        { query: { enabled: enabled && mode === "filtered" }, client: { headers: authHeaders } }
    );

    const activePosts = mode === "search" ? searchData?.posts
        : mode === "filtered" ? filteredData?.posts
        : data?.posts;

    const activeTotal = mode === "search" ? (searchData?.count ?? 0)
        : mode === "filtered" ? (filteredData?.total ?? 0)
        : (data?.total ?? 0);

    const isFetching = mode === "search" ? fetchingSearch
        : mode === "filtered" ? fetchingFiltered
        : fetchingAll;

    const isPostsLoading = isFetching && accPosts.length === 0;

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

    function toggleTag(tag: GetTagResponse) {
        setActiveTags((prev) => prev.some((t) => t.id === tag.id) ? prev.filter((t) => t.id !== tag.id) : [...prev, tag]);
    }
    function toggleCollege(col: GetCollegeResponse) {
        setActiveColleges((prev) => prev.some((c) => c.id === col.id) ? prev.filter((c) => c.id !== col.id) : [...prev, col]);
    }
    function toggleSport(sport: SportResponse) {
        setActiveSports((prev) => prev.some((s) => s.id === sport.id) ? prev.filter((s) => s.id !== sport.id) : [...prev, sport]);
    }

    return (
        <div className="h-screen bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60 w-full">
            {showFilterPopup && (
                <div className="fixed inset-0 z-50 flex items-center justify-center">
                    <div className="absolute inset-0 bg-black/40" onClick={() => setShowFilterPopup(false)} />
                    <div className="relative z-10">
                        <SearchPopup
                            activeTags={activeTags}
                            setActiveTagsAction={setActiveTags}
                            onBackAction={() => setShowFilterPopup(false)}
                        />
                    </div>
                </div>
            )}

            <div className="flex h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10 overflow-y-auto">
                    <div className="w-full flex flex-col gap-6">
                        <SearchBar
                            value={searchQuery}
                            onChange={setSearchQuery}
                            placeholder="Search insider content..."
                            className="w-full"
                        />

                        <div className="flex items-center justify-between">
                            <h1 className="text-4xl font-bold text-gray-900">Insider Content</h1>
                            {isAdmin && (
                                <button
                                    onClick={() => setShowCreatePopup(true)}
                                    className="flex items-center justify-center w-10 h-10 rounded-full bg-[#2C649A] text-white hover:bg-[#245580] transition-colors"
                                    aria-label="Create premium post"
                                >
                                    <Plus size={22} />
                                </button>
                            )}
                        </div>

                        <div className="flex items-center gap-2 flex-wrap">
                            {debouncedQuery ? (
                                <p className="text-sm text-zinc-500">
                                    Showing results for &ldquo;<span className="font-medium text-zinc-800">{debouncedQuery}</span>&rdquo;
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
                                        <CancellableTag key={tag.id} label={tag.name} onRemove={() => toggleTag(tag)} />
                                    ))}
                                    {activeColleges.map((col) => (
                                        <CancellableTag key={col.id} label={col.name} onRemove={() => toggleCollege(col)} />
                                    ))}
                                    {activeSports.map((sport) => (
                                        <CancellableTag key={sport.id} label={sport.name} onRemove={() => toggleSport(sport)} />
                                    ))}
                                </>
                            )}
                        </div>

                        {isError && <p className="text-red-500">Failed to load insider content.</p>}
                        <div className="flex flex-col gap-4">
                            {isPostsLoading ? (
                                <Spinner className="size-6 mx-auto text-gray-400" />
                            ) : accPosts.length === 0 ? (
                                <p className="text-gray-500">No insider content available.</p>
                            ) : (
                                accPosts.map((post) => (
                                    <PremiumSmallPost key={post.id} post={post} />
                                ))
                            )}
                            <div ref={sentinelRef} className="h-1" />
                            {isFetching && accPosts.length > 0 && (
                                <p className="text-sm text-zinc-400 text-center py-2">Loading more...</p>
                            )}
                        </div>
                    </div>
                </main>
            </div>

            {showCreatePopup && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
                    <CreatePremiumPostPopup
                        onClose={() => {
                            setShowCreatePopup(false);
                            refetch();
                        }}
                    />
                </div>
            )}

            {!hasPremium && session && (
                <PremiumPaymentPopup onClose={() => router.push("/")} />
            )}
        </div>
    );
}
