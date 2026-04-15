import { cn } from "@/lib/utils";

interface TagProps {
  label: string;
  className?: string;
}

export function Tag({ label, className }: TagProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-lg border-1 border-[#D4E94B] bg-[#FCFDF1] px-3 py-1 text-xs text-zinc-800",
        className,
      )}
    >
      {label}
    </span>
  );
}
