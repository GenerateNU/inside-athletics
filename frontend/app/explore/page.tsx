"use client";

// Currently the popular tags are just replaced with user tag follows!

import { useState } from "react";
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
    getApiV1TagByIdQueryOptions,
} from "@/api/hooks";
import { CancellableTag } from "@/components/filtering/CancellableTag";
import { GetTagResponse } from "@/api";


export default function ExplorePage() {
    const session = useSession();
    const enabled = !!session?.access_token;

    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const router = useRouter();
    const [query, setQuery] = useState("");
    const [activeTags, setActiveTags] = useState<GetTagResponse[]>([]);

    function toggleTag(tag: GetTagResponse) {
        setActiveTags((prev) =>
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

    const tagIds = tagsfollowsData?.tag_ids ?? [];
    const tagQueries = useQueries({
        queries: tagIds.map((id) =>
            getApiV1TagByIdQueryOptions(id, { headers: authHeaders }),
        ),
    });
    const followedTags = tagQueries
        .map((q) => q.data)
        .filter((t) => t !== undefined);

    

    const { data: allPostsData, isLoading: loadingAllPosts } = useGetApiV1Posts(
        { },
        { 
            query: { enabled },
            client: { headers: authHeaders } },
    );
    console.log(allPostsData)

    const { data: filteredPostsData, isLoading: loadingFilteredPosts } = useGetApiV1PostsFilter(
        { tag_ids: activeTags.map((t) => t.id).join(",") },
        {
            query: { enabled: activeTags.length > 0 },
            client: { headers: authHeaders },
        },
    );

    const posts = activeTags.length > 0
        ? (filteredPostsData?.posts ?? [])
        : (allPostsData?.posts ?? []);
    const isLoading = activeTags.length > 0 ? loadingFilteredPosts : loadingAllPosts;


    return (
        <div className="min-h-screen bg-zinc-50">
            <div className="flex min-h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
                    <div className="flex w-full max-w-5xl flex-col gap-6">
                        <SearchBar
                            value={query}
                            onChange={setQuery}
                            onSubmit={() => query.trim() && router.push(`/search?q=${encodeURIComponent(query.trim())}`)}
                            placeholder="Search posts..."
                            className="w-full"
                        />

                        <div className="flex flex-col gap-3 w-full">
                            <h2 className="font-semibold text-base">Popular Tags</h2>
                            <div className="flex flex-wrap gap-2">
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
