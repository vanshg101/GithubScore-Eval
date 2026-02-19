"use client";

import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, type PieLabelRenderProps } from "recharts";

interface LanguageDistributionChartProps {
  languages: string[];
}

const COLORS = [
  "#3b82f6", "#8b5cf6", "#ec4899", "#f59e0b", "#22c55e",
  "#06b6d4", "#f97316", "#6366f1", "#14b8a6", "#e11d48",
  "#84cc16", "#a855f7",
];

export function LanguageDistributionChart({ languages }: LanguageDistributionChartProps) {
  if (!languages || languages.length === 0) {
    return (
      <div className="flex h-[250px] items-center justify-center">
        <p className="text-sm text-zinc-500 dark:text-zinc-400">No language data</p>
      </div>
    );
  }

  // Count occurrences or treat each as equal weight
  const langCounts = languages.reduce<Record<string, number>>((acc, lang) => {
    acc[lang] = (acc[lang] || 0) + 1;
    return acc;
  }, {});

  const data = Object.entries(langCounts)
    .map(([name, value]) => ({ name, value }))
    .sort((a, b) => b.value - a.value)
    .slice(0, 10);

  return (
    <div className="h-[250px] w-full">
      <ResponsiveContainer width="100%" height="100%">
        <PieChart>
          <Pie
            data={data}
            cx="50%"
            cy="50%"
            outerRadius={80}
            dataKey="value"
            label={({ name, percent }: PieLabelRenderProps) =>
              `${name ?? ""} ${(((percent as number) ?? 0) * 100).toFixed(0)}%`
            }
            labelLine={{ stroke: "#71717a" }}
            strokeWidth={2}
            stroke="transparent"
          >
            {data.map((_, index) => (
              <Cell key={index} fill={COLORS[index % COLORS.length]} />
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
        </PieChart>
      </ResponsiveContainer>
    </div>
  );
}
