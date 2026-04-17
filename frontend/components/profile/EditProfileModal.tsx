"use client";

import { PenSquare, UserRound, X } from "lucide-react";
import * as React from "react";

import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";

type Props = {
  open: boolean;
  onClose: () => void;
  isSaving: boolean;
  user: {
    firstName: string;
    lastName: string;
    pronouns: string;
    about: string;
  };
  onSave: (values: {
    firstName: string;
    lastName: string;
    pronouns: string;
    about: string;
    selectedTagIds: string[];
  }) => Promise<void>;
  availableTags: Array<{ id: string; name: string }>;
  selectedTagIds: string[];
};

function TagSection({
  title,
  tags,
  selected,
  onToggle,
}: {
  title: string;
  tags: Array<{ id: string; name: string }>;
  selected: string[];
  onToggle: (tagId: string) => void;
}) {
  return (
    <div>
      <h4 className="text-sm font-semibold">{title}</h4>
      <div className="mt-2 flex flex-wrap gap-2">
        {tags.map((tag) => (
          <button key={tag.id} type="button" onClick={() => onToggle(tag.id)}>
            <Badge
              className={
                selected.includes(tag.id)
                  ? "h-[30px] rounded-[12px] border border-[#7F8C2D] bg-[#D4E94B]/50 px-[8px] py-[5px] text-sm text-slate-900"
                  : "h-[30px] rounded-[12px] border border-[#D4E94B] bg-[#FCFDF1] px-[8px] py-[5px] text-sm text-slate-700"
              }
            >
              {selected.includes(tag.id) ? `x ${tag.name}` : tag.name}
            </Badge>
          </button>
        ))}
      </div>
    </div>
  );
}

const PRONOUN_OPTIONS = ["She/her", "He/him", "They/them", "Other"];

export function EditProfileModal({
  open,
  onClose,
  user,
  isSaving,
  onSave,
  availableTags,
  selectedTagIds,
}: Props) {
  const [firstName, setFirstName] = React.useState(user.firstName);
  const [lastName, setLastName] = React.useState(user.lastName);
  const [pronouns, setPronouns] = React.useState(user.pronouns);
  const [about, setAbout] = React.useState(user.about);
  const [selectedTags, setSelectedTags] = React.useState<string[]>([]);
  const [saveError, setSaveError] = React.useState<string | null>(null);

  React.useEffect(() => {
    if (!open) return;
    setFirstName(user.firstName);
    setLastName(user.lastName);
    setPronouns(user.pronouns);
    setAbout(user.about);
    setSelectedTags(selectedTagIds);
    setSaveError(null);
  }, [open, selectedTagIds, user.about, user.firstName, user.lastName, user.pronouns]);

  if (!open) return null;

  const toggleTag = (tagId: string) => {
    setSelectedTags((current) =>
      current.includes(tagId)
        ? current.filter((item) => item !== tagId)
        : [...current, tagId],
    );
  };

  const handleSave = async () => {
    setSaveError(null);
    try {
      await onSave({
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        pronouns: pronouns.trim(),
        about: about.trim(),
        selectedTagIds: selectedTags,
      });
    } catch {
      setSaveError("Unable to save profile changes. Please try again.");
    }
  };

  const midpoint = Math.ceil(availableTags.length / 2);
  const generalTags = availableTags.slice(0, midpoint);
  const priorityTags = availableTags.slice(midpoint);

  return (
    <div className="fixed inset-0 z-40 bg-black/35 p-6">
      <div className="mx-auto max-h-[90vh] w-full max-w-[980px] overflow-y-auto rounded-3xl bg-[#e9f0f7] p-6 shadow-2xl md:p-8">
        <div className="mb-6 flex items-center justify-between">
          <h2 className="text-5xl font-black">Edit Profile</h2>
          <button
            type="button"
            className="rounded-full p-1 text-slate-600 hover:bg-slate-100"
            onClick={onClose}
            disabled={isSaving}
          >
            <X className="h-8 w-8" />
          </button>
        </div>

        <div className="space-y-6">
          <div className="flex justify-center">
            <div className="relative">
              <Avatar className="h-24 w-24 border-slate-500 text-slate-600">
                <AvatarFallback>
                  <UserRound className="h-12 w-12" />
                </AvatarFallback>
              </Avatar>
              <button
                type="button"
                className="absolute right-0 bottom-0 rounded-full bg-slate-200 p-1"
              >
                <PenSquare className="h-4 w-4 text-slate-700" />
              </button>
            </div>
          </div>

          <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
            <input
              type="text"
              value={firstName}
              onChange={(e) => setFirstName(e.target.value)}
              className="h-12 rounded-xl border border-[#2f74b5] bg-[#f5f7fb] px-3 text-xl"
              placeholder="First name"
            />
            <input
              type="text"
              value={lastName}
              onChange={(e) => setLastName(e.target.value)}
              className="h-12 rounded-xl border border-[#2f74b5] bg-[#f5f7fb] px-3 text-xl"
              placeholder="Last name"
            />
            <select
              value={pronouns}
              onChange={(e) => setPronouns(e.target.value)}
              className="h-12 rounded-xl border border-[#2f74b5] bg-[#f5f7fb] px-3 text-xl md:col-span-2"
            >
              {PRONOUN_OPTIONS.map((option) => (
                <option key={option} value={option}>
                  {option}
                </option>
              ))}
            </select>
            <div className="md:col-span-2">
              <textarea
                value={about}
                onChange={(e) => setAbout(e.target.value)}
                className="min-h-[180px] w-full rounded-xl border border-[#2f74b5] bg-[#f5f7fb] px-3 py-3 text-xl leading-8"
                placeholder="Tell the community about yourself"
              />
            </div>
          </div>

          <TagSection
            title="General"
            tags={generalTags}
            selected={selectedTags}
            onToggle={toggleTag}
          />
          <TagSection
            title="Priorities"
            tags={priorityTags}
            selected={selectedTags}
            onToggle={toggleTag}
          />

          {saveError ? (
            <p className="text-sm text-red-600" role="alert">
              {saveError}
            </p>
          ) : null}

          <div className="flex justify-end pt-2">
            <Button
              className="rounded-xl bg-[#2d6ca6] px-6 py-2 text-xl font-semibold text-white hover:bg-[#235a8a]"
              onClick={handleSave}
              disabled={isSaving}
            >
              {isSaving ? "Saving..." : "Save"}
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
