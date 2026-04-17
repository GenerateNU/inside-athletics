"use client";

import { useGetApiV1PostsPremium } from "@/api/hooks";
import PremiumSmallPost from "@/components/post/PremiumSmallPost";
import { Navbar } from "@/components/ui/navbar";
import { useSession, usePermissions } from "@/utils/SessionContext";
import { useState } from "react";
import { Plus } from "lucide-react";
import CreatePremiumPostPopup from "@/components/ui/create-premium-post-popup";

export default function InsiderContentPage() {
    const session = useSession();
    const { isAdmin } = usePermissions();
    const enabled = !!session?.access_token;
    const authHeaders = session?.access_token
        ? { Authorization: `Bearer ${session.access_token}` }
        : undefined;

    const [showCreatePopup, setShowCreatePopup] = useState(false);

    const { data, isLoading, isError, refetch } = useGetApiV1PostsPremium(
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
                        <div className="flex items-center justify-between mb-6">
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
        </div>
    );
}
