"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Plus, X } from "lucide-react";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";
import { Textarea } from "@/components/ui/textarea";
import { Settings, Image, Link, Video, BarChart2, File } from "lucide-react";

//Component for an individual tag
function TagButton({ tag, active, onClick }: { tag: Tag; active: boolean; onClick: () => void }) {
  return (
    <div className={`p-[0.5px] rounded-md ${active ? "bg-gradient-to-b from-[#377DC0] to-[#DBE64C]" : "bg-gray-300"}`}>
      <Button
        variant="ghost"
        onClick={onClick}
        className={`rounded-md bg-white flex items-center gap-2 w-full h-full px-1 py-1 text-[#001225] ${active ? "text-[#F4F8FA] hover:text-[#F4F8FA]" : "text-gray-500 hover:text-gray-500"} ${tag.IsSchool ? "bg-gradient-to-b from-[#001F3E] to-[#377DC0]" : "bg-gradient-to-b from-[#164779] to-[#00804D]"}`}
      >
        {active ? <X size={16} /> : <Plus size={16} />}
        {tag.Name}
      </Button>
    </div>
  );
}

type Tag = {
  Name: string;
  IsSchool: boolean;
};

export default function CreatePostPopup() {
  const searchParams = useSearchParams();
  const initialTagNames = searchParams.getAll("tag_name");
  const initialTagIsSchools = searchParams.getAll("tag_is_school").map((val) => val === "true");
  const initialActiveTags: Tag[] = initialTagNames.map((name, i) => ({
    Name: name,
    IsSchool: initialTagIsSchools[i] ?? false,
  }));
  const [activeTags, setActiveTags] = useState<Tag[]>(initialActiveTags);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [isAnonymous, setIsAnonymous] = useState(true);
  const [commentVisibility, setCommentVisibility] = useState(true);
  const [isShareable, setIsShareable] = useState(true);
  const router = useRouter();

  //handler function that takes tags off the top bar
  const removeTag = (tag: Tag) => {
    setActiveTags((prev) => prev.filter((name) => name.Name !== tag.Name));
  };

  return (
    <div className="flex p-4">
      <div className="max-w-lg w-full space-y-4">
        <div className="flex justify-between">
          <label className="block text-3xl text-[#001225] font-bold">Create Post</label>
          <Button
            variant="ghost"
            onClick={() => { }}
          >
            <X className="!w-8 !h-8" />
          </Button>
        </div>
        <hr className="border-t border-[#001F3E]" />
        <Input
          type="text"
          placeholder="New Post Title"
          value={title}
          className="block text-1xl text-[#000000]"
          onChange={(e) => setTitle(e.target.value)}
          required
        />
        {/* Added tags */}
        <div className="flex flex-wrap gap-2 mt-2">
          <div className={`p-[0.5px] rounded-md bg-gray-300`}>
            <Button
              variant="ghost"
              onClick={() => {
                const params = new URLSearchParams();
                activeTags.forEach((tag) => {
                  params.append("tag_name", tag.Name);
                  params.append("tag_is_school", String(tag.IsSchool));
                });
                router.push(`/search_popup?${params.toString()}`);
              }}
              className={`rounded-md bg-white flex items-center gap-2 w-full h-full px-1 py-1 text-[#001225]}`}
            >
              <Plus size={16} />
              Add Tag
            </Button>
          </div>
          <div className="flex flex-wrap gap-2">
            {[...activeTags].map((tag) => (
              <TagButton key={tag.Name} tag={tag} active={true} onClick={() => removeTag(tag)} />
            ))}
          </div>
        </div>
        <Textarea
          placeholder="Begin typing..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="min-h-[200px] text-[#000000]"
          required
        />
        <div className="flex justify-between">
          <div className="flex gap-1">
            <Button variant="ghost"><Settings size={20} /></Button>
            <Button variant="ghost"><Image size={20} /></Button>
            <Button variant="ghost"><Link size={20} /></Button>
            <Button variant="ghost"><Video size={20} /></Button>
            <Button variant="ghost"><BarChart2 size={20} /></Button>
            <Button variant="ghost"><File size={20} /></Button>
            <div className="flex gap-1 items-center">
              <input
                type="checkbox"
                checked={isAnonymous}
                onChange={(e) => setIsAnonymous(e.target.checked)}
                className="accent-[#377DC0]"
              />
              <div className="text-xs">Anonymous</div>
              <input
                type="checkbox"
                checked={commentVisibility}
                onChange={(e) => setCommentVisibility(e.target.checked)}
                className="accent-[#377DC0]"
              />
              <div className="text-xs">Comments</div>
              <input
                type="checkbox"
                checked={isShareable}
                onChange={(e) => setIsShareable(e.target.checked)}
                className="accent-[#377DC0]"
              />
              <div className="text-xs">Shareable</div>
            </div>
          </div>
          <Button
            variant="ghost"
            className={`rounded-md text-gray-400 border-[#001F3E] text-[#001F3E] bg-white flex items-center gap-2 h-full px-1 py-1`}
            onClick={() => { }}
          >
            Post
          </Button>
        </div>
      </div>
    </div>
  );
}