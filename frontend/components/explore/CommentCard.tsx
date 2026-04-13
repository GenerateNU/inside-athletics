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
        <div className="mt-3 flex gap-2">
          <textarea
            value={replyText}
            onChange={(e) => setReplyText(e.target.value)}
            placeholder="Write a reply..."
            rows={2}
            className="min-h-0 flex-1 resize-none rounded-md border border-zinc-200 bg-zinc-50 px-3 py-2 text-base text-zinc-900 placeholder:text-zinc-400 focus:border-zinc-400 focus:bg-white focus:outline-none"
          />
          <div className="flex flex-col gap-1">
            <Button
              size="sm"
              onClick={handleReplySubmit}
              disabled={!replyText.trim() || submittingReply}
            >
              {submittingReply ? "Sending..." : "Send"}
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => { setReplyOpen(false); setReplyText(""); }}
            >
              Cancel
            </Button>
          </div>
        </div>
      )}

      {repliesOpen && (
        <div className="mt-3 border-l-2 border-zinc-100 pl-4">
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
  );
}
