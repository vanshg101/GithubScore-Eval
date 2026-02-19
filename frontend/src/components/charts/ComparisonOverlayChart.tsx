"use client";

import {
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
  ResponsiveContainer,
  Tooltip,
  Legend,
} from "recharts";
import type { Score } from "@/lib/api";

const INDICATOR_LABELS: Record<string, string> = {
  total_commits: "Commits",
  pr_merge_rate: "PR Merge",
  issue_resolution: "Issues",
  code_review: "Reviews",
  contribution_consistency: "Consistency",
  repo_diversity: "Repo Diversity",
  stars_earned: "Stars",
  fork_impact: "Forks",
  avg_pr_size: "PR Size",
  issue_response_time: "Response",
  commit_trend: "Trend",
  language_diversity: "Languages",
};

const COLORS = [
  "#3b82f6", "#ef4444", "#22c55e", "#f59e0b", "#8b5cf6",
  "#ec4899", "#06b6d4", "#f97316", "#6366f1", "#14b8a6",
];

interface ComparisonOverlayChartProps {
  scores: Score[];
}

export function ComparisonOverlayChart({ scores }: ComparisonOverlayChartProps) {
  if (scores.length === 0) return null;

  // Build unified data: one entry per indicator with each user's normalized value
  const indicatorKeys = Object.keys(scores[0].indicator_scores);
  const data = indicatorKeys.map((key) => {
    const entry: Record<string, string | number> = {
      indicator: INDICATOR_LABELS[key] || key,
    };
    scores.forEach((s) => {
      entry[s.username] = Math.round(
        (s.indicator_scores[key]?.normalized || 0) * 100
      );
    });
    return entry;
  });

  return (
    <div className="h-[400px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <RadarChart data={data} cx="50%" cy="50%" outerRadius="65%">
          <PolarGrid stroke="#a1a1aa" strokeOpacity={0.3} />
          <PolarAngleAxis
            dataKey="indicator"
            tick={{ fill: "#71717a", fontSize: 10 }}
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 100]}
            tick={{ fill: "#a1a1aa", fontSize: 9 }}
          />
          {scores.map((s, i) => (
            <Radar
              key={s.username}
              name={s.username}
              dataKey={s.username}
              stroke={COLORS[i % COLORS.length]}
              fill={COLORS[i % COLORS.length]}
              fillOpacity={0.1}
              strokeWidth={2}
            />
          ))}
          <Legend
            wrapperStyle={{ fontSize: "12px" }}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: "rgba(24,24,27,0.95)",
              border: "1px solid #3f3f46",
              borderRadius: "8px",
              color: "#fff",
              fontSize: "12px",
            }}
          />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  );
}
