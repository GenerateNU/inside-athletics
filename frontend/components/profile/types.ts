import type { CommentResponse, PostResponse } from "@/api/models";

export type ProfileType = "athlete" | "regular";
export type FeedView = "posts" | "comments" | "likes";

export type ProfilePageData = {
  profileType: ProfileType;
  showSurveyPrompt: boolean;
  user: {
    id: string;
    username: string;
    firstName: string;
    lastName: string;
    pronouns: string;
    email?: string;
    about: string;
    divisionTag?: string;
    sportTag?: string;
    collegeTag?: string;
  };
  communities: string[];
  interests: string[];
  posts: PostResponse[];
  likedPosts: PostResponse[];
  comments: CommentResponse[];
};
