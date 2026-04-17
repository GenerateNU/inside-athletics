"use client";

import type { CommentResponse, PostResponse } from "@/api/models";
import { ProfileCommentsFeed } from "@/components/profile/ProfileCommentsFeed";
import { ProfileLikesFeed } from "@/components/profile/ProfileLikesFeed";
import { ProfilePostsFeed } from "@/components/profile/ProfilePostsFeed";
import type { FeedView } from "@/components/profile/types";

type Props = {
  posts: PostResponse[];
  likedPosts: PostResponse[];
  comments: CommentResponse[];
  activeView: FeedView;
};

export function ProfileFeed({
  posts,
  likedPosts,
  comments,
  activeView,
}: Props) {
  if (activeView === "comments") {
    return <ProfileCommentsFeed comments={comments} />;
  }
  if (activeView === "likes") {
    return <ProfileLikesFeed likedPosts={likedPosts} />;
  }
  return <ProfilePostsFeed posts={posts} />;
}
