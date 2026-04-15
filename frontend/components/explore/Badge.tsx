import { cn } from "@/lib/utils";

interface BadgeProps {
  icon: React.ReactNode;
  count: number;
  onClick?: () => void;
  active?: boolean;
  className?: string;
}

export function Badge({ icon, count, onClick, active, className }: BadgeProps) {
  return (
    <button
      type="button"
      onClick={onClick}
      disabled={!onClick}
      className={cn(
        "flex items-center gap-5 rounded-full border border-[#3E7DBB] bg-white px-2 w-17 py-2 text-xs text-zinc-500 transition-colors text-md",
        onClick && "hover:border-red-200 hover:bg-red-50 hover:text-red-500",
        !onClick && "cursor-default",
        active && "text-red-500",
        className,
      )}
    >
      {icon}
      <span>{count}</span>
    </button>
  );
}
