"use client";

import Link from "next/link";
import { use, useState, useEffect } from "react";
import { useSession } from "@/utils/SessionContext";
import {
  useGetApiV1CollegeById,
  useGetApiV1CollegesSearch,
  useGetApiV1UserCollegeFollows,
  usePostApiV1UserCollege,
  useDeleteApiV1UserCollegeById,
} from "@/api/hooks";
import { Navbar } from "@/components/ui/navbar";
import { SearchBar } from "@/components/post/SearchBar";
import Image from "next/image";
import { Plus, Check } from "lucide-react";


function CircleStat({ label, value, max, unit }: { label: string; value: number | null; max: number; unit?: string }) {
  const pct = value != null ? Math.min(value / max, 1) : 0;
  const r = 36;
  const circ = 2 * Math.PI * r;
  const dash = pct * circ;

  return (
    <div className="flex flex-col items-center gap-2">
      <div className="relative w-24 h-24">
        <svg className="w-full h-full -rotate-90" viewBox="0 0 88 88">
          <circle cx="44" cy="44" r={r} fill="none" stroke="#e5e7eb" strokeWidth="8" />
          <circle
            cx="44" cy="44" r={r} fill="none"
            stroke="#3b82f6" strokeWidth="8"
            strokeDasharray={`${dash} ${circ}`}
            strokeLinecap="round"
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center text-sm font-semibold text-black">
          {value != null ? `${value}${unit ?? ""}` : "N/A"}
        </div>
      </div>
      <span className="text-xs text-zinc-500">{label}</span>
    </div>
  );
}

export default function CollegePage({
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

  const [search, setSearch] = useState("");
  const [description, setDescription] = useState<string | null>(null);
  const [followId, setFollowId] = useState<string | null>(null);

  const { data: college, isLoading, error } = useGetApiV1CollegeById(id, {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { data: collegeFollows, refetch: refetchFollows } = useGetApiV1UserCollegeFollows({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const isFollowing = collegeFollows?.college_ids?.includes(id) ?? false;

  const { mutate: followCollege } = usePostApiV1UserCollege({
    client: { headers: authHeaders },
  });

  const { mutate: unfollowCollege } = useDeleteApiV1UserCollegeById({
    client: { headers: authHeaders },
  });

  function handleFollowToggle() {
    if (isFollowing) {
        if (!followId) {
        // Can't unfollow without the follow ID — this is an API limitation
        // Would need backend to return follow ID in the college follows endpoint
        console.warn("No follow ID available to unfollow");
        return;
        }
        unfollowCollege(
        { id: followId },
        { onSuccess: () => { setFollowId(null); refetchFollows(); } }
        );
    } else {
        followCollege(
        { data: { college_id: id } },
        { onSuccess: (data) => { setFollowId(data.id); refetchFollows(); } }
        );
    }
  }

  useEffect(() => {
    if (!college?.name) return;
    const slug = college.name.replace(/\s+/g, "_");
    fetch(`https://en.wikipedia.org/api/rest_v1/page/summary/${slug}`)
      .then((res) => res.json())
      .then((data) => setDescription(data.extract ?? null))
      .catch(() => setDescription(null));
  }, [college?.name]);

  const { data: searchResults } = useGetApiV1CollegesSearch(
    { search_str: search },
    {
      query: { enabled: enabled && search.length > 0 },
      client: { headers: authHeaders },
    }
  );

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error loading college</div>;

  return (
    <div className="flex min-h-screen bg-white bg-linear-to-b from-[#A8C8E8]/60 to-[#E8F1FA]/60">
      <Navbar className="sticky top-0 h-screen shrink-0" />

      <div className="flex min-w-0 w-full pt-10 px-10 flex-col gap-4 m-10 rounded-4xl">
        {/* Search */}
        <div className="relative">
          <SearchBar value={search} onChange={setSearch} placeholder="Search colleges..." />
          {search.length > 0 && searchResults && (
            <ul className="absolute z-50 top-full left-0 right-0 mt-1 bg-white rounded-2xl shadow-lg border border-zinc-100 max-h-64 overflow-y-auto">
              {searchResults.results?.map((c) => (
                <li key={c.id} className="border-b border-zinc-100 last:border-0">
                  <Link href={`/college/${c.id}`} className="block py-2 px-4 text-black hover:bg-zinc-50">
                    {c.name}
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </div>

        {/* Main college module */}
        <main className="flex min-w-0 w-full p-10 flex-col bg-white rounded-4xl">
          {/* Banner + logo */}
          <div className="relative w-full rounded-2xl overflow-visible mb-6">
            <div className="w-full h-20 rounded-2xl bg-linear-to-r from-blue-500 to-teal-400" />
            <div className="absolute -bottom-8 left-6">
              {college?.logo ? (
                <Image
                  src={college.logo}
                  alt={college.name}
                  width={80}
                  height={80}
                  className="rounded-md shadow-lg border-2 p-2 border-white bg-gray-100 object-contain"
                />
              ) : (
                <div className="w-20 h-20 rounded-md shadow-lg border-2 border-white bg-gray-100 flex items-center justify-center text-zinc-400 text-xs font-medium">
                  No logo
                </div>
              )}
            </div>
          </div>

          {/* Name + follow + metadata */}
          <div className="mt-5 flex items-start justify-between">
            <div className="flex-1">
              <h1 className="text-xl font-semibold text-black">{college?.name}</h1>
              <div className="flex gap-4 mt-1 text-xs text-zinc-400">
                {college?.city && college?.state && (
                  <span>{college.city}, {college.state}</span>
                )}
                {college?.website && (
                  <a href={college.website} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                    Website
                  </a>
                )}
              </div>
              <p className="mt-2 text-sm text-zinc-500 leading-relaxed">
                {description ?? "No description available."}
              </p>
            </div>

            {/* Follow button */}
            <button
            onClick={handleFollowToggle}
            className={`ml-6 px-4 py-1.5 cursor-pointer rounded-full text-sm font-medium transition-colors flex items-center gap-1 ${
                isFollowing
                ? "bg-zinc-100 text-zinc-600 hover:bg-zinc-200"
                : "bg-blue-500 text-white hover:bg-blue-600"
            }`}
            >
            {isFollowing ? (
                <><Check className="size-4" /> Following</>
            ) : (
                <><Plus className="size-4" /> Follow</>
            )}
            </button>
          </div>
        </main>

        {/* Stats module */}
        <main className="w-full min-w-0 p-10 flex flex-col bg-white rounded-4xl">
          <h2 className="text-lg font-semibold text-black mb-6">Rankings</h2>
          <div className="flex gap-10">
            <CircleStat
              label="Academic Rank"
              value={college?.academic_rank ?? null}
              max={500}
            />
            <CircleStat
              label="NCAA Division"
              value={college?.division_rank ?? null}
              max={3}
            />
          </div>
        </main>
      </div>
    </div>
  );
}