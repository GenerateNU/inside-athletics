import { Button } from "./button";
import { Plus, X } from "lucide-react";

export type Tag = {
  id: string;
  name: string;
  type: string;
};

export function TagButton({
  tag,
  active,
  onClick,
}: {
  tag: Tag;
  active: boolean;
  onClick: () => void;
}) {
  return (
    <div className="p-[0.5px]">
      <Button
        variant="ghost"
        onClick={onClick}
        className={`rounded-lg border border-[#7F8C2D] flex items-center gap-2 w-full h-full px-1 py-1 font-light ${active ? "bg-[#D4E94B80]" : "bg-[#FCFDF1]"}`}
      >
        {active ? <X size={16} /> : <Plus size={16} />}
        {tag.name}
      </Button>
    </div>
  );
}
