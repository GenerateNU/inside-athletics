"use client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { useState } from "react";
import { Plus, X, Search, ArrowLeft } from "lucide-react";
import { useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";

//Component for an individual tag
function TagButton({ tag, active, onClick }: { tag: Tag; active: boolean; onClick: () => void }) {
  return (
    <div className={`p-[0.5px] rounded-md ${active ? "bg-gradient-to-b from-[#377DC0] to-[#DBE64C]" : "bg-gray-300"}`}>
      <Button
        variant="ghost"
        onClick={onClick}
        className={`rounded-md bg-white flex items-center gap-2 w-full h-full px-1 py-1 text-[#001225] ${active ? "text-[#001225] hover:text-[#001225]" : "text-gray-500 hover:text-gray-500"}`}
      >
        {active ? <X size={16} /> : <Plus size={16} />}
        {tag.Name}
      </Button>
    </div>
  );
}


//component for a header and list of corresponding tags
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


export default function SearchPopup() {
  const searchParams = useSearchParams();
  const initialTagNames = searchParams.getAll("tag_name");
  const initialTagIsSchools = searchParams.getAll("tag_is_school").map((val) => val === "true");
  const initialActiveTags: Tag[] = initialTagNames.map((name, i) => ({
    Name: name,
    IsSchool: initialTagIsSchools[i] ?? false,
  }));

  const [activeTags, setActiveTags] = useState<Tag[]>(initialActiveTags);
  const [search, setSearch] = useState("");
  const router = useRouter();

  //filter logic for searchbar (matches based on whether tag label contains string input)
  const filter = (tags: Tag[]) =>
    tags.filter((tag) => tag.Name.toLowerCase().includes(search.toLowerCase()));

  //handler function that takes tags on and off the section beneath the search bar when clicking on header tags
  const toggleTag = (tag: Tag) => {
    setActiveTags((prev) => {
      const exists = prev.find((t) => t.Name === tag.Name);
      if (exists) {
        return prev.filter((t) => t.Name !== tag.Name);
      } else {
        return [...prev, tag];
      }
    });
  };

  //handler function that takes tags off the section beneath the search bar when clicking on search bar tags
  const removeTag = (tag: Tag) => {
    setActiveTags((prev) => prev.filter((name) => name.Name !== tag.Name));
  };

  //mock data
  const schools = [
    { Name: "Northeastern University", IsSchool: true }, 
    { Name: "Northwestern University", IsSchool: true },
    { Name: "Southwestern University", IsSchool: true },
    { Name: "Southeastern University", IsSchool: true },
  ];

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
    { Name: "NCAA Compliance", IsSchool: false },
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
    <div className="flex p-4">
      <div className="max-w-lg w-full space-y-4">
        <div className="flex gap-4">
          {/* onclick function left undefined for later popup logic */}
          <Button
            variant="ghost"
            onClick={() => {
              const params = new URLSearchParams();
              activeTags.forEach((tag) => {
                params.append("tag_name", tag.Name);
                params.append("tag_is_school", String(tag.IsSchool));
              });
              router.push(`/create_post_popup?${params.toString()}`);
            }}
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
        {activeTags.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {[...activeTags].map((tag) => (
              <TagButton key={tag.Name} tag={tag} active={true} onClick={() => removeTag(tag)} />
            ))}
          </div>
        )}

        {[...schools.keys()].length > 0 && (
          <TagSection header="Schools" tags={filter(schools)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {[...sports.keys()].length > 0 && (
          <TagSection header="Sports" tags={filter(sports)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {[...recruitments.keys()].length > 0 && (
          <TagSection header="Recruitment" tags={filter(recruitments)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {[...academics.keys()].length > 0 && (
          <TagSection header="Academics" tags={filter(academics)} activeTags={activeTags} onToggle={toggleTag} />
        )}
        {[...advice.keys()].length > 0 && (
          <TagSection header="Advice" tags={filter(advice)} activeTags={activeTags} onToggle={toggleTag} />
        )}
      </div>
    </div>
  );
}