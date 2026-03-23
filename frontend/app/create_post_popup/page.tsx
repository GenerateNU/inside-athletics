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
function TagButton({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <div className={`p-[1px] rounded-md ${active ? "bg-gradient-to-b from-blue-500 to-yellow-400" : "bg-gray-300"}`}>
      <Button
        variant="ghost"
        onClick={onClick}
        className={`rounded-md text-gray-400 border-black flex items-center gap-2 h-full px-1 py-1 text-black bg-white`}
      >
        {active ? <X size={16} /> : <Plus size={16} />}
        {label}
      </Button>
    </div>
  );
}

export default function CreatePostPopup() {
  const searchParams = useSearchParams();
  const initialTags = searchParams.getAll("tag");
  const [activeTags, setActiveTags] = useState<Set<string>>(new Set(initialTags));
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const router = useRouter();

  //handler function that takes tags off the section
  const removeTag = (tag: string) => {
    setActiveTags((prev) => {
      const next = new Set(prev);
      next.delete(tag);
      return next;
    });
  };

  return (
    <div className="flex p-4">
      <div className="max-w-lg w-full space-y-4">
        <div className="flex justify-between">
          <label className="block text-3xl text-black font-bold">Create Post</label>
          <Button
            variant="ghost"
            onClick={() => { }}
          >
            <X className="!w-8 !h-8" />
          </Button>
        </div>
        <hr className="border-t border-gray-400" />
        <Input
          type="text"
          placeholder="New Post Title"
          value={title}
          className="block text-1xl text-gray-400"
          onChange={(e) => setTitle(e.target.value)}
        />
        {/* Added tags */}
        <div className="flex flex-wrap gap-2 mt-2">
          <TagButton
            key={"add_tag"}
            label={"Tag"}
            active={false}
            onClick={() => {
              const params = new URLSearchParams();
              [...activeTags].forEach(tag => params.append("tag", tag));
              router.push(`/search_popup?${params.toString()}`);
            }}
          />
          <div className="flex flex-wrap gap-2">
            {[...activeTags].map((tag) => (
              <TagButton key={tag} label={tag} active={true} onClick={() => removeTag(tag)} />
            ))}
          </div>
        </div>
        <Textarea
          placeholder="Begin typing..."
          value={content}
          onChange={(e) => setContent(e.target.value)}
          className="min-h-[200px]"
        />
        <div className="flex justify-between">
          <div className="flex gap-1">
            <Button variant="ghost"><Settings size={20} /></Button>
            <Button variant="ghost"><Image size={20} /></Button>
            <Button variant="ghost"><Link size={20} /></Button>
            <Button variant="ghost"><Video size={20} /></Button>
            <Button variant="ghost"><BarChart2 size={20} /></Button>
            <Button variant="ghost"><File size={20} /></Button>
          </div>
          <Button
            variant="ghost"
            className={`rounded-md text-gray-400 border-black bg-white flex items-center gap-2 h-full px-1 py-1`}
            onClick={() => { }}
          >
            Post
          </Button>
        </div>
      </div>
    </div>
  );
}