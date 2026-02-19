"use client";

interface ContributionHeatmapProps {
  activeWeeks: number;
  totalCommits: number;
}

export function ContributionHeatmap({
  activeWeeks,
  totalCommits,
}: ContributionHeatmapProps) {
  // Generate 52 weeks of simulated activity
  const weeks = 52;
  const avgPerWeek = totalCommits / Math.max(activeWeeks, 1);

  const cells = Array.from({ length: weeks }, (_, i) => {
    const isActive = i < activeWeeks;
    if (!isActive) return 0;
    // Random-ish distribution around avg
    const seed = Math.sin(i * 127.1 + 311.7) * 43758.5453;
    const rand = seed - Math.floor(seed);
    return Math.max(0, Math.round(avgPerWeek * (0.3 + rand * 1.4)));
  });

  const maxVal = Math.max(...cells, 1);

  const getColor = (val: number) => {
    if (val === 0) return "bg-zinc-100 dark:bg-zinc-800";
    const intensity = val / maxVal;
    if (intensity > 0.75) return "bg-green-600 dark:bg-green-500";
    if (intensity > 0.5) return "bg-green-500 dark:bg-green-600";
    if (intensity > 0.25) return "bg-green-400 dark:bg-green-700";
    return "bg-green-300 dark:bg-green-800";
  };

  return (
    <div>
      <div className="flex flex-wrap gap-1">
        {cells.map((val, i) => (
          <div
            key={i}
            className={`h-3 w-3 rounded-sm ${getColor(val)}`}
            title={`Week ${i + 1}: ${val} commits`}
          />
        ))}
      </div>
      <div className="mt-3 flex items-center gap-2 text-xs text-zinc-500 dark:text-zinc-400">
        <span>Less</span>
        <div className="flex gap-0.5">
          <div className="h-3 w-3 rounded-sm bg-zinc-100 dark:bg-zinc-800" />
          <div className="h-3 w-3 rounded-sm bg-green-300 dark:bg-green-800" />
          <div className="h-3 w-3 rounded-sm bg-green-400 dark:bg-green-700" />
          <div className="h-3 w-3 rounded-sm bg-green-500 dark:bg-green-600" />
          <div className="h-3 w-3 rounded-sm bg-green-600 dark:bg-green-500" />
        </div>
        <span>More</span>
      </div>
    </div>
  );
}
