"use client";

import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

interface CommitTrendChartProps {
  totalCommits: number;
  activeWeeks: number;
  trend: string;
}

export function CommitTrendChart({
  totalCommits,
  activeWeeks,
  trend,
}: CommitTrendChartProps) {
  // Simulate monthly commit data from total commits and trend
  const months = [
    "Jan", "Feb", "Mar", "Apr", "May", "Jun",
    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
  ];

  const avgPerMonth = totalCommits / 12;
  const data = months.map((month, i) => {
    let factor: number;
    if (trend === "increasing") {
      factor = 0.5 + (i / 11) * 1.0;
    } else if (trend === "decreasing") {
      factor = 1.5 - (i / 11) * 1.0;
    } else {
      factor = 0.85 + Math.sin(i * 0.8) * 0.3;
    }
    return {
      month,
      commits: Math.max(0, Math.round(avgPerMonth * factor)),
    };
  });

  return (
    <div className="h-[280px] w-full">
      <div className="mb-2 flex items-center gap-2">
        <span className="text-xs text-zinc-500 dark:text-zinc-400">
          Trend:{" "}
          <span
            className={
              trend === "increasing"
                ? "font-medium text-green-600 dark:text-green-400"
                : trend === "decreasing"
                ? "font-medium text-red-600 dark:text-red-400"
                : "font-medium text-amber-600 dark:text-amber-400"
            }
          >
            {trend}
          </span>
        </span>
        <span className="text-xs text-zinc-400">·</span>
        <span className="text-xs text-zinc-500 dark:text-zinc-400">
          {activeWeeks} active weeks
        </span>
      </div>
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data}>
          <CartesianGrid strokeDasharray="3 3" stroke="#3f3f46" strokeOpacity={0.3} />
          <XAxis dataKey="month" tick={{ fill: "#71717a", fontSize: 11 }} />
          <YAxis tick={{ fill: "#71717a", fontSize: 11 }} />
          <Tooltip
            contentStyle={{
              backgroundColor: "rgba(24,24,27,0.95)",
              border: "1px solid #3f3f46",
              borderRadius: "8px",
              color: "#fff",
              fontSize: "12px",
            }}
          />
          <Line
            type="monotone"
            dataKey="commits"
            stroke="#8b5cf6"
            strokeWidth={2}
            dot={{ fill: "#8b5cf6", r: 3 }}
            activeDot={{ r: 5 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}
