"use client";

import { useEffect, useState } from "react";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import Link from "next/link";
import api from "@/lib/api";
import type { Ranking } from "@/lib/api";
import { PageLoading, ErrorState, EmptyState } from "@/components/StateDisplays";
import { Building2, Search } from "lucide-react";

export default function OrgListPage() {
  const { isAuthenticated, loading: authLoading } = useAuth();
  const router = useRouter();

  const [orgName, setOrgName] = useState("");
  const [ranking, setRanking] = useState<Ranking | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [authLoading, isAuthenticated, router]);

  const handleSearch = async () => {
    const trimmed = orgName.trim();
    if (!trimmed) return;

    setLoading(true);
    setError(null);
    setRanking(null);

    try {
      const data = await api.evaluateOrg(trimmed);
      setRanking(data);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to evaluate organization"
      );
    } finally {
      setLoading(false);
    }
  };

  if (authLoading || !isAuthenticated) return null;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-5xl">
        <div className="mb-6">
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
            Organization Evaluation
          </h1>
          <p className="mt-1 text-sm text-zinc-500 dark:text-zinc-400">
            Enter a GitHub organization name to evaluate all its public members.
          </p>
        </div>

        {/* Search */}
        <div className="mb-6">
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleSearch();
            }}
            className="relative"
          >
            <Building2 className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-zinc-400" />
            <input
              type="text"
              value={orgName}
              onChange={(e) => setOrgName(e.target.value)}
              placeholder="GitHub organization (e.g., google, facebook)"
              className="w-full rounded-xl border border-zinc-200 bg-white py-2.5 pl-10 pr-24 text-sm text-zinc-900 placeholder:text-zinc-400 focus:border-zinc-400 focus:outline-none focus:ring-2 focus:ring-zinc-200 dark:border-zinc-700 dark:bg-zinc-900 dark:text-white dark:placeholder:text-zinc-500 dark:focus:border-zinc-600 dark:focus:ring-zinc-800"
            />
            <button
              type="submit"
              disabled={loading}
              className="absolute right-2 top-1/2 -translate-y-1/2 rounded-lg bg-zinc-900 px-3 py-1 text-xs font-medium text-white hover:bg-zinc-700 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
            >
              {loading ? "Evaluating…" : "Evaluate"}
            </button>
          </form>
        </div>

        {/* Results */}
        {loading && (
          <div>
            <PageLoading />
            <p className="mt-2 text-center text-sm text-zinc-500 dark:text-zinc-400">
              Fetching and scoring organization members… This may take a moment.
            </p>
          </div>
        )}

        {error && <ErrorState message={error} onRetry={handleSearch} />}

        {!loading && ranking && (
          <div>
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-lg font-semibold text-zinc-900 dark:text-white">
                {ranking.total_developers} members scored
              </h2>
              <span className="text-xs text-zinc-400">
                {ranking.snapshot_date}
              </span>
            </div>
            <div className="overflow-hidden rounded-xl border border-zinc-200 dark:border-zinc-800">
              <div className="overflow-x-auto">
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
            </div>
          </div>
        )}

        {!loading && !ranking && !error && (
          <EmptyState
            title="Search an organization"
            description="Enter a GitHub organization name to evaluate its members."
            icon={<Search className="h-6 w-6 text-zinc-400" />}
          />
        )}
      </div>
    </div>
  );
}
