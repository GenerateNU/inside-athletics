"use client";

import { useState } from "react";
import { Heart, MessageCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import {
  usePostApiV1CommentLike,
  useDeleteApiV1CommentLikeById,
} from "@/api/hooks";
import type { CommentResponse } from "@/api/models/CommentResponse";

interface CommentItemProps {
  comment: CommentResponse;
  authHeaders?: Record<string, string>;
  onReply?: () => void;
  showRepliesToggle?: boolean;
  repliesOpen?: boolean;
  onToggleReplies?: () => void;
}

export function CommentItem({
  comment,
  authHeaders,
  onReply,
  showRepliesToggle,
  repliesOpen,
  onToggleReplies,
}: CommentItemProps) {
  const [isLiked, setIsLiked] = useState(comment.is_liked);
  const [likeCount, setLikeCount] = useState(comment.like_count);

  const authorName = comment.is_anonymous
    ? "Anonymous"
    : comment.user
      ? `${comment.user.first_name} ${comment.user.last_name}`
      : "Unknown";

  const { mutate: likeComment } = usePostApiV1CommentLike({
    client: { headers: authHeaders },
  });

  const { mutate: unlikeComment } = useDeleteApiV1CommentLikeById({
    client: { headers: authHeaders },
  });

  function handleLikeToggle() {
    if (isLiked) {
      setIsLiked(false);
      setLikeCount((c) => c - 1);
      unlikeComment(
        { id: comment.id },
        {
          onError: () => {
            setIsLiked(true);
            setLikeCount((c) => c + 1);
          },
        },
      );
    } else {
      setIsLiked(true);
      setLikeCount((c) => c + 1);
      likeComment(
        { data: { comment_id: comment.id } },
        {
          onError: () => {
            setIsLiked(false);
            setLikeCount((c) => c - 1);
          },
        },
      );
    }
  }

  return (
    <div>
      <p className="mb-1 text-base font-semibold text-black">{authorName}</p>
      <p className="text-base leading-relaxed text-black">
        {comment.description}
      </p>
      <div className="mt-2 flex items-center gap-4">
        <button
          type="button"
          onClick={handleLikeToggle}
          className={
            "flex items-center gap-1 text-xs text-black transition-colors"
          }
        >
          <Heart
            className={"size-3.5"}
            stroke="url(#green-gradient)"
            fill={isLiked ? "url(#green-gradient)" : "none"}
          />
          {likeCount}
        </button>

        {onReply && (
          <button
            type="button"
            onClick={onReply}
            className="flex items-center gap-1 text-xs text-black"
          >
            <MessageCircle
              className={"size-3.5"}
              stroke="url(#green-gradient)"
            />
            Reply
          </button>
        )}

        {showRepliesToggle && onToggleReplies && (
          <button
            type="button"
            onClick={onToggleReplies}
            className="text-xs font-medium text-zinc-400 transition-colors hover:text-zinc-700"
          >
            {repliesOpen ? "Hide replies" : "View replies"}
          </button>
        )}
      </div>
    </div>
  );
}
