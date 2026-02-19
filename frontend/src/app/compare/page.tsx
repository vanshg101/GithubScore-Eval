"use client";

import { useEffect, useState, useCallback } from "react";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import Link from "next/link";
import api from "@/lib/api";
import type { Ranking, Score } from "@/lib/api";
import { PageLoading, ErrorState, EmptyState } from "@/components/StateDisplays";
import { ComparisonOverlayChart } from "@/components/charts";
import { GitCompareArrows, X, Plus } from "lucide-react";

export default function ComparePage() {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const router = useRouter();

  const [usernames, setUsernames] = useState<string[]>([""]);
  const [ranking, setRanking] = useState<Ranking | null>(null);
  const [scores, setScores] = useState<Score[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  const addUsername = () => {
    if (usernames.length < 10) {
      setUsernames([...usernames, ""]);
    }
  };

  const removeUsername = (index: number) => {
    setUsernames(usernames.filter((_, i) => i !== index));
  };

  const updateUsername = (index: number, value: string) => {
    const updated = [...usernames];
    updated[index] = value;
    setUsernames(updated);
  };

  const handleCompare = useCallback(async () => {
    const validUsernames = usernames
      .map((u) => u.trim())
      .filter((u) => u.length > 0);

    if (validUsernames.length < 2) {
      setError("Enter at least 2 usernames to compare");
      return;
    }

    setLoading(true);
    setError(null);
    setScores([]);

    try {
      const data = await api.compareDevelopers(validUsernames);
      setRanking(data);

      // Fetch individual scores for radar overlay
      const scorePromises = validUsernames.map(async (u) => {
        try {
          return await api.getScore(u);
        } catch {
          return null;
        }
      });
      const scoreResults = await Promise.all(scorePromises);
      setScores(scoreResults.filter((s): s is Score => s !== null));
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to compare developers"
      );
    } finally {
      setLoading(false);
    }
  }, [usernames]);

  if (authLoading || !isAuthenticated) return null;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
            Compare Developers
          </h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Enter 2–10 GitHub usernames for side-by-side comparison.
          </p>
        </div>

        {/* Input Form */}
        <div className="mb-6 rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
          <div className="space-y-2">
            {usernames.map((username, index) => (
              <div key={index} className="flex items-center gap-2">
                <input
                  type="text"
                  value={username}
                  onChange={(e) => updateUsername(index, e.target.value)}
                  placeholder={`GitHub username ${index + 1}`}
                  className="flex-1 rounded-lg border border-zinc-200 bg-white px-3 py-2 text-sm text-zinc-900 placeholder:text-zinc-400 focus:border-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-200 dark:border-zinc-700 dark:bg-zinc-800 dark:text-white dark:placeholder:text-zinc-500 dark:focus:border-zinc-600 dark:focus:ring-zinc-700"
                  onKeyDown={(e) => {
                    if (e.key === "Enter") handleCompare();
                  }}
                />
                {usernames.length > 1 && (
                  <button
                    onClick={() => removeUsername(index)}
                    className="rounded-lg p-2 text-zinc-400 hover:bg-zinc-100 hover:text-zinc-600 dark:hover:bg-zinc-800 dark:hover:text-zinc-300"
                  >
                    <X className="h-4 w-4" />
                  </button>
                )}
              </div>
            ))}
          </div>
          <div className="mt-3 flex items-center gap-3">
            {usernames.length < 10 && (
              <button
                onClick={addUsername}
                className="inline-flex items-center gap-1.5 rounded-lg border border-dashed border-zinc-300 px-3 py-1.5 text-xs font-medium text-zinc-500 hover:border-zinc-400 hover:text-zinc-700 dark:border-zinc-700 dark:text-zinc-400 dark:hover:border-zinc-600 dark:hover:text-zinc-300"
              >
                <Plus className="h-3.5 w-3.5" />
                Add User
              </button>
            )}
            <button
              onClick={handleCompare}
              disabled={loading}
              className="inline-flex items-center gap-2 rounded-lg bg-zinc-900 px-4 py-2 text-sm font-medium text-white hover:bg-zinc-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              <GitCompareArrows className="h-4 w-4" />
              {loading ? "Comparing…" : "Compare"}
            </button>
          </div>
          {error && (
            <p className="mt-3 text-sm text-red-600 dark:text-red-400">
              {error}
            </p>
          )}
        </div>

        {/* Results */}
        {loading && <PageLoading />}

        {!loading && ranking && (
          <div className="space-y-6">
            {/* Ranking Table */}
            <div className="overflow-hidden rounded-xl border border-zinc-200 dark:border-zinc-800">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-zinc-200 bg-zinc-50 dark:border-zinc-800 dark:bg-zinc-900/50">
                    <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      Rank
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      Developer
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      Score
                    </th>
                    <th className="px-4 py-3 text-right text-xs font-medium uppercase tracking-wider text-zinc-500 dark:text-zinc-400">
                      ML Score
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {ranking.rankings.map((entry) => (
                    <tr
                      key={entry.username}
                      className="border-b border-zinc-100 transition-colors hover:bg-zinc-50 dark:border-zinc-800/50 dark:hover:bg-zinc-900/50"
                    >
                      <td className="px-4 py-3 text-sm font-medium text-zinc-900 dark:text-white">
                        #{entry.rank}
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
                  ))}
                </tbody>
              </table>
            </div>

            {/* Overlay Chart */}
            {scores.length >= 2 && (
              <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
                <h3 className="mb-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
                  Score Breakdown Comparison
                </h3>
                <ComparisonOverlayChart scores={scores} />
              </div>
            )}
          </div>
        )}

        {!loading && !ranking && !error && (
          <EmptyState
            title="Ready to compare"
            description="Enter usernames above and click Compare to see side-by-side analysis."
            icon={<GitCompareArrows className="h-6 w-6 text-zinc-400" />}
          />
        )}
      </div>
    </div>
  );
}
