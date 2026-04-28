"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState, useRef } from "react";
import { Plus, X, Paperclip } from "lucide-react";
import { Textarea } from "@/components/ui/textarea";
import Select from "react-select";
import SearchPopup from "./search-popup.tsx";
import { useSession } from "@/utils/SessionContext";
import {
  useGetApiV1Colleges,
  usePostApiV1PostPremium,
  usePostApiV1ContentUploadUrl,
  usePostApiV1ContentConfirmUpload,
  usePostApiV1Media,
} from "@/api/hooks";
import { Tag, TagButton } from "./tag-button.tsx";
import { errorModelSchema } from "@/api/zod/errorModelSchema";

interface CreatePremiumPostPopupProps {
  onClose: () => void;
}

export default function CreatePremiumPostPopup({ onClose }: CreatePremiumPostPopupProps) {
  const [activeTags, setActiveTags] = useState<Tag[]>([]);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [selectedSchool, setSelectedSchool] = useState<{ value: string; label: string } | null>(null);
  const [showSearchPopup, setShowSearchPopup] = useState(false);
  const [file, setFile] = useState<File | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const session = useSession();
  const enabled = !!session?.access_token;
  const authHeaders = session?.access_token
    ? { Authorization: `Bearer ${session.access_token}` }
    : undefined;

  const { data: collegesData } = useGetApiV1Colleges(undefined, {
    query: { enabled },
    client: { headers: authHeaders },
  });

  const { mutateAsync: getUploadUrl } = usePostApiV1ContentUploadUrl({
    client: { headers: authHeaders },
  });
  const { mutateAsync: confirmUpload } = usePostApiV1ContentConfirmUpload({
    client: { headers: authHeaders },
  });
  const { mutateAsync: createMedia } = usePostApiV1Media({
    client: { headers: authHeaders },
  });
  const { mutateAsync: createPremiumPost } = usePostApiV1PostPremium({
    client: { headers: authHeaders },
  });

  const handleSubmit = async () => {
    setError(null);

    const schoolTag = activeTags.find((t) => t.type === "schools");
    const sportTag = activeTags.find((t) => t.type === "sports");
    const otherTags = activeTags.filter((t) => t.type !== "schools" && t.type !== "sports");

    if (!title.trim()) {
      setError("Please enter a title for your post.");
      return;
    }
    if (!content.trim()) {
      setError("Please enter a message for your post.");
      return;
    }
    if (!schoolTag && !sportTag && otherTags.length === 0) {
      setError("Please add at least one school, sport, or tag.");
      return;
    }

    setIsSubmitting(true);
    try {
      let mediaId: string | undefined;

      if (file) {
        const s3Key = `premium/content/${Date.now()}/${file.name}`;
        const fileType = file.type || "application/octet-stream";

        const uploadUrlData = await getUploadUrl({
          data: { key: s3Key, fileType, fileName: file.name },
        });

        const s3PutResponse = await fetch(uploadUrlData.upload_url, {
          method: "PUT",
          body: file,
          headers: {
            "Content-Type": fileType,
          },
        });

        if (!s3PutResponse.ok) {
          throw new Error("Failed to upload file to storage.");
        }

        await confirmUpload({ data: { key: uploadUrlData.key } });

        const mediaData = await createMedia({
          data: { s3key: uploadUrlData.key, media_type: fileType, title: file.name },
        });

        mediaId = mediaData.id;
      }

      await createPremiumPost({
        data: {
          title,
          content,
          tag: otherTags.map((t) => t.id),
          ...(schoolTag ? { college_id: schoolTag.id } : {}),
          ...(sportTag ? { sport_id: sportTag.id } : {}),
          ...(mediaId ? { media_id: mediaId } : {}),
        } as any,
      });

      onClose();
    } catch (err: unknown) {
      const responseData = (err as { response?: { data?: unknown } })?.response?.data;
      const apiError = errorModelSchema.safeParse(responseData);
      if (apiError.success) {
        const firstDetail = apiError.data.errors?.[0]?.message;
        setError(firstDetail ?? apiError.data.detail ?? apiError.data.title ?? "Failed to create post.");
      } else if (err instanceof Error) {
        setError(err.message);
      } else {
        setError("Failed to create post. Please try again.");
      }
    } finally {
      setIsSubmitting(false);
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
    <div className="flex w-[60vw] bg-white rounded-2xl justify-center items-center py-10 px-10 overflow-scroll">
      <div className="w-full flex flex-col justify-between space-y-3 overflow-y-auto max-h-[90vh]">
        <div className="flex justify-between">
          <label className="block text-3xl text-[#001225] font-bold">
            Create Premium Post
          </label>
          <Button variant="ghost" onClick={onClose}>
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
        />
        <label className="block text-1xl text-[#001225] font-bold">Message</label>
        <Textarea
          placeholder="Body Text"
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="min-h-[200px] text-[#000000]"
        />
        <label className="block text-1xl text-[#001225] font-bold">Add School</label>
        <label className="block text-xs text-[#001225]">Select ONE School</label>
        <Select
          instanceId="school-select-premium"
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
            <TagButton tag={firstSchoolTag} active={true} onClick={() => removeTag(firstSchoolTag)} />
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
            {activeTags
              .filter((t) => t.type !== "schools")
              .map((tag) => (
                <TagButton key={tag.id} tag={tag} active={true} onClick={() => removeTag(tag)} />
              ))}
          </div>
        </div>
        <label className="block text-1xl text-[#001225] font-bold">Attach File</label>
        <div
          className="flex items-center gap-3 border border-dashed border-[#2C649A] rounded-lg p-3 cursor-pointer hover:bg-[#F4F8FA] transition-colors"
          onClick={() => fileInputRef.current?.click()}
        >
          <Paperclip size={18} className="text-[#2C649A] shrink-0" />
          <span className="text-sm text-[#001225]">
            {file ? file.name : "Click to attach a file (image, video, PDF…)"}
          </span>
          {file && (
            <button
              className="ml-auto text-gray-400 hover:text-red-500"
              onClick={(e) => { e.stopPropagation(); setFile(null); }}
            >
              <X size={16} />
            </button>
          )}
        </div>
        <input
          ref={fileInputRef}
          type="file"
          className="hidden"
          accept="image/*,video/*,application/pdf"
          onChange={(e) => setFile(e.target.files?.[0] ?? null)}
        />
        {error && <p className="text-sm text-red-500 text-center">{error}</p>}
        <div className="flex justify-end">
          <Button
            variant="ghost"
            className="rounded-full bg-[#2C649A] text-[#F4F8FA] hover:text-[#F4F8FA] hover:bg-[#245580] flex items-center gap-2 h-full px-4 py-1"
            onClick={handleSubmit}
            disabled={isSubmitting}
          >
            {isSubmitting ? "Posting…" : "Post"}
          </Button>
        </div>
      </div>
    </div>
  );
}
