"use client";

import { useEffect, useState, useCallback } from "react";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import Link from "next/link";
import api from "@/lib/api";
import type { Ranking, RankEntry } from "@/lib/api";
import { PageLoading, ErrorState, EmptyState } from "@/components/StateDisplays";
import { Trophy, ChevronUp, ChevronDown, ChevronLeft, ChevronRight } from "lucide-react";

type SortKey = "rank" | "username" | "score" | "ml_score";
type SortDir = "asc" | "desc";

export default function LeaderboardPage() {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const router = useRouter();

  const [ranking, setRanking] = useState<Ranking | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [sortKey, setSortKey] = useState<SortKey>("rank");
  const [sortDir, setSortDir] = useState<SortDir>("asc");

  const pageSize = 20;

  const loadData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await api.getRankings(page, pageSize);
      setRanking(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load rankings");
    } finally {
      setLoading(false);
    }
  }, [page]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
      return;
    }
    if (!authLoading && isAuthenticated) {
      loadData();
    }
  }, [authLoading, isAuthenticated, router, loadData]);

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir(key === "rank" ? "asc" : "desc");
    }
  };

  const sortedRankings = ranking
    ? [...ranking.rankings].sort((a, b) => {
        const aVal = a[sortKey];
        const bVal = b[sortKey];
        if (typeof aVal === "string" && typeof bVal === "string") {
          return sortDir === "asc"
            ? aVal.localeCompare(bVal)
            : bVal.localeCompare(aVal);
        }
        return sortDir === "asc"
          ? (aVal as number) - (bVal as number)
          : (bVal as number) - (aVal as number);
      })
    : [];

  const totalPages = ranking
    ? Math.ceil(ranking.total_developers / pageSize)
    : 1;

  if (authLoading) return null;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6 flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
              Leaderboard
            </h1>
            <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
              {ranking
                ? `${ranking.total_developers} developers ranked`
                : "Loading rankings…"}
            </p>
          </div>
          {ranking && (
            <span className="text-xs text-zinc-400">
              Snapshot: {ranking.snapshot_date}
            </span>
          )}
        </div>

        {loading ? (
          <PageLoading />
        ) : error ? (
          <ErrorState message={error} onRetry={loadData} />
        ) : !ranking || ranking.rankings.length === 0 ? (
          <EmptyState
            title="No rankings yet"
            description="Search and score developers to populate the leaderboard."
            icon={<Trophy className="h-6 w-6 text-zinc-400" />}
          />
        ) : (
          <>
            {/* Table */}
            <div className="overflow-hidden rounded-xl border border-zinc-200 dark:border-zinc-800">
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
                      <SortableHeader
                        label="Rank"
                        sortKey="rank"
                        currentKey={sortKey}
                        direction={sortDir}
                        onClick={handleSort}
                        className="w-20"
                      />
                      <SortableHeader
                        label="Developer"
                        sortKey="username"
                        currentKey={sortKey}
                        direction={sortDir}
                        onClick={handleSort}
                      />
                      <SortableHeader
                        label="Score"
                        sortKey="score"
                        currentKey={sortKey}
                        direction={sortDir}
                        onClick={handleSort}
                        className="text-right"
                      />
                      <SortableHeader
                        label="ML Score"
                        sortKey="ml_score"
                        currentKey={sortKey}
                        direction={sortDir}
                        onClick={handleSort}
                        className="text-right"
                      />
                    </tr>
                  </thead>
                  <tbody>
                    {sortedRankings.map((entry) => (
                      <RankRow key={entry.username} entry={entry} />
                    ))}
                  </tbody>
                </table>
              </div>
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
              <div className="mt-4 flex items-center justify-between">
                <span className="text-xs text-zinc-500 dark:text-zinc-400">
                  Page {page} of {totalPages}
                </span>
                <div className="flex gap-2">
                  <button
                    onClick={() => setPage((p) => Math.max(1, p - 1))}
                    disabled={page === 1}
                    className="inline-flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium text-zinc-700 hover:bg-zinc-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-zinc-700 dark:text-zinc-300 dark:hover:bg-zinc-800"
                  >
                    <ChevronLeft className="h-3.5 w-3.5" />
                    Prev
                  </button>
                  <button
                    onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                    disabled={page === totalPages}
                    className="inline-flex items-center gap-1 rounded-lg border border-zinc-200 px-3 py-1.5 text-xs font-medium text-zinc-700 hover:bg-zinc-50 disabled:cursor-not-allowed disabled:opacity-50 dark:border-zinc-700 dark:text-zinc-300 dark:hover:bg-zinc-800"
                  >
                    Next
                    <ChevronRight className="h-3.5 w-3.5" />
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}

function SortableHeader({
  label,
  sortKey,
  currentKey,
  direction,
  onClick,
  className = "",
}: {
  label: string;
  sortKey: SortKey;
  currentKey: SortKey;
  direction: SortDir;
  onClick: (key: SortKey) => void;
  className?: string;
}) {
  const isActive = currentKey === sortKey;
  return (
    <th
      className={`cursor-pointer px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-zinc-500 hover:text-zinc-700 dark:text-zinc-400 dark:hover:text-zinc-200 ${className}`}
      onClick={() => onClick(sortKey)}
    >
      <span className="inline-flex items-center gap-1">
        {label}
        {isActive &&
          (direction === "asc" ? (
            <ChevronUp className="h-3 w-3" />
          ) : (
            <ChevronDown className="h-3 w-3" />
          ))}
      </span>
    </th>
  );
}

function RankRow({ entry }: { entry: RankEntry }) {
  const getRankBadge = (rank: number) => {
    if (rank === 1) return "🥇";
    if (rank === 2) return "🥈";
    if (rank === 3) return "🥉";
    return `#${rank}`;
  };

  return (
    <tr className="border-b border-zinc-100 transition-colors hover:bg-zinc-50 dark:border-zinc-800/50 dark:hover:bg-zinc-900/50">
      <td className="px-4 py-3 text-sm font-medium text-zinc-900 dark:text-white">
        {getRankBadge(entry.rank)}
      </td>
      <td className="px-4 py-3">
        <Link
          href={`/developer/${entry.username}`}
          className="text-sm font-medium text-blue-600 hover:underline dark:text-blue-400"
        >
          {entry.username}
        </Link>
      </td>
      <td className="px-4 py-3 text-right text-sm font-semibold text-zinc-900 dark:text-white">
        {entry.score.toFixed(1)}
      </td>
      <td className="px-4 py-3 text-right text-sm text-zinc-600 dark:text-zinc-400">
        {entry.ml_score.toFixed(1)}
      </td>
    </tr>
  );
}
