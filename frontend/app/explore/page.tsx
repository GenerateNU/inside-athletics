"use client";

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
    useGetApiV1Sports,
    useGetApiV1CollegesSearch,
    useGetApiV1Posts,
    useGetApiV1UserTagFollows,
    useGetApiV1PostsSearch,
    useGetApiV1TagsSearch,

    getApiV1TagByIdQueryOptions,
} from "@/api/hooks";
import type { GetTagFollowsByUserResponse } from "@/api/models/GetTagFollowsByUserResponse";
import type { PostResponse } from "@/api/models/PostResponse";
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
    const [activeTag, setActiveTag] = useState<GetTagResponse | null>(null);

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
        { tag_ids: activeTag?.id },
        {
            query: { enabled: !!activeTag },
            client: { headers: authHeaders },
        },
    );
    console.log(filteredPostsData)
    
    const posts = activeTag
        ? (filteredPostsData?.posts ?? [])
        : (allPostsData?.posts ?? []);
    const isLoading = activeTag ? loadingFilteredPosts : loadingAllPosts;


    return (
        <div className="min-h-screen bg-zinc-50">
            <div className="flex min-h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
                    <div className="flex w-full max-w-5xl flex-col gap-10">
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
                                        onClick={() => setActiveTag(activeTag?.id === tag.id ? null : tag)}
                                    >
                                        <Tag label={tag.name} />
                                    </button>
                                ))}
                            </div>
                        </div>

                        <div className="flex items-center gap-2 w-full">
                            <span className="font-semibold text-base">Explore</span>
                            {activeTag && <CancellableTag
                                label={activeTag.name}
                                onRemove={() => setActiveTag(null)}
                            />}
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
