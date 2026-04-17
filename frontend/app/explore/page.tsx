"use client";

// Currently the popular tags are just replaced with user tag follows!

import { useState, useEffect } from "react";
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
import { GetCollegeResponse, GetTagResponse, SportResponse } from "@/api";


export default function ExplorePage() {
    const session = useSession();
    const enabled = !!session?.access_token;

    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const router = useRouter();
    const [query, setQuery] = useState("");
    const [activeTags, setActiveTags] = useState<GetTagResponse[]>([]);

    useEffect(() => {
        if (!query.trim()) return;
        const t = setTimeout(() => {
            router.push(`/?q=${encodeURIComponent(query.trim())}`);
        }, 300);
        return () => clearTimeout(t);
    }, [query, router]);
    const [activeColleges, setActiveCollege] = useState<GetCollegeResponse[]>([]);
    const [activeSports, setActiveSports] = useState<SportResponse[]>([]);

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


    const { data: allPostsData, isLoading: loadingAllPosts } = useGetApiV1Posts(
        { },
        {
            query: { enabled },
            client: { headers: authHeaders } },
    );
    console.log("all posts: " + allPostsData)

    const hasActiveFilters = activeTags.length > 0 || activeColleges.length > 0 || activeSports.length > 0;

    const { data: filteredPostsData, isLoading: loadingFilteredPosts } = useGetApiV1PostsFilter(
        { sport_ids: activeSports.map((t) => t.id).join(","),
            tag_ids: activeTags.map((t) => t.id).join(","),
          college_ids: activeColleges.map((t) => t.id).join(","),
        },
        {
            query: { enabled: hasActiveFilters },
            client: { headers: authHeaders },
        },
    );
    console.log("filtered: " + filteredPostsData)

    const posts = hasActiveFilters
        ? (filteredPostsData?.posts ?? [])
        : (allPostsData?.posts ?? []);
    const isLoading = hasActiveFilters ? loadingFilteredPosts : loadingAllPosts;


    return (
        <div className="min-h-screen bg-zinc-50">
            <div className="flex min-h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
                    <div className="flex w-full max-w-5xl flex-col gap-6">
                        <SearchBar
                            value={query}
                            onChange={setQuery}
                            placeholder="Search posts..."
                            className="w-full"
                        />

                        <div className="flex flex-col gap-3 w-full">
                            <h2 className="font-semibold text-base">Popular Tags</h2>
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
                            <span className="font-semibold text-base">Explore</span>
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
