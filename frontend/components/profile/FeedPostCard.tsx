"use client";

import { Heart, Bookmark } from "lucide-react";
import Link from "next/link";

import type { PostResponse } from "@/api/models";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";

type Props = {
  post: PostResponse;
  filledHeart: boolean;
};

export function FeedPostCard({ post, filledHeart }: Props) {
  const avatar = (
    <Avatar className="h-9 w-9 bg-slate-300">
      <AvatarFallback />
    </Avatar>
  );

  return (
    <article className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
      {post.author?.id ? (
        <Link
          href={`/profile/${post.author.id}`}
          className="shrink-0 rounded-full outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
          aria-label={`View ${
            post.author.username || "user"
          }'s profile`}
        >
          {avatar}
        </Link>
      ) : (
        avatar
      )}
      <div>
        <div className="mb-1 flex items-baseline gap-3">
          {post.author?.id ? (
            <Link
              href={`/profile/${post.author.id}`}
              className="text-sm font-semibold hover:underline"
            >
              {post.author.username || "anonymous"}
            </Link>
          ) : (
            <p className="text-sm font-semibold">anonymous</p>
          )}
          <p className="text-sm leading-6 text-slate-700">{post.content}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          {post.college?.name ? (
            <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
              <Bookmark className="mr-1 h-3 w-3" />
              {post.college.name}
            </Badge>
          ) : null}
          {post.sport?.name ? (
            <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
              {post.sport.name}
            </Badge>
          ) : null}
        </div>
      </div>
      <button
        type="button"
        className="mt-1 text-[#2f6aa0] transition-opacity hover:opacity-70"
      >
        <Heart className={cn("h-5 w-5", filledHeart && "fill-current")} />
      </button>
    </article>
  );
}
