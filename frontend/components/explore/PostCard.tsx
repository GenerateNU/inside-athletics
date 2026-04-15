"use client";

import Link from "next/link";
import { Heart, MessageCircle } from "lucide-react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";
import type { PostResponse } from "@/api/models/PostResponse";

interface PostCardProps {
  post: PostResponse;
}

export function PostCard({ post }: PostCardProps) {
  const authorName = post.is_anonymous
    ? "Anonymous"
    : `${post.author.first_name} ${post.author.last_name}`;

  const authorInitials = post.is_anonymous
    ? "A"
    : `${post.author.first_name[0] ?? ""}${post.author.last_name[0] ?? ""}`.toUpperCase();

  const tagLabels = [
    ...(post.sport?.name ? [post.sport.name] : []),
    ...(post.tags?.map((t) => t.name) ?? []),
  ];

  return (
    <Link href={`/posts/${post.id}`} className="block border-b border-zinc-100 px-6 py-5 hover:bg-zinc-50 transition-colors">
      {/* Sport / tag labels */}
      <div className="mb-1.5 flex flex-wrap gap-3 text-xs text-zinc-400">
        {tagLabels.map((label) => (
          <span key={label}>{label}</span>
        ))}
      </div>

      {/* Title */}
      <h3 className="mb-2 text-sm font-semibold leading-snug text-zinc-900">
        {post.title}
      </h3>

      {/* Content preview */}
      <p className="mb-4 line-clamp-3 text-sm leading-relaxed text-zinc-500">
        {post.content}
      </p>

      {/* Footer */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4 text-zinc-400">
          <span className="flex items-center gap-1.5 text-xs">
            <Heart
              className={cn(
                "size-3.5",
                post.is_liked && "fill-red-500 text-red-500",
              )}
            />
            {post.like_count ?? 0}
          </span>
          <span className="flex items-center gap-1.5 text-xs">
            <MessageCircle className="size-3.5" />
            {post.comment_count ?? 0}
          </span>
        </div>

        <div className="flex items-center gap-2">
          <span className="text-xs text-zinc-500">{authorName}</span>
          <Avatar size="sm">
            <AvatarFallback className="bg-zinc-100 text-[10px] font-medium text-zinc-600">
              {authorInitials}
            </AvatarFallback>
          </Avatar>
        </div>
      </div>
    </Link>
  );
}
