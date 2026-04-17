export default function Loading({ lines = 6 }: { lines?: number }) {
  const widths = ["w-full", "w-5/6", "w-4/6"];
  return (
    <div className="flex h-full w-full flex-col space-y-3">
      {Array.from({ length: lines }).map((_, index) => {
        const width = widths[index % widths.length];
        return (
          <div
            key={index}
            className={`h-5 rounded-full bg-gray-200 animate-pulse ${width}`}
          />
        );
      })}
    </div>
  );
}

export function LargeLoading() {
  return (
    <div className="flex h-full w-full flex-col justify-between">
      <div className="space-y-3">
        <div className="h-5 rounded-full bg-gray-200 animate-pulse" />
        <div className="h-5 w-5/6 rounded-full bg-gray-200 animate-pulse" />
        <div className="h-5 w-4/6 rounded-full bg-gray-200 animate-pulse" />
      </div>

      <div className="flex-1" />

      <div className="flex gap-3">
        <div
          style={{ height: "20px", width: "180px", backgroundColor: "#E5E7EB" }}
          className="rounded-full animate-pulse"
        />
        <div
          style={{ height: "20px", width: "180px", backgroundColor: "#E5E7EB" }}
          className="rounded-full animate-pulse"
        />
      </div>
    </div>
  );
}
