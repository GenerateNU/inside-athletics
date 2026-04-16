"use client";

import { useGetApiV1PostsPremium } from "@/api/hooks";
import PremiumSmallPost from "@/components/post/PremiumSmallPost";
import { Navbar } from "@/components/ui/navbar";
import { useSession } from "@/utils/SessionContext";

export default function InsiderContentPage() {
    const session = useSession();
    const enabled = !!session?.access_token;
    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const { data, isLoading, isError } = useGetApiV1PostsPremium(
        {},
        {
            query: { enabled },
            client: { headers: authHeaders },
        }
    );
    const posts = data?.posts ?? [];
    return (
        <div className="min-h-screen bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60 w-full">
            <div className="flex min-h-screen">
                <Navbar className="h-screen shrink-0" />
                <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10">
                    <div className="w-full">
                        <h1 className="text-4xl font-bold text-gray-900 mb-6">Insider Content</h1>
                        {isLoading && <p className="text-gray-500">Loading...</p>}
                        {isError && <p className="text-red-500">Failed to load insider content.</p>}
                        {!isLoading && !isError && posts.length === 0 && (
                            <p className="text-gray-500">No insider content available.</p>
                        )}
                        <div className="flex flex-col gap-4">
                            {posts.map((post) => (
                                <PremiumSmallPost key={post.id} post={post} />
                            ))}
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}
