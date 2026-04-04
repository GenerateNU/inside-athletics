"use client";

import { Heart, Bookmark } from "lucide-react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import type { CommentResponse, PostResponse } from "@/api/models";

type Props = {
  posts: PostResponse[];
  likedPosts: PostResponse[];
  comments: CommentResponse[];
  activeView: "posts" | "comments" | "likes";
};

function FeedPost({
  post,
  filledHeart,
}: {
  post: PostResponse;
  filledHeart: boolean;
}) {
  return (
    <article className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
      <Avatar className="h-9 w-9 bg-slate-300">
        <AvatarFallback />
      </Avatar>
      <div>
        <div className="mb-1 flex items-baseline gap-3">
          <p className="text-sm font-semibold">
            {post.author.username || "username"}
          </p>
          <p className="text-sm leading-6 text-slate-700">{post.content}</p>
        </div>
        <div className="flex flex-wrap gap-2">
          <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
            <Bookmark className="mr-1 h-3 w-3" />
            {post.college.name}
          </Badge>
          <Badge variant="outline" className="rounded-sm px-2 py-1 text-[11px]">
            {post.sport.name}
          </Badge>
        </div>
      </div>
      <button
        type="button"
        className="mt-1 text-[#2f6aa0] transition-opacity hover:opacity-70"
      >
        <Heart className={cn("h-5 w-5", filledHeart && "fill-current")} />
      </button>
    </article>
  );
}

function FeedComment({ comment }: { comment: CommentResponse }) {
  return (
    <article className="space-y-2">
      <div className="grid grid-cols-[auto_1fr_auto] items-start gap-3">
        <Avatar className="h-9 w-9 bg-slate-300">
          <AvatarFallback />
        </Avatar>
        <div className="mb-1 flex items-baseline gap-3">
          <p className="text-sm font-semibold">
            {comment.user?.username || "otherperson"}
          </p>
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

export function ProfileFeed({
  posts,
  likedPosts,
  comments,
  activeView,
}: Props) {
  if (activeView === "comments") {
    return (
      <div className="mt-4 space-y-6">
        {comments.map((comment) => (
          <FeedComment key={comment.id} comment={comment} />
        ))}
      </div>
    );
  }

  const rows = activeView === "likes" ? likedPosts : posts;

  return (
    <div className="mt-4 space-y-6">
      {rows.map((post) => (
        <FeedPost
          key={post.id}
          post={post}
          filledHeart={activeView === "likes" || Boolean(post.is_liked)}
        />
      ))}
    </div>
  );
}
