"use client";

import {
  RadarChart,
  PolarGrid,
  PolarAngleAxis,
  PolarRadiusAxis,
  Radar,
  ResponsiveContainer,
  Tooltip,
} from "recharts";
import type { IndicatorScore } from "@/lib/api";

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
  issue_response_time: "Response Time",
  commit_trend: "Trend",
  language_diversity: "Languages",
};

interface ScoreRadarChartProps {
  indicatorScores: Record<string, IndicatorScore>;
}

export function ScoreRadarChart({ indicatorScores }: ScoreRadarChartProps) {
  const data = Object.entries(indicatorScores).map(([key, value]) => ({
    indicator: INDICATOR_LABELS[key] || key,
    score: Math.round(value.normalized * 100),
    weighted: Math.round(value.weighted * 100) / 100,
  }));

  return (
    <div className="h-[350px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <RadarChart data={data} cx="50%" cy="50%" outerRadius="70%">
          <PolarGrid stroke="#a1a1aa" strokeOpacity={0.3} />
          <PolarAngleAxis
            dataKey="indicator"
            tick={{ fill: "#71717a", fontSize: 11 }}
          />
          <PolarRadiusAxis
            angle={90}
            domain={[0, 100]}
            tick={{ fill: "#a1a1aa", fontSize: 10 }}
          />
          <Radar
            name="Score"
            dataKey="score"
            stroke="#3b82f6"
            fill="#3b82f6"
            fillOpacity={0.25}
            strokeWidth={2}
          />
          <Tooltip
            contentStyle={{
              backgroundColor: "rgba(24,24,27,0.95)",
              border: "1px solid #3f3f46",
              borderRadius: "8px",
              color: "#fff",
              fontSize: "12px",
            }}
            formatter={(value: number | undefined) => [`${value ?? 0}%`, "Normalized"]}
          />
        </RadarChart>
      </ResponsiveContainer>
    </div>
  );
}
