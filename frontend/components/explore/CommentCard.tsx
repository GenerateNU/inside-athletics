"use client";

import { useState } from "react";
import { MoreHorizontal } from "lucide-react";
import { useQueryClient } from "@tanstack/react-query";
import { Button } from "@/components/ui/button";
import {
  usePostApiV1Comment,
  useListApiV1CommentByIdReplies,
  listApiV1PostByPostIdCommentsQueryKey,
  listApiV1CommentByIdRepliesQueryKey,
} from "@/api/hooks";
import type { CommentResponse } from "@/api/models/CommentResponse";
import { CommentItem } from "./CommentItem";

interface CommentCardProps {
  comment: CommentResponse;
  postId: string;
  authHeaders?: Record<string, string>;
}

export function CommentCard({ comment, postId, authHeaders }: CommentCardProps) {
  const queryClient = useQueryClient();

  const [replyOpen, setReplyOpen] = useState(false);
  const [replyText, setReplyText] = useState("");
  const [repliesOpen, setRepliesOpen] = useState(false);
  const [localHasReplies, setLocalHasReplies] = useState(comment.has_replies);

  const { data: replies, isLoading: loadingReplies } = useListApiV1CommentByIdReplies(
    comment.id,
    {
      query: { enabled: repliesOpen },
      client: { headers: authHeaders },
    },
  );

  const { mutate: submitReply, isPending: submittingReply } =
    usePostApiV1Comment({ client: { headers: authHeaders } });

  function handleReplySubmit() {
    if (!replyText.trim()) return;
    submitReply(
      {
        data: {
          description: replyText.trim(),
          is_anonymous: false,
          parent_comment_id: comment.id,
          post_id: postId,
        },
      },
      {
        onSuccess: () => {
          setReplyText("");
          setReplyOpen(false);
          setLocalHasReplies(true);
          setRepliesOpen(true);
          queryClient.invalidateQueries({
            queryKey: listApiV1PostByPostIdCommentsQueryKey(postId),
          });
          queryClient.invalidateQueries({
            queryKey: listApiV1CommentByIdRepliesQueryKey(comment.id),
          });
        },
      },
    );
  }

  return (
    <div className="py-4">
      <svg width="0" height="0" style={{ position: "absolute" }}>
        <defs>
          <linearGradient id="green-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stopColor="#00804D" />
            <stop offset="100%" stopColor="#043D26" />
          </linearGradient>
        </defs>
      </svg>
      <div className={`relative ${repliesOpen ? "pl-4" : ""}`}>
        {repliesOpen && (
          <div className="absolute inset-y-0 left-0 w-0.5 rounded-full bg-zinc-200" />
        )}

        <div className="flex items-start justify-between gap-2">
          <div className="min-w-0 flex-1">
            <CommentItem
              comment={comment}
              authHeaders={authHeaders}
              onReply={() => setReplyOpen((o) => !o)}
              showRepliesToggle={localHasReplies}
              repliesOpen={repliesOpen}
              onToggleReplies={() => setRepliesOpen((o) => !o)}
            />
          </div>

          <button
            type="button"
            className="shrink-0 rounded-md p-1 text-zinc-300 transition-colors hover:bg-zinc-100 hover:text-zinc-500"
            aria-label="More options"
          >
            <MoreHorizontal className="size-4" />
          </button>
        </div>

        {replyOpen && (
          <div className="mt-3 flex gap-2 relative">
            <textarea
              value={replyText}
              onChange={(e) => setReplyText(e.target.value)}
              placeholder="Write a reply..."
              rows={4}
              className="min-h-0 flex-1 resize-none rounded-2xl border border-[#3E7DBB] bg-white px-3 py-2 text-base text-zinc-900 placeholder:text-zinc-400"
            />
            <Button
                className={"absolute bottom-1 right-2 rounded-3xl bg-[#A8C8E8] text-[#E8F1FA]"}
                size="lg"
                onClick={handleReplySubmit}
                disabled={!replyText.trim() || submittingReply}
              >
                {submittingReply ? "Posting..." : "Post"}
              </Button>
          </div>
        )}

        {repliesOpen && (
          <div className="mt-3 ml-5">
            {loadingReplies ? (
              <div className="space-y-3">
                {[1, 2].map((i) => (
                  <div key={i} className="space-y-1">
                    <div className="h-3 w-20 animate-pulse rounded bg-zinc-100" />
                    <div className="h-8 animate-pulse rounded bg-zinc-100" />
                  </div>
                ))}
              </div>
            ) : !replies || replies.length === 0 ? (
              <p className="text-xs text-zinc-400">No replies yet.</p>
            ) : (
              <div className="space-y-3">
                {replies.map((reply) => (
                  <CommentItem key={reply.id} comment={reply} authHeaders={authHeaders} />
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
