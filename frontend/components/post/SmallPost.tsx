"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { Heart, MessageSquareText, UserRound } from "lucide-react";
import { Badge } from "@/components/post/Badge"
import { Tag } from "@/components/post/Tag"
import { cn } from "@/lib/utils";
import type { PostResponse } from "@/api/models/PostResponse";

type SmallPostProps = React.ComponentProps<"div"> & {
    post: PostResponse;
};

export default function SmallPost({ post, className, ...props }: SmallPostProps) {
    const router = useRouter();

    const authorName = post.is_anonymous
        ? "Anonymous"
        : `${post.author.first_name} ${post.author.last_name}`;

    const pfpURL = post.author?.profile_picture

    const tags = post.tags
    const tagsEnable = tags != null
    const sportEnable = post.sport != null 
    const collegeEnable = post.college != null

    const profileId = !post.is_anonymous && post.author?.id ? post.author.id : null;

    return (
        <div
            className={cn("bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm hover:shadow-md transition-shadow cursor-pointer", className)}
            onClick={() => router.push(`/posts/${post.id}`)}
            onKeyDown={(e) => {
                if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    router.push(`/posts/${post.id}`);
                }
            }}
            role="button"
            tabIndex={0}
            {...props}
        >
                <h2 className="font-bold text-gray-900 text-base mb-3 text-left">{post.title}</h2>

                <div className="flex gap-2 mb-3">
                    {sportEnable && <Tag label={post.sport.name} />}
                    {collegeEnable && <Tag label={post.college.name} />}
                    {tagsEnable && tags?.map((tag) => (
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

                    {profileId ? (
                        <Link
                            href={`/profile/${profileId}`}
                            onClick={(e) => e.stopPropagation()}
                            className="flex items-center gap-2 rounded-md outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
                        >
                        {pfpURL
                            ? <img
                                src={pfpURL}
                                alt=""
                                className="w-7 h-7 rounded-full object-cover shrink-0"
                              />
                            : <div className="w-7 h-7 rounded-full bg-zinc-200 shrink-0" />
                        }
                        <span className="text-sm font-semibold text-gray-800">{authorName}</span>
                        </Link>
                    ) : (
                        <div className="flex items-center gap-2">
                            {post.is_anonymous
                                ? <div className="w-7 h-7 rounded-full bg-zinc-200 flex items-center justify-center shrink-0">
                                    <UserRound size={16} className="text-zinc-500" />
                                  </div>
                                : pfpURL
                                    ? <img
                                        src={pfpURL}
                                        alt=""
                                        className="w-7 h-7 rounded-full object-cover shrink-0"
                                      />
                                    : <div className="w-7 h-7 rounded-full bg-zinc-200 shrink-0" />
                            }
                            <span className="text-sm font-semibold text-gray-800">{authorName}</span>
                        </div>
                    )}
                </div>
        </div>
    );
}