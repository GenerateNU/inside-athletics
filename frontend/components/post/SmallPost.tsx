"use client";

import Link from "next/link";
import { Heart, MessageSquareText, UserRound } from "lucide-react";
import { Badge } from "@/components/post/Badge"
import { Tag } from "@/components/post/Tag"
import { cn } from "@/lib/utils";
import type { PostResponse } from "@/api/models/PostResponse";

type SmallPostProps = React.ComponentProps<"div"> & {
    post: PostResponse;
};

export default function SmallPost({ post, className, ...props }: SmallPostProps) {
    const authorName = post.is_anonymous
        ? "Anonymous"
        : `${post.author.first_name} ${post.author.last_name}`;

    const pfpURL = post.author?.profile_picture

    return (
        <Link href={`/posts/${post.id}`} className="block w-full">
            <div className={cn("bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm hover:shadow-md transition-shadow cursor-pointer", className)} {...props}>
                <h2 className="font-bold text-gray-900 text-base mb-3 text-left">{post.title}</h2>

                <div className="flex flex-wrap gap-2 mb-3">
                    {post.sport && <Tag label={post.sport.name} />}
                    {post.college && <Tag label={post.college.name} />}
                    {post.tags?.map((tag) => (
                        <Tag key={tag.id} label={tag.name} />
                    ))}
                </div>

                <p className="text-gray-700 text-sm leading-relaxed mb-4 text-left">{post.content}</p>

                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <Badge
                            icon={<Heart className="text-[#3E7DBB] size-5 shrink-0" />}
                            count={post.like_count ?? 0}
                        />
                        <Badge
                            icon={<MessageSquareText className="text-[#3E7DBB] size-5 shrink-0" />}
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
        </Link>
    );
}