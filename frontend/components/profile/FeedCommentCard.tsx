"use client";

import { Heart } from "lucide-react";
import Link from "next/link";

import type { CommentResponse } from "@/api/models";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { cn } from "@/lib/utils";

type Props = {
  comment: CommentResponse;
};

export function FeedCommentCard({ comment }: Props) {
  const avatar = (
    <Avatar className="h-9 w-9 bg-slate-300">
      <AvatarFallback />
    </Avatar>
  );

  return (
    <article className="space-y-2">
      <div className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
        {comment.user?.id ? (
          <Link
            href={`/profile/${comment.user.id}`}
            className="shrink-0 rounded-full outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
            aria-label={`View ${
              comment.user.username || "user"
            }'s profile`}
          >
            {avatar}
          </Link>
        ) : (
          avatar
        )}
        <div className="mb-1 flex items-baseline gap-3">
          {comment.user?.id ? (
            <Link
              href={`/profile/${comment.user.id}`}
              className="text-sm font-semibold hover:underline"
            >
              {comment.user.username || "otherperson"}
            </Link>
          ) : (
            <p className="text-sm font-semibold">
              {comment.user?.username || "otherperson"}
            </p>
          )}
          <p className="text-sm leading-6 text-slate-700">
            {comment.description}
          </p>
        </div>
        <button
          type="button"
          className="mt-1 text-slate-800 transition-opacity hover:opacity-70"
        >
          <Heart
            className={cn("h-5 w-5", comment.is_liked && "fill-current")}
          />
        </button>
      </div>
      <div className="ml-12 border-l border-slate-300 bg-slate-100/75 px-3 py-2">
        <p className="text-xs text-slate-600">
          Reply thread starts here (Reddit-style spacing + divider)
        </p>
      </div>
    </article>
  );
}
