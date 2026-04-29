"use client";

import { use, useState, useEffect } from "react";
import { ArrowLeft, Heart, MessageCircle } from "lucide-react";
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Navbar } from "@/components/ui/navbar";
import { useSession } from "@/utils/SessionContext";
import { useQueryClient } from "@tanstack/react-query";
import {
  useGetApiV1PostById,
  useListApiV1PostByPostIdComments,
  usePostApiV1PostLike,
  useDeleteApiV1PostLikeById,
  usePostApiV1Comment,
  listApiV1PostByPostIdCommentsQueryKey,
  getApiV1PostByIdQueryKey,
} from "@/api/hooks";
import { CommentCard } from "@/components/post/CommentCard";
import { Badge } from "@/components/post/Badge";
import { Tag } from "@/components/post/Tag";
import { cn } from "@/lib/utils";
import { SearchBar } from "@/components/post/SearchBar";

export default function PostPage({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const { data: post, isLoading: loadingPost } = useGetApiV1PostById(id, {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { data: comments, isLoading: loadingComments } =
    useListApiV1PostByPostIdComments(id, {
      query: { enabled },
      client: { headers: authHeaders },
    });

  const queryClient = useQueryClient();
  const [isLiked, setIsLiked] = useState(false);
  const [likeCount, setLikeCount] = useState(0);
  const [commentOpen, setCommentOpen] = useState(false);
  const [commentText, setCommentText] = useState("");

  const { mutate: submitComment, isPending: submittingComment } =
    usePostApiV1Comment({
      client: { headers: authHeaders },
    });

  function handleCommentSubmit() {
    if (!commentText.trim()) return;
    submitComment(
      {
        data: {
          description: commentText.trim(),
          is_anonymous: false,
          post_id: id,
        },
      },
      {
        onSuccess: () => {
          setCommentText("");
          setCommentOpen(false);
          queryClient.invalidateQueries({ queryKey: listApiV1PostByPostIdCommentsQueryKey(id) });
          queryClient.invalidateQueries({ queryKey: getApiV1PostByIdQueryKey(id) });
        },
      },
    );
  }

  useEffect(() => {
    if (post) {
      setIsLiked(post.is_liked ?? false);
      setLikeCount(post.like_count ?? 0);
    }
  }, [post?.id]);

  const { mutate: likePost } = usePostApiV1PostLike({
    client: { headers: authHeaders },
  });

  const { mutate: unlikePost } = useDeleteApiV1PostLikeById({
    client: { headers: authHeaders },
  });

  function handleLikeToggle() {
    if (isLiked) {
      setIsLiked(false);
      setLikeCount((c) => c - 1);
      unlikePost(
        { id },
        {
          onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: getApiV1PostByIdQueryKey(id) });
          },
          onError: () => {
            setIsLiked(true);
            setLikeCount((c) => c + 1);
          },
        },
      );
    } else {
      setIsLiked(true);
      setLikeCount((c) => c + 1);
      likePost(
        { data: { post_id: id } },
        {
          onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: getApiV1PostByIdQueryKey(id) });
          },
          onError: () => {
            setIsLiked(false);
            setLikeCount((c) => c - 1);
          },
        },
      );
    }
  }

  const authorName = post?.is_anonymous
    ? "Anonymous"
    : post
      ? `${post.author.first_name} ${post.author.last_name}`
      : "";

  const authorInitials = post?.is_anonymous
    ? "A"
    : post
      ? `${post.author.first_name[0] ?? ""}${post.author.last_name[0] ?? ""}`.toUpperCase()
      : "";

  const tagLabels = post
    ? [
        ...(post.college?.name ? [post.college.name] : []),
        ...(post.sport?.name ? [post.sport.name] : []),
        ...(post.tags?.map((t) => t.name) ?? []),
      ]
    : [];

  return (
    <div className="flex min-h-screen bg-white bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60">
      <Navbar className="sticky top-0 h-screen shrink-0" />

      <main className="flex min-w-0 pt-10 px-10 flex-1 flex-col bg-white m-10 rounded-4xl">
        {/* Back + title */}
        <div className="flex items-center gap-3 px-2">
          <Link
            href="/explore"
            className="shrink-0 text-zinc-400 transition-colors hover:text-zinc-700"
            aria-label="Back"
          >
            <ArrowLeft className="size-12" color="black" />
          </Link>
          {loadingPost ? (
            <div className="h-4 w-56 animate-pulse rounded bg-zinc-100" />
          ) : (
            <h1 className="text-xl font-semibold text-black">{post?.title}</h1>
          )}
        </div>

        {loadingPost ? (
          <div className="space-y-3 px-6 py-6">
            <div className="h-3 w-24 animate-pulse rounded bg-zinc-100" />
            <div className="h-16 animate-pulse rounded bg-zinc-100" />
          </div>
        ) : post ? (
          <>
            {/* Post body */}
            <div className="px-15 py-5">
              {/* Author */}
              {!post.is_anonymous && post.author?.id ? (
                <Link
                  href={`/profile/${post.author.id}`}
                  className="mb-3 inline-flex items-center gap-2 rounded-md outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
                >
                  <Avatar size="lg">
                    {post.author.profile_picture ? (
                      <AvatarImage
                        src={post.author.profile_picture}
                        alt=""
                      />
                    ) : null}
                    <AvatarFallback className="bg-zinc-100 text-[10px] font-medium text-zinc-600">
                      {authorInitials}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-lg font-semibold text-black">
                    {authorName}
                  </span>
                </Link>
              ) : (
                <div className="mb-3 flex items-center gap-2">
                  <Avatar size="lg">
                    {!post.is_anonymous && post.author.profile_picture && (
                      <AvatarImage
                        src={post.author.profile_picture}
                        alt={authorName}
                      />
                    )}
                    <AvatarFallback className="bg-zinc-100 text-[10px] font-medium text-zinc-600">
                      {authorInitials}
                    </AvatarFallback>
                  </Avatar>
                  <span className="text-lg font-semibold text-black">
                    {authorName}
                  </span>
                </div>
              )}

              {/* Content */}
              <p className="text-md leading-relaxed text-black">
                {post.content}
              </p>

              {/* Like + comment counts */}
              <div className="mt-4 flex items-center gap-2">
                <Badge
                  icon={
                    <Heart
                      className={cn(
                        "size-5 shrink-0",
                        isLiked
                          ? "fill-red-500 text-red-500"
                          : "text-[#3E7DBB]",
                      )}
                    />
                  }
                  count={likeCount}
                  active={isLiked}
                  onClick={handleLikeToggle}
                />
                <Badge
                  icon={
                    <MessageCircle className="size-5 shrink-0 text-[#3E7DBB]" />
                  }
                  count={post.comment_count ?? 0}
                  onClick={() => setCommentOpen((o) => !o)}
                />
              </div>
            </div>

            {/* Comments */}
            <div className=" border-zinc-100 px-15 pb-8">
              <p className="py-4 text-lg font-semibold text-black">Comments</p>

              {commentOpen && (
                <div className="mb-4 flex gap-2 relative">
                  <textarea
                    value={commentText}
                    onChange={(e) => setCommentText(e.target.value)}
                    placeholder="Write a comment..."
                    rows={4}
                    className="min-h-0 flex-1 resize-none rounded-2xl border border-[#3E7DBB] bg-white px-3 py-2 text-base text-zinc-900 placeholder:text-zinc-400"
                  />
                  <button
                    type="button"
                    onClick={handleCommentSubmit}
                    disabled={!commentText.trim() || submittingComment}
                    className="absolute bottom-2 right-2 rounded-3xl bg-[#A8C8E8] px-4 py-1.5 text-sm font-medium text-[#E8F1FA] disabled:opacity-50"
                  >
                    {submittingComment ? "Posting..." : "Post"}
                  </button>
                </div>
              )}

              {loadingComments ? (
                <div className="space-y-5">
                  {[1, 2, 3].map((i) => (
                    <div key={i} className="space-y-2">
                      <div className="h-3 w-24 animate-pulse rounded bg-zinc-100" />
                      <div className="h-10 animate-pulse rounded bg-zinc-100" />
                    </div>
                  ))}
                </div>
              ) : !comments || comments.length === 0 ? (
                <p className="py-6 text-center text-sm text-zinc-400">
                  No comments yet.
                </p>
              ) : (
                <>
                  <div className="divide-y divide-zinc-100">
                    {comments.map((comment) => (
                      <CommentCard
                        key={comment.id}
                        comment={comment}
                        postId={id}
                        authHeaders={authHeaders}
                      />
                    ))}
                  </div>
                </>
              )}
            </div>
          </>
        ) : (
          <div className="flex items-center justify-center py-16 text-sm text-zinc-400">
            Post not found.
          </div>
        )}
      </main>
    </div>
  );
}
