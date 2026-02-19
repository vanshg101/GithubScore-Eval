"use client";

import { Zap } from "lucide-react";

interface MLImpactScoreProps {
  score: number;
}

export function MLImpactScore({ score }: MLImpactScoreProps) {
  const clampedScore = Math.min(100, Math.max(0, score));
  const getColorClass = () => {
    if (clampedScore >= 80) return "text-green-600 dark:text-green-400";
    if (clampedScore >= 60) return "text-blue-600 dark:text-blue-400";
    if (clampedScore >= 40) return "text-amber-600 dark:text-amber-400";
    return "text-red-600 dark:text-red-400";
  };

  const getBarColor = () => {
    if (clampedScore >= 80) return "bg-green-500";
    if (clampedScore >= 60) return "bg-blue-500";
    if (clampedScore >= 40) return "bg-amber-500";
    return "bg-red-500";
  };

  const getLabel = () => {
    if (clampedScore >= 80) return "Excellent";
    if (clampedScore >= 60) return "Good";
    if (clampedScore >= 40) return "Average";
    return "Below Average";
  };

  return (
    <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
      <div className="flex items-center gap-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
        <Zap className="h-4 w-4" />
        ML Impact Score
      </div>
      <div className="mt-3 flex items-baseline gap-2">
        <span className={`text-4xl font-bold ${getColorClass()}`}>
          {clampedScore.toFixed(1)}
        </span>
        <span className="text-sm text-zinc-400">/100</span>
      </div>
      <div className="mt-3">
        <div className="h-2 w-full overflow-hidden rounded-full bg-zinc-100 dark:bg-zinc-800">
          <div
            className={`h-full rounded-full ${getBarColor()} transition-all duration-500`}
            style={{ width: `${clampedScore}%` }}
          />
        </div>
      </div>
      <p className={`mt-2 text-xs font-medium ${getColorClass()}`}>{getLabel()}</p>
    </div>
  );
}
