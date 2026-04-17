"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
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
  const [activeTags, setActiveTags] = useState<Tag[]>([]);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [selectedSchool, setSelectedSchool] = useState<{
    value: string;
    label: string;
  } | null>(null);
  const [showSearchPopup, setShowSearchPopup] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const { data: collegesData } = useGetApiV1Colleges(undefined, 
    {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { mutateAsync: createPost, isPending } = usePostApiV1Post({
    client: { headers: authHeaders },
  });

  const handleSubmit = async () => {
    setError(null);

    const schoolTag = activeTags.find((t) => t.type === "schools");
    const sportTag = activeTags.find((t) => t.type === "sports");
    const otherTags = activeTags.filter(
      (t) => t.type !== "schools" && t.type !== "sports",
    );

    const payload = {
      title,
      content,
      is_anonymous: isAnonymous,
      college_id: schoolTag?.id,
      sport_id: sportTag?.id,
      tags:
        otherTags.length > 0
          ? otherTags.map((t) => ({ id: t.id, name: t.name }))
          : null,
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
      const responseData = (err as { response?: { data?: unknown } })?.response
        ?.data;
      const apiError = errorModelSchema.safeParse(responseData);
      if (apiError.success) {
        setError(
          apiError.data.detail ??
            apiError.data.title ??
            "Failed to create post.",
        );
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
      if (prev.find((t) => t.id === tag.id))
        return prev.filter((t) => t.id !== tag.id);
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
    <div className="flex w-[60vw] bg-white rounded-2xl flex justify-center items-center py-10 px-10 max-h-[60vh] overflow-scroll">
      <div className=" w-full flex flex-col justify-between space-y-3 overflow-y-auto max-h-[90vh]">
        <div className="flex justify-between">
          <label className="block text-3xl text-[#001225] font-bold">
            Create Post
          </label>
          <Button
            variant="ghost"
            onClick={() => router.push("/?createPost=false")}
          >
            <X className="!w-8 !h-8" />
          </Button>
        </div>
        <hr className="border-t border-[#001F3E]" />
        <label className="block text-1xl text-[#001225] font-bold">
          New Post Title
        </label>
        <Input
          type="text"
          placeholder="New Post Title"
          value={title}
          className="block text-1xl text-[#000000]"
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        <label className="block text-1xl text-[#001225] font-bold">
          Message
        </label>
        <Textarea
          placeholder="Body Text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="min-h-[200px] text-[#000000]"
          required
        />
        <label className="block text-1xl text-[#001225] font-bold">
          Add School
        </label>
        <label className="block text-xs text-[#001225]">
          Select ONE School
        </label>
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
              "&:hover": {
                borderColor: state.isFocused ? "#2C649A" : base.borderColor,
              },
            }),
            menu: (base) => ({
              ...base,
              fontFamily: "inherit",
              fontSize: "0.875rem",
            }),
            option: (base) => ({
              ...base,
              fontFamily: "inherit",
              fontSize: "0.875rem",
            }),
            placeholder: (base) => ({ ...base, fontSize: "0.875rem" }),
          }}
        />
        {activeTags.filter((t) => t.type === "schools").length > 0 && (
          <div className="flex flex-wrap gap-2">
            {firstSchoolTag && (
              <TagButton
                key={firstSchoolTag.id}
                tag={firstSchoolTag}
                active={true}
                onClick={() => removeTag(firstSchoolTag)}
              />
            )}
          </div>
        )}
        <div className="flex flex-wrap gap-2 mt-2">
          <div className="">
            <Button
              variant="ghost"
              onClick={() => setShowSearchPopup(true)}
              className="rounded-lg border border-[#D4E94B] bg-[#FCFDF1] flex font-light items-center gap-2 w-full h-full px-1 py-1"
            >
              <Plus size={16} />
              Add Tag
            </Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {activeTags
              .filter((t) => t.type !== "schools")
              .map((tag) => (
                <TagButton
                  key={tag.id}
                  tag={tag}
                  active={true}
                  onClick={() => removeTag(tag)}
                />
              ))}
          </div>
        </div>
        {error && <p className="text-sm text-red-500 text-center">{error}</p>}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-1">
            <label className="text-xs font-bold text-[#001225]">
              Post Anonymously
            </label>
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
