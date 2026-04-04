"use client";

import { BookOpen, Heart, MessageSquareText } from "lucide-react";
import { useSession } from "@/utils/SessionContext";
import { cn } from "@/lib/utils";
import { useGetApiV1PostById, } from "@/api/hooks";

type SmallPostProps = React.ComponentProps<"div"> & {
    id: string;
};


function unwrapBody<T>(value: unknown): T | undefined {
    let current = value;
    for (let depth = 0; depth < 5; depth += 1) {
        if (!current || typeof current !== "object") return current as T | undefined;
        if ("body" in current && current.body !== undefined) { current = current.body; continue; }
        if ("Body" in current && current.Body !== undefined) { current = current.Body; continue; }
        return current as T | undefined;
    }
    return current as T | undefined;
}


export default function SmallPost({ id, className, ...props }: SmallPostProps) {
    const session = useSession();

    const { data: raw, isLoading, isError } = useGetApiV1PostById(id, {
        client: {
            headers: { Authorization: `Bearer ${session?.access_token}` },
        },
    });

    const post = unwrapBody<import("@/api").PostResponse>(raw);

    if (isLoading) {
        return (
            <div className={cn("bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm animate-pulse", className)} {...props}>
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-3" />
                <div className="flex gap-2 mb-3">
                    <div className="h-7 bg-gray-200 rounded-md w-24" />
                    <div className="h-7 bg-gray-200 rounded-md w-20" />
                </div>
                <div className="space-y-2 mb-4">
                    <div className="h-3 bg-gray-200 rounded w-full" />
                    <div className="h-3 bg-gray-200 rounded w-5/6" />
                </div>
            </div>
        );
    }

    if (isError || !post) {
        return (
            <div className={cn("bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm", className)} {...props}>
                <p className="text-sm text-gray-400">Failed to load post.</p>
            </div>
        );
    }

    const authorName = post.is_anonymous
        ? "Anonymous"
        : `${post.author.first_name} ${post.author.last_name}`;


    return (
        <div className="bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm">
            <h2 className="font-bold text-gray-900 text-base mb-3 text-left">{post.title}</h2>

            <div className="flex gap-2 mb-3">
                <span className="flex items-center gap-1.5 border border-gray-300 rounded-md px-3 py-1 text-sm text-gray-700">
                    <BookOpen size={15} className="text-blue-600" />
                    {post.college.name}
                </span>
                <span className="border border-gray-300 rounded-md px-3 py-1 text-sm text-gray-700">
                    {post.sport.name}
                </span>
            </div>

            <p className="text-gray-700 text-sm leading-relaxed mb-4 text-left">{post.content}</p>

            <div className="flex flex-col items-end gap-3">
                <div className="flex items-center gap-3">
                    <button className="flex items-center gap-1.5 border border-gray-200 rounded-full px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50">
                        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
                            <defs>
                                <linearGradient id="heartGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                                    <stop offset="0%" stopColor="#3b82f6" />
                                    <stop offset="100%" stopColor="#22c55e" />
                                </linearGradient>
                            </defs>
                            <path d="M19 14c1.49-1.46 3-3.21 3-5.5A5.5 5.5 0 0 0 16.5 3c-1.76 0-3 .5-4.5 2-1.5-1.5-2.74-2-4.5-2A5.5 5.5 0 0 0 2 8.5c0 2.3 1.5 4.05 3 5.5l7 7Z" stroke="url(#heartGradient)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" fill="url(#heartGradient)" />
                        </svg>
                        {post.like_count ?? 0}
                    </button>
                    <button className="flex items-center gap-1.5 border border-gray-200 rounded-full px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-50">
                        <MessageSquareText size={16} className="text-blue-500" />
                        {post.comment_count ?? 0}
                    </button>
                </div>

                <div className="flex items-center gap-2">
                    {!post.is_anonymous && (
                        <div className="w-7 h-7 rounded-full bg-zinc-200 shrink-0" />
                    )}
                    <span className="text-sm font-semibold text-gray-800">{authorName}</span>
                </div>
            </div>
        </div>
    );
}