"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useEffect, useState } from "react";
import { Plus, X } from "lucide-react";
import { Textarea } from "@/components/ui/textarea";
import Select from "react-select";
import SearchPopup from "./search-popup.tsx";
import { useSession } from "@/utils/SessionContext";
import { useGetApiV1Colleges, usePostApiV1Post } from "@/api/hooks";
import { useRouter } from "next/navigation";
import { Tag, TagButton } from "./tag-button.tsx";
import { createPostRequestSchema } from "@/api/zod/createPostRequestSchema";
import { errorModelSchema } from "@/api/zod/errorModelSchema";

export default function CreatePostPopup() {
  const [isAdmin, setIsAdmin] = useState(false);
  const [_, setLoading] = useState(true);
  const [activeTags, setActiveTags] = useState<Tag[]>([]);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [isInsiderContent, setIsInsiderContent] = useState(false);
  const [selectedSchool, setSelectedSchool] = useState<{ value: string; label: string } | null>(null);
  const [showSearchPopup, setShowSearchPopup] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  useEffect(() => {
    fetch("/api/v1/role/roles", {
      headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
    })
      .then((r) => r.json())
      .then((data) => {
        setIsAdmin(data.roles.some((r: { name: string }) => r.name === "admin"));
      })
      .finally(() => setLoading(false));
  }, []);

  const { data: collegesData } = useGetApiV1Colleges({
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { mutateAsync: createPost } = usePostApiV1Post({
    client: { headers: authHeaders },
  });

  const handleSubmit = async () => {
    setError(null);

    const schoolTag = activeTags.find((t) => t.type === "schools");
    const sportTag = activeTags.find((t) => t.type === "sports");
    const otherTags = activeTags.filter((t) => t.type !== "schools" && t.type !== "sports");

    const payload = {
      title,
      content,
      is_anonymous: isAnonymous,
      is_premium: isInsiderContent,
      college_id: schoolTag?.id,
      sport_id: sportTag?.id,
      tags: otherTags.length > 0 ? otherTags.map((t) => ({ id: t.id, name: t.name })) : null,
    };

    const parsed = createPostRequestSchema.safeParse(payload);
    if (!parsed.success) {
      const firstIssue = parsed.error.issues[0];
      const field = firstIssue?.path[0];
      const fieldMessages: Record<string, string> = {
        title: "Please enter a title for your post.",
        content: "Please enter a message for your post.",
      };
      setError(fieldMessages[field as string] ?? "Please fill out all required fields.");
      return;
    }

    try {
      await createPost({ data: parsed.data });
      router.push("/");
    } catch (err: unknown) {
      const responseData = (err as { response?: { data?: unknown } })?.response?.data;
      const apiError = errorModelSchema.safeParse(responseData);
      if (apiError.success) {
        setError(apiError.data.detail ?? apiError.data.title ?? "Failed to create post.");
      } else {
        setError("Failed to create post. Please try again.");
      }
    }
  };

  const schools = (collegesData?.colleges ?? []).map((c) => ({
    value: c.id,
    label: c.name,
  }));

  const toggleSchoolTag = (option: { value: string; label: string } | null) => {
    if (!option) return;
    const tag: Tag = { id: option.value, name: option.label, type: "schools" };
    setActiveTags((prev) => {
      if (prev.find((t) => t.id === tag.id)) return prev.filter((t) => t.id !== tag.id);
      return [...prev.filter((t) => t.type !== "schools"), tag];
    });
    setSelectedSchool(null);
  };

  const removeTag = (tag: Tag) => {
    setActiveTags((prev) => prev.filter((t) => t.id !== tag.id));
  };

  const firstSchoolTag = activeTags.find((t) => t.type === "schools");

  if (showSearchPopup) {
    return (
      <SearchPopup
        activeTags={activeTags}
        setActiveTagsAction={setActiveTags}
        onBackAction={() => setShowSearchPopup(false)}
      />
    );
  }

  return (
    <div className="flex w-[60vw] bg-white rounded-2xl justify-center items-center py-10 px-10 max-h-[90vh]">
      <div className="w-full flex flex-col justify-between space-y-3">
        <div className="flex justify-between items-center">
          <label className="block text-3xl text-[#001225] font-bold">Create Post</label>
          <Button variant="ghost" onClick={() => router.push("/?createPost=false")}>
            <X className="!w-8 !h-8" />
          </Button>
        </div>
        <hr className="border-t border-[#001F3E]" />
        <label className="block text-1xl text-[#001225] font-bold">New Post Title</label>
        <Input
          type="text"
          placeholder="New Post Title"
          value={title}
          className="block text-1xl text-[#000000]"
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        <label className="block text-1xl text-[#001225] font-bold">Message</label>
        {isAdmin && (
          <div className="flex items-center gap-1">
            <label className="text-xs font-bold text-[#001225]">Mark as Insider Content</label>
            <button
              onClick={() => setIsInsiderContent(!isInsiderContent)}
              className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors ${isInsiderContent ? "bg-[#2C649A]" : "bg-gray-300"}`}
            >
              <span
                className={`inline-block h-3 w-3 transform rounded-full bg-white transition-transform ${isInsiderContent ? "translate-x-5" : "translate-x-1"}`}
              />
            </button>
          </div>
        )}
        <Textarea
          placeholder="Body Text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="min-h-[200px] text-[#000000]"
          required
        />
        <label className="block text-1xl text-[#001225] font-bold">Add School</label>
        <label className="block text-xs text-[#001225]">Select ONE School</label>
        <Select
          instanceId="school-select"
          options={schools}
          value={selectedSchool}
          onChange={toggleSchoolTag}
          isSearchable={true}
          placeholder="Select a school..."
          styles={{
            control: (base, state) => ({
              ...base,
              fontFamily: "inherit",
              fontSize: "0.875rem",
              borderColor: state.isFocused ? "#2C649A" : base.borderColor,
              boxShadow: state.isFocused ? "0 0 0 1px #2C649A" : base.boxShadow,
              "&:hover": { borderColor: state.isFocused ? "#2C649A" : base.borderColor },
            }),
            menu: (base) => ({ ...base, fontFamily: "inherit", fontSize: "0.875rem" }),
            option: (base) => ({ ...base, fontFamily: "inherit", fontSize: "0.875rem" }),
            placeholder: (base) => ({ ...base, fontSize: "0.875rem" }),
          }}
        />
        {firstSchoolTag && (
          <div className="flex flex-wrap gap-2">
            <TagButton
              key={firstSchoolTag.id}
              tag={firstSchoolTag}
              active={true}
              showAdminView={false}
              onClick={() => removeTag(firstSchoolTag)}
            />
          </div>
        )}
        <div className="flex flex-wrap gap-2 mt-2">
          <Button
            variant="ghost"
            onClick={() => setShowSearchPopup(true)}
            className="rounded-lg border border-[#D4E94B] bg-[#FCFDF1] flex font-light items-center gap-2 px-1 py-1"
          >
            <Plus size={16} />
            Add Tag
          </Button>
          <div className="flex flex-wrap gap-2">
            {activeTags.filter((t) => t.type !== "schools").map((tag) => (
              <TagButton
                key={tag.id}
                tag={tag}
                active={true}
                showAdminView={false}
                onClick={() => removeTag(tag)}
              />
            ))}
          </div>
        </div>
        {error && <p className="text-sm text-red-500 text-center">{error}</p>}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1">
            <label className="text-xs font-bold text-[#001225]">Post Anonymously</label>
            <button
              onClick={() => setIsAnonymous(!isAnonymous)}
              className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors ${isAnonymous ? "bg-[#2C649A]" : "bg-gray-300"}`}
            >
              <span
                className={`inline-block h-3 w-3 transform rounded-full bg-white transition-transform ${isAnonymous ? "translate-x-5" : "translate-x-1"}`}
              />
            </button>
          </div>
          <Button
            variant="ghost"
            className="rounded-full bg-[#2C649A] text-[#F4F8FA] hover:text-[#F4F8FA] hover:bg-[#245580] flex items-center gap-2 h-full px-4 py-1"
            onClick={handleSubmit}
          >
            Post
          </Button>
        </div>
      </div>
    </div>
  );
}