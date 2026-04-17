"use client";

import type { PostResponse } from "@/api/models";
import { FeedPostCard } from "@/components/profile/FeedPostCard";

type Props = {
  posts: PostResponse[];
};

export function ProfilePostsFeed({ posts }: Props) {
  return (
    <div className="mt-4 space-y-6">
      {posts.map((post) => (
        <FeedPostCard
          key={post.id}
          post={post}
          filledHeart={Boolean(post.is_liked)}
        />
      ))}
    </div>
  );
}
