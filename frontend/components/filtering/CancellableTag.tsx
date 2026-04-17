import { X } from "lucide-react";
import { cn } from "@/lib/utils";

interface CancellableTagProps {
  label: string;
  onRemove: () => void;
  className?: string;
}

export function CancellableTag({ label, onRemove, className }: CancellableTagProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-lg border border-[#A8C96A] bg-[#D4E896] px-3 py-1 text-xs text-zinc-800",
        className,
      )}
    >
      <button
        onClick={onRemove}
        className="flex items-center justify-center hover:text-zinc-500"
        aria-label={`Remove ${label}`}
      >
        <X className="size-3" />
      </button>
      {label}
    </span>
  );
}