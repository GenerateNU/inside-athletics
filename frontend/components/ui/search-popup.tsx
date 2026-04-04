"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Plus, X, Search } from "lucide-react";

function TagButton({ tag, active, onClick }: { tag: Tag; active: boolean; onClick: () => void }) {
  return (
    <div className={`p-[0.5px] rounded-md bg-[#00804D]`}>
      <Button
        variant="ghost"
        onClick={onClick}
        className={`rounded-md bg-white flex items-center gap-2 w-full h-full px-1 py-1 text-[#001225] ${active ? "text-[#F4F8FA] hover:text-[#F4F8FA] bg-gradient-to-b from-[#00804D] to-[#043D26]" : "text-[#001225] hover:text-[#001225]"}`}
      >
        {active ? <X size={16} /> : <Plus size={16} />}
        {tag.Name}
      </Button>
    </div>
  );
}

function TagSection({ header, tags, activeTags, onToggle }: {
  header: string;
  tags: Tag[];
  activeTags: Tag[];
  onToggle: (tag: Tag) => void;
}) {
  return (
    <div className="mt-8">
      <label className="block text-1xl text-black font-bold">{header}</label>
      <div className="flex flex-wrap gap-2 mt-2">
        {tags.map((tag) => (
          <TagButton key={tag.Name} tag={tag} active={activeTags.some((t) => t.Name === tag.Name)} onClick={() => onToggle(tag)} />
        ))}
      </div>
    </div>
  );
}

type Tag = {
  Name: string;
  IsSchool: boolean;
};

type SearchPopupProps = {
  activeTags: Tag[];
  setActiveTagsAction: React.Dispatch<React.SetStateAction<Tag[]>>;
  onBackAction: () => void;
};

export default function SearchPopup({ activeTags, setActiveTagsAction, onBackAction }: SearchPopupProps) {
  const [search, setSearch] = useState("");

  const filter = (tags: Tag[]) =>
    tags.filter((tag) => tag.Name.toLowerCase().includes(search.toLowerCase()));

  const toggleTag = (tag: Tag) => {
    setActiveTagsAction((prev) => {
      const exists = prev.find((t) => t.Name === tag.Name);
      if (exists) {
        return prev.filter((t) => t.Name !== tag.Name);
      } else {
        return [...prev, tag];
      }
    });
  };

  const removeTag = (tag: Tag) => {
    setActiveTagsAction((prev) => prev.filter((t) => t.Name !== tag.Name));
  };

  const sports = [
    { Name: "Swim", IsSchool: false },
    { Name: "Basketball", IsSchool: false },
    { Name: "Football", IsSchool: false },
    { Name: "Soccer", IsSchool: false },
    { Name: "Tennis", IsSchool: false },
    { Name: "Track & Field", IsSchool: false },
    { Name: "Lacrosse", IsSchool: false },
    { Name: "Rowing", IsSchool: false },
    { Name: "Rugby", IsSchool: false },
    { Name: "Wrestling", IsSchool: false },
    { Name: "Gymnastics", IsSchool: false },
  ];

  const recruitments = [
    { Name: "Transfer", IsSchool: false },
    { Name: "Walk On", IsSchool: false },
    { Name: "Recruit", IsSchool: false },
    { Name: "Scouting", IsSchool: false },
    { Name: "D1 Recruiting", IsSchool: false },
    { Name: "D2 Recruiting", IsSchool: false },
    { Name: "D3 Recruiting", IsSchool: false },
    { Name: "Highlight Film Review", IsSchool: false },
  ];

  const academics = [
    { Name: "Application/Testing", IsSchool: false },
    { Name: "Highlight Film Review", IsSchool: false },
    { Name: "NCAA Compliance", IsSchool: false },
  ];

  const advice = [
    { Name: "Parent Advice", IsSchool: false },
    { Name: "Ask a Recruit", IsSchool: false },
    { Name: "Decision Help", IsSchool: false },
    { Name: "Coach/Recruiter Perspective", IsSchool: false },
    { Name: "Email Guidance", IsSchool: false },
  ];

  return (
    <div className="flex w-[600px]">
      <div className="max-w-lg w-full space-y-3 justify-between">
        <div className="flex gap-4">
          <label className="block text-1xl text-black font-bold">Add Tag</label>
        </div>
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
        {activeTags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {activeTags.map((tag) => (
              <TagButton key={tag.Name} tag={tag} active={true} onClick={() => removeTag(tag)} />
            ))}
          </div>
        )}
        <TagSection header="Sports" tags={filter(sports)} activeTags={activeTags} onToggle={toggleTag} />
        <TagSection header="Recruitment" tags={filter(recruitments)} activeTags={activeTags} onToggle={toggleTag} />
        <TagSection header="Academics" tags={filter(academics)} activeTags={activeTags} onToggle={toggleTag} />
        <TagSection header="Advice" tags={filter(advice)} activeTags={activeTags} onToggle={toggleTag} />
        <div className="flex justify-end">
          <Button
            variant="ghost"
            className="rounded-full bg-[#2C649A] hover:bg-[#245580] hover:text-[#F4F8FA] text-[#F4F8FA] flex items-center gap-2 h-full px-4 py-1"
            onClick={onBackAction}
          >
            Done
          </Button>
        </div>
      </div>
    </div>
  );
}