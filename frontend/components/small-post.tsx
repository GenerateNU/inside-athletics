"use client";

import { Heart, MessageSquareText, UserRound } from "lucide-react";
import { Badge } from "@/components/explore/Badge"
import { Tag } from "@/components/explore/Tag"
import { useSession } from "@/utils/SessionContext";
import { cn } from "@/lib/utils";
import { useGetApiV1PostById, useGetApiV1UserById} from "@/api/hooks";

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


    const { } = useGetApiV1UserById(id, {
        client: {
            headers: { Authorization: `Bearer ${session?.access_token}` },
        },
    });

    const user = unwrapBody<import("@/api").GetUserResponse>(raw);

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

    const pfpURL = user?.profile_picture_url

    return (
        <div className="bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm">
            <h2 className="font-bold text-gray-900 text-base mb-3 text-left">{post.title}</h2>

            <div className="flex gap-2 mb-3">
                <Tag label={post.sport.name} />
                <Tag label={post.college.name} />
            </div>

            <p className="text-gray-700 text-sm leading-relaxed mb-4 text-left">{post.content}</p>

            <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                    <Badge
                        icon={<Heart size={16} className="text-blue-500" />}
                        count={post.like_count ?? 0}
                    />
                    <Badge
                        icon={<MessageSquareText size={16} className="text-blue-500" />}
                        count={post.comment_count ?? 0}
                    />
                </div>

                <div className="flex items-center gap-2">
                    {post.is_anonymous
                        ? <div className="w-7 h-7 rounded-full bg-zinc-200 flex items-center justify-center shrink-0">
                            <UserRound size={16} className="text-zinc-500" />
                          </div>
                        : pfpURL
                            ? <img
                                src={pfpURL}
                                alt={authorName}
                                className="w-7 h-7 rounded-full object-cover shrink-0"
                              />
                            : <div className="w-7 h-7 rounded-full bg-zinc-200 shrink-0" />
                    }
                    <span className="text-sm font-semibold text-gray-800">{authorName}</span>
                </div>
            </div>
        </div>
    );
}