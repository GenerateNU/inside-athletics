"use client";

import Link from "next/link";
import { Heart, MessageSquareText, UserRound } from "lucide-react";
import { Badge } from "@/components/post/Badge"
import { Tag } from "@/components/post/Tag"
import { cn } from "@/lib/utils";
import type { PostResponse } from "@/api/models/PostResponse";
import { PremiumPostResponse } from "@/api";
import MediaDisplay from "./MediaDisplay";

type SmallPostProps = React.ComponentProps<"div"> & {
    post: PremiumPostResponse;
};

export default function PremiumSmallPost({ post, className, ...props }: SmallPostProps) {
    const authorName = `${post.author.first_name} ${post.author.last_name}`;
    const pfpURL = post.author?.profile_picture
    

    return (
            <div className={cn("bg-white rounded-2xl border border-gray-200 p-5 w-full shadow-sm hover:shadow-md transition-shadow cursor-pointer", className)} {...props}>
                <h2 className="font-bold text-gray-900 text-base mb-3 text-left">{post.title}</h2>

                

                <p className="text-gray-700 text-sm leading-relaxed mb-4 text-left">{post.content}</p>

                {post.media && <MediaDisplay media={post.media}/> }
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2">
                        { pfpURL
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