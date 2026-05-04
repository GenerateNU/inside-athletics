"use client";

import { Heart } from "lucide-react";
import Link from "next/link";

import type { CommentResponse } from "@/api/models";
import { cn } from "@/lib/utils";

type Props = {
  comment: CommentResponse;
};

export function FeedCommentCard({ comment }: Props) {
  return (
    <article className="space-y-2">
      <div className="grid grid-cols-[1fr_auto] items-start gap-3">
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

    </article>
  );
}
