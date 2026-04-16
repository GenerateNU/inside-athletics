"use client";

import { useState } from "react";
// import { useSession } from "@/utils/SessionContext";

import { SearchBar } from "@/components/post/SearchBar";
import SmallPost from "@/components/post/SmallPost";
import { Navbar } from "@/components/ui/navbar";
import { Tag } from "@/components/post/Tag";

// import {
//     useGetApiV1Posts,
//     useGetApiV1TagsSearch,
//     useGetApiV1TagByTagIdPosts,
// } from "@/api/hooks";
import type { GetTagResponse } from "@/api/models/GetTagResponse";
import type { PostResponse } from "@/api/models/PostResponse";
import { CancellableTag } from "@/components/filtering/CancellableTag";

// --- MOCK DATA (remove fallbacks once API is fixed) ---
const MOCK_TAGS: GetTagResponse[] = [
    { id: "1", name: "Basketball" },
    { id: "2", name: "Swimming" },
    { id: "3", name: "Football" },
    { id: "4", name: "Northeastern University" },
    { id: "5", name: "Binghamton University" },
    { id: "6", name: "Baseball" },
    { id: "7", name: "Diving" },
];

const MOCK_USER = {
    id: "u1",
    first_name: "Alex",
    last_name: "Johnson",
    username: "alexj",
    email: "alex@example.com",
    bio: null,
    college: "Northeastern University",
    sport: "Basketball",
    division: 1,
    account_type: false,
    expected_grad_year: 2026,
    profile_picture: "",
    verified_athelete_status: "verified",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    deleted_at: null,
};

const MOCK_COLLEGE = {
    id: "c1",
    name: "Northeastern University",
    city: "Boston",
    state: "MA",
    division_rank: 1 as const,
    academic_rank: 49,
    logo: "",
    website: "https://northeastern.edu",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    deleted_at: null,
};

const MOCK_SPORT = {
    id: "s1",
    name: "Basketball",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
};

const MOCK_POSTS: PostResponse[] = [
    {
        id: "p1",
        title: "Thoughts on our last home game",
        content: "The energy in the arena was incredible. We came back from 12 down in the second half and the crowd never stopped believing.",
        author: MOCK_USER,
        college: MOCK_COLLEGE,
        sport: MOCK_SPORT,
        tags: [{ id: "1", name: "Basketball", created_at: "2024-01-01T00:00:00Z", updated_at: "2024-01-01T00:00:00Z", deleted_at: null }],
        is_anonymous: false,
        is_verified_athlete: true,
        like_count: 34,
        comment_count: 12,
    },
    {
        id: "p2",
        title: "Recruiting tips for D1 swimmers",
        content: "Start reaching out to coaches junior year. Send a swim resume with your best times and a short highlight reel. Persistence matters more than perfection.",
        author: { ...MOCK_USER, id: "u2", first_name: "Jamie", last_name: "Lee", sport: "Swimming" },
        college: { ...MOCK_COLLEGE, id: "c2", name: "Binghamton University" },
        sport: { ...MOCK_SPORT, id: "s2", name: "Swimming" },
        tags: [{ id: "2", name: "Swimming", created_at: "2024-01-01T00:00:00Z", updated_at: "2024-01-01T00:00:00Z", deleted_at: null }],
        is_anonymous: false,
        is_verified_athlete: true,
        like_count: 57,
        comment_count: 8,
    },
    {
        id: "p3",
        title: "Balancing practice and midterms",
        content: "Block study sessions right after practice while the routine is still fresh. Professors are usually understanding if you communicate early.",
        author: { ...MOCK_USER, id: "u3", first_name: "Morgan", last_name: "Smith" },
        college: MOCK_COLLEGE,
        sport: { ...MOCK_SPORT, id: "s3", name: "Football" },
        tags: [],
        is_anonymous: true,
        is_verified_athlete: false,
        like_count: 21,
        comment_count: 5,
    },
];
// --- END MOCK DATA ---

export default function ExplorePage() {
    // const session = useSession();
    // const authHeaders = session?.access_token
    //     ? { Authorization: `Bearer ${session.access_token}` }
    //     : undefined;

    const [query, setQuery] = useState("");
    const [activeTag, setActiveTag] = useState<GetTagResponse | null>(null);

    // const { data: tagsData } = useGetApiV1TagsSearch(
    //     {},
    //     { client: { headers: authHeaders } },
    // );

    // const { data: allPostsData, isLoading: loadingAllPosts } = useGetApiV1Posts(
    //     {},
    //     { client: { headers: authHeaders } },
    // );

    // const { data: tagPostsData, isLoading: loadingTagPosts } = useGetApiV1TagByTagIdPosts(
    //     activeTag?.id ?? "",
    //     {},
    //     {
    //         query: { enabled: !!activeTag },
    //         client: { headers: authHeaders },
    //     },
    // );

    // const popularTags = tagsData?.results ?? [];
    // const posts = activeTag
    //     ? (tagPostsData?.post_ids ?? [])
    //     : (allPostsData?.posts ?? []);
    // const isLoading = activeTag ? loadingTagPosts : loadingAllPosts;
    const popularTags = MOCK_TAGS;
    const posts = MOCK_POSTS;
    const isLoading = false;

    return (
        <div className="min-h-screen bg-zinc-50">
            <div className="flex min-h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
                    <div className="flex w-full max-w-3xl flex-col gap-10">
                        <SearchBar
                            value={query}
                            onChange={setQuery}
                            placeholder="Search posts..."
                            className="w-full"
                        />

                        <div className="flex flex-col gap-3 w-full">
                            <h2 className="font-semibold text-base">Popular Tags</h2>
                            <div className="flex flex-wrap gap-2">
                                {popularTags.map((tag) => (
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
