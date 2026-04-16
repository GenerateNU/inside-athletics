"use client";

import { useState } from "react";
import { useSearchParams } from "next/navigation";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Button } from "@/components/ui/button";
import { Navbar } from "@/components/ui/navbar";
import CreatePostPopup from "@/components/ui/create-post-popup";
import { RatingPanel } from "@/components/ui/rating-panel";
import { CiUser } from "react-icons/ci";
import SmallPost from "@/components/post/SmallPost";
import type { PostResponse } from "@/api/models/PostResponse";
import PremiumSmallPost from "@/components/post/PremiumSmallPost";

const examplePost: PostResponse = {
  id: "example-post-1",
  title: "My Experience on the Track Team",
  content: "Being on the track team has been one of the best decisions I've made at Northeastern. The coaches are incredibly supportive and the team culture is amazing.",
  is_anonymous: false,
  is_verified_athlete: true,
  like_count: 24,
  comment_count: 7,
  sport: {
    id: "sport-1",
    name: "Track & Field",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
  },
  college: {
    id: "college-1",
    name: "Northeastern University",
    city: "Boston",
    state: "MA",
    academic_rank: 49,
    division_rank: 1,
    logo: "",
    website: "https://northeastern.edu",
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z",
    deleted_at: null,
  },
  author: {
    id: "ac26eeb9-2877-425f-a1cd-887a9c979578",
    first_name: "Zainab",
    last_name: "Imadulla",
    profile_picture: "IMG_1363.JPG",
    username: "nubslovesdubs",
    email: "Zainab.imadulla@icloud.com",
    bio: null,
    college: "",
    sport: "",
    account_type: true,
    division: null,
    expected_grad_year: 0,
    verified_athelete_status: "False",
    created_at: "",
    updated_at: "",
    deleted_at: null,
  },
  tags: null,
};

const northeasternCollegeId = "014d2c09-4023-445d-9779-66aff4824245";

export default function Page() {
  const [selectedCollegeId, setSelectedCollegeId] = useState(
    northeasternCollegeId,
  );
  const searchParams = useSearchParams();
  const showCreatePost = searchParams.get("createPost") === "true";

  return (
    <div className="min-h-screen bg-zinc-50">
      {showCreatePost && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
          <CreatePostPopup />
        </div>
      )}
      <div className="flex min-h-screen">
        <Navbar className="h-screen shrink-0" />
        <main className="flex min-w-0 flex-1 justify-center p-6 md:p-10 overflow-scroll max-h-screen">
          <div className="flex w-full max-w-5xl flex-col items-center gap-10">
            <div className="flex flex-col items-center gap-4 text-center">
              <p className="text-6xl">🐸</p>
              <h1 className="text-4xl font-bold">Welcome to Inside Athletics</h1>
              <p className="text-muted-foreground">
                Under Construction! but here are some components:
              </p>

              <div className="flex flex-row gap-5">
                <Avatar />
                <Avatar>
                  <AvatarFallback>
                    <CiUser strokeWidth={1.3} />
                  </AvatarFallback>
                </Avatar>
              </div>
              <div className="flex flex-row gap-5">
                <Button variant="outline"> Click here</Button>
              </div>
            </div>
            <div className="flex w-full max-w-3xl flex-col items-start gap-3">
              <div className="flex flex-wrap items-center gap-3">
                <Button
                  type="button"
                  onClick={() => setSelectedCollegeId(northeasternCollegeId)}
                >
                  Use Northeastern Ratings
                </Button>
              </div>
              <p className="text-sm text-zinc-600">
                Active college: Northeastern
              </p>
            </div>
            <RatingPanel collegeId={selectedCollegeId} />
            <div className="flex w-full max-w-3xl flex-col gap-3">
              <h2 className="text-lg font-semibold">SmallPost Example</h2>
              <SmallPost post={examplePost} />
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}
