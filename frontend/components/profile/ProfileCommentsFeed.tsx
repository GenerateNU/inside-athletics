"use client";

import type { CommentResponse } from "@/api/models";
import { FeedCommentCard } from "@/components/profile/FeedCommentCard";

type Props = {
  comments: CommentResponse[];
};

export function ProfileCommentsFeed({ comments }: Props) {
  return (
    <div className="mt-4 space-y-6">
      {comments.map((comment) => (
        <FeedCommentCard key={comment.id} comment={comment} />
      ))}
    </div>
  );
}
