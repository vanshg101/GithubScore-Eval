"use client";

import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts";

interface PRMergeRateChartProps {
  totalPRs: number;
  mergedPRs: number;
}

const COLORS = ["#22c55e", "#ef4444"];

export function PRMergeRateChart({ totalPRs, mergedPRs }: PRMergeRateChartProps) {
  const unmerged = totalPRs - mergedPRs;
  const mergeRate = totalPRs > 0 ? ((mergedPRs / totalPRs) * 100).toFixed(1) : "0";

  const data = [
    { name: "Merged", value: mergedPRs },
    { name: "Unmerged", value: unmerged },
  ];

  if (totalPRs === 0) {
    return (
      <div className="flex h-[250px] items-center justify-center">
        <p className="text-sm text-zinc-500 dark:text-zinc-400">No PRs found</p>
      </div>
    );
  }

  return (
    <div className="h-[250px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            innerRadius={55}
            outerRadius={80}
            dataKey="value"
            strokeWidth={2}
            stroke="transparent"
          >
            {data.map((_, index) => (
              <Cell key={index} fill={COLORS[index]} />
            ))}
          </Pie>
          <Tooltip
            contentStyle={{
              backgroundColor: "rgba(24,24,27,0.95)",
              border: "1px solid #3f3f46",
              borderRadius: "8px",
              color: "#fff",
              fontSize: "12px",
            }}
          />
          <text
            x="50%"
            y="48%"
            textAnchor="middle"
            dominantBaseline="middle"
            className="fill-zinc-900 text-2xl font-bold dark:fill-white"
          >
            {mergeRate}%
          </text>
          <text
            x="50%"
            y="58%"
            textAnchor="middle"
            dominantBaseline="middle"
            className="fill-zinc-500 text-xs"
          >
            merge rate
          </text>
        </PieChart>
      </ResponsiveContainer>
      <div className="flex justify-center gap-4 text-xs">
        <span className="flex items-center gap-1">
          <span className="h-2 w-2 rounded-full bg-green-500" />
          Merged ({mergedPRs})
        </span>
        <span className="flex items-center gap-1">
          <span className="h-2 w-2 rounded-full bg-red-500" />
          Unmerged ({unmerged})
        </span>
      </div>
    </div>
  );
}
