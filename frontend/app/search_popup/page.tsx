"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Plus, X, Search, ArrowLeft } from "lucide-react";

//Component for an individual tag
function TagButton({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
<div className={`p-[1px] rounded-md ${active ? "bg-gradient-to-b from-blue-500 to-yellow-400" : "bg-gray-300"}`}>
  <Button
    variant="ghost"
    onClick={onClick}
    className={`rounded-md bg-white flex items-center gap-2 w-full h-full px-1 py-1 ${active ? "text-black" : "text-gray-500"}`}
  >
    {active ? <X size={16} /> : <Plus size={16} />}
    {label}
  </Button>
</div>
  );
}

//component for a header and list of corresponding tags
function TagSection({ header, tags, activeTags, onToggle }: {
  header: string;
  tags: string[];
  activeTags: Set<string>;
  onToggle: (tag: string) => void;
}) {
  return (
    <div className="mt-8">
      <label className="block text-1xl text-black font-bold">{header}</label>
      <div className="flex flex-wrap gap-2 mt-2">
        {tags.map((tag) => (
          <TagButton key={tag} label={tag} active={activeTags.has(tag)} onClick={() => onToggle(tag)} />
        ))}
      </div>
    </div>
  );
}

export default function SearchPopup() {
  const [activeTags, setActiveTags] = useState<Set<string>>(new Set());
  const [search, setSearch] = useState("");

  //filter logic for searchbar (matches based on whether tag label contains string input)
  const filter = (tags: string[]) =>
    tags.filter((tag) => tag.toLowerCase().includes(search.toLowerCase()));

  //handler function that takes tags on and off the section beneath the search bar when clicking on header tags
  const toggleTag = (tag: string) => {
    setActiveTags((prev) => {
      const next = new Set(prev);
      if (next.has(tag)) next.delete(tag);
      else next.add(tag);
      return next;
    });
  };

  //handler function that takes tags off the section beneath the search bar when clicking on search bar tags
  const removeTag = (tag: string) => {
    setActiveTags((prev) => {
      const next = new Set(prev);
      next.delete(tag);
      return next;
    });
  };

  //tag labels for each section
  const sports = ["Swim", "Basketball", "Football", "Soccer", "Tennis", "Track & Field", "Lacrosse", "Rowing", "Rugby", "Wrestling", "Gymnastics"];
  const recruitments = ["Transfer", "Walk On", "Recruit", "Scouting", "D1 Recruiting", "D2 Recruiting", "D3 Recruiting", "Highlight Film Review", "NCAA Compliance"];
  const academics = ["Application/Testing", "Highlight Film Review", "NCAA Compliance"];
  const advice = ["Parent Advice", "Ask a Recruit", "Decision Help", "Coach/Recruiter Perspective", "Email Guidance"];

  return (
    <div className="flex p-4">
      <div className="max-w-lg w-full space-y-4">
        <div className="flex gap-4">
          {/* onclick function left undefined for later popup logic */}
          <Button
            variant="ghost"
            onClick={() => { }}
          >
            <ArrowLeft className="!w-8 !h-8" />
          </Button>
          <label className="block text-3xl text-black font-bold">Add Tag</label>
        </div>
        {/* Searchbar */}
        <div className="relative">
          <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
          <Input
            id="search_tags"
            name="search_tags"
            type="text"
            placeholder="Search Tags"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            className="pl-9"
          />
        </div>

        {/* Searchbar tags */}
        {activeTags.size > 0 && (
          <div className="flex flex-wrap gap-2">
            {[...activeTags].map((tag) => (
              <TagButton key={tag} label={tag} active={true} onClick={() => removeTag(tag)} />
            ))}
          </div>
        )}

        {/* Tag sections */}
        {filter(sports).length > 0 && (
          <TagSection header="Sports" tags={filter(sports)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {filter(recruitments).length > 0 && (
          <TagSection header="Recruitment" tags={filter(recruitments)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {filter(academics).length > 0 && (
          <TagSection header="Academics" tags={filter(academics)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {filter(advice).length > 0 && (
          <TagSection header="Advice" tags={filter(advice)} activeTags={activeTags} onToggle={toggleTag} />
        )}
      </div>
    </div>
  );
}