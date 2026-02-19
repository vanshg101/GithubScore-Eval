"use client";

import { useEffect, useState, useCallback, use } from "react";
import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import api from "@/lib/api";
import type { Developer, Score } from "@/lib/api";
import { PageLoading, ErrorState } from "@/components/StateDisplays";
import { SearchBar } from "@/components/SearchBar";
import {
  ScoreRadarChart,
  CommitTrendChart,
  PRMergeRateChart,
  LanguageDistributionChart,
  ContributionHeatmap,
  MLImpactScore,
} from "@/components/charts";
import {
  GitCommit,
  GitPullRequest,
  CircleDot,
  MessageSquare,
  Star,
  GitFork,
  BookOpen,
  Users,
} from "lucide-react";

export default function DeveloperProfilePage({
  params,
}: {
  params: Promise<{ username: string }>;
}) {
  const { username } = use(params);
  const { isAuthenticated, loading: authLoading } = useAuth();
  const router = useRouter();

  const [developer, setDeveloper] = useState<Developer | null>(null);
  const [score, setScore] = useState<Score | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [fetching, setFetching] = useState(false);

  const loadData = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      // Try to get existing developer data
      let dev: Developer;
      try {
        dev = await api.getDeveloper(username);
      } catch {
        // If not found, fetch from GitHub
        setFetching(true);
        dev = await api.fetchDeveloper(username);
        setFetching(false);
      }
      setDeveloper(dev);

      // Try to get existing score, compute if not found
      let scoreData: Score;
      try {
        scoreData = await api.getScore(username);
      } catch {
        scoreData = await api.computeScore(username);
      }
      setScore(scoreData);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load developer data");
    } finally {
      setLoading(false);
      setFetching(false);
    }
  }, [username]);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
      return;
    }
    if (!authLoading && isAuthenticated) {
      loadData();
    }
  }, [authLoading, isAuthenticated, router, loadData]);

  if (authLoading || loading) {
    return (
      <div className="flex-1 p-6 lg:p-8">
        <PageLoading />
        {fetching && (
          <p className="mt-4 text-center text-sm text-zinc-500 dark:text-zinc-400">
            Fetching GitHub data for {username}…
          </p>
        )}
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex-1 p-6 lg:p-8">
        <div className="mx-auto max-w-5xl">
          <SearchBar className="mb-6" />
          <ErrorState message={error} onRetry={loadData} />
        </div>
      </div>
    );
  }

  if (!developer || !score) return null;

  const m = developer.metrics;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-6xl">
        <SearchBar className="mb-6" />

        {/* Profile Header */}
        <div className="mb-8 flex flex-col items-start gap-6 sm:flex-row sm:items-center">
          {developer.profile.avatar_url && (
            <img
              src={developer.profile.avatar_url}
              alt={developer.username}
              className="h-20 w-20 rounded-full border-2 border-zinc-200 dark:border-zinc-700"
            />
          )}
          <div className="flex-1">
            <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
              {developer.profile.name || developer.username}
            </h1>
            <p className="text-sm text-zinc-500 dark:text-zinc-400">
              @{developer.username}
            </p>
            {developer.profile.bio && (
              <p className="mt-1 text-sm text-zinc-600 dark:text-zinc-300">
                {developer.profile.bio}
              </p>
            )}
            <div className="mt-2 flex flex-wrap gap-4 text-xs text-zinc-500 dark:text-zinc-400">
              <span className="flex items-center gap-1">
                <BookOpen className="h-3.5 w-3.5" />
                {developer.profile.public_repos} repos
              </span>
              <span className="flex items-center gap-1">
                <Users className="h-3.5 w-3.5" />
                {developer.profile.followers} followers
              </span>
            </div>
          </div>
          <div className="text-right">
            <div className="text-3xl font-bold text-zinc-900 dark:text-white">
              {score.weighted_score.toFixed(1)}
            </div>
            <div className="text-xs text-zinc-500 dark:text-zinc-400">
              Weighted Score
            </div>
            {score.percentile > 0 && (
              <div className="mt-1 text-xs font-medium text-blue-600 dark:text-blue-400">
                Top {(100 - score.percentile).toFixed(0)}%
              </div>
            )}
          </div>
        </div>

        {/* Stats Grid */}
        <div className="mb-8 grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-4">
          <StatCard icon={<GitCommit className="h-4 w-4" />} label="Commits" value={m.total_commits} />
          <StatCard icon={<GitPullRequest className="h-4 w-4" />} label="PRs" value={`${m.merged_prs}/${m.total_prs}`} />
          <StatCard icon={<CircleDot className="h-4 w-4" />} label="Issues" value={`${m.total_issues_closed}/${m.total_issues_opened}`} />
          <StatCard icon={<MessageSquare className="h-4 w-4" />} label="Reviews" value={m.review_comments} />
          <StatCard icon={<Star className="h-4 w-4" />} label="Stars" value={m.total_stars} />
          <StatCard icon={<GitFork className="h-4 w-4" />} label="Forks" value={m.total_forks} />
          <StatCard icon={<BookOpen className="h-4 w-4" />} label="Repos" value={m.repos_contributed} />
          <StatCard label="Avg PR Lines" value={m.avg_pr_lines_changed.toFixed(0)} />
        </div>

        {/* Charts Grid */}
        <div className="grid gap-6 lg:grid-cols-2">
          {/* ML Impact Score */}
          <MLImpactScore score={score.ml_impact_score} />

          {/* Score Radar */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="mb-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              Score Breakdown
            </h3>
            <ScoreRadarChart indicatorScores={score.indicator_scores} />
          </div>

          {/* Commit Trend */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="mb-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              Commit Trend (12 months)
            </h3>
            <CommitTrendChart
              totalCommits={m.total_commits}
              activeWeeks={m.active_weeks}
              trend={m.commit_trend}
            />
          </div>

          {/* PR Merge Rate */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="mb-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              PR Merge Rate
            </h3>
            <PRMergeRateChart totalPRs={m.total_prs} mergedPRs={m.merged_prs} />
          </div>

          {/* Language Distribution */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="mb-2 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              Language Distribution
            </h3>
            <LanguageDistributionChart languages={m.languages} />
          </div>

          {/* Contribution Heatmap */}
          <div className="rounded-xl border border-zinc-200 bg-white p-5 dark:border-zinc-800 dark:bg-zinc-900">
            <h3 className="mb-3 text-sm font-medium text-zinc-500 dark:text-zinc-400">
              Contribution Heatmap
            </h3>
            <ContributionHeatmap
              activeWeeks={m.active_weeks}
              totalCommits={m.total_commits}
            />
          </div>
        </div>
      </div>
    </div>
  );
}

function StatCard({
  icon,
  label,
  value,
}: {
  icon?: React.ReactNode;
  label: string;
  value: string | number;
}) {
  return (
    <div className="rounded-xl border border-zinc-200 bg-white p-4 dark:border-zinc-800 dark:bg-zinc-900">
      <div className="flex items-center gap-1.5 text-zinc-400">
        {icon}
        <span className="text-xs">{label}</span>
      </div>
      <div className="mt-1 text-lg font-bold text-zinc-900 dark:text-white">
        {value}
      </div>
    </div>
  );
}
