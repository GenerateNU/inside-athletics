import { Search } from "lucide-react";
import { Input } from "@/components/ui/input";

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function SearchBar({
  value,
  onChange,
  placeholder = "Search posts...",
  className,
}: SearchBarProps) {
  return (
    <div className="border-n">
      <Search className="pointer-events-none absolute left-9 top-1/2 size-5 -translate-y-1/2 text-zinc-400" />
      <Input
        type="search"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder={placeholder}
        className="w-full h-12 pl-9 border-white placeholder:text-zinc-400 placeholder:text-base bg-white"
        aria-label={placeholder}
      />
    </div>
  );
}
