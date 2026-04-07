"use client";

import type { PostResponse } from "@/api/models";
import { FeedPostCard } from "@/components/profile/FeedPostCard";

type Props = {
  likedPosts: PostResponse[];
};

export function ProfileLikesFeed({ likedPosts }: Props) {
  return (
    <div className="mt-4 space-y-6">
      {likedPosts.map((post) => (
        <FeedPostCard key={post.id} post={post} filledHeart />
      ))}
    </div>
  );
}
