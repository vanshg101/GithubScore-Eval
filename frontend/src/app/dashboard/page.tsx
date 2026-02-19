"use client";

import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { SearchBar } from "@/components/SearchBar";
import { Search, Trophy, GitCompareArrows, Building2 } from "lucide-react";
import Link from "next/link";

export default function DashboardPage() {
  const { user, isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-zinc-300 border-t-zinc-900 dark:border-zinc-700 dark:border-t-white" />
      </div>
    );
  }

  if (!isAuthenticated) return null;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-5xl">
        {/* Welcome header */}
        <div className="mb-8">
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
            Welcome back, {user?.display_name || user?.username}!
          </h1>
          <p className="mt-1 text-zinc-500 dark:text-zinc-400">
            Search for a GitHub username to analyze their contributions.
          </p>
        </div>

        <SearchBar className="mb-8" />

        {/* Quick stats cards */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
          <DashCard
            title="Search Developer"
            description="Look up any GitHub user to view their score."
            href="/search"
            icon={<Search className="h-5 w-5" />}
          />
          <DashCard
            title="Leaderboard"
            description="View the ranked leaderboard of scored developers."
            href="/leaderboard"
            icon={<Trophy className="h-5 w-5" />}
          />
          <DashCard
            title="Compare"
            description="Compare 2–10 developers side by side."
            href="/compare"
            icon={<GitCompareArrows className="h-5 w-5" />}
          />
          <DashCard
            title="Organizations"
            description="Evaluate all members of a GitHub org."
            href="/org"
            icon={<Building2 className="h-5 w-5" />}
          />
        </div>
      </div>
    </div>
  );
}

function DashCard({
  title,
  description,
  href,
  icon,
}: {
  title: string;
  description: string;
  href: string;
  icon: React.ReactNode;
}) {
  return (
    <Link
      href={href}
      className="group rounded-xl border border-zinc-200 bg-white p-5 transition-all hover:border-zinc-300 hover:shadow-md dark:border-zinc-800 dark:bg-zinc-900 dark:hover:border-zinc-700"
    >
      <div className="mb-3 flex h-9 w-9 items-center justify-center rounded-lg bg-zinc-100 text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
        {icon}
      </div>
      <h2 className="text-sm font-semibold text-zinc-900 dark:text-white">
        {title}
      </h2>
      <p className="mt-1 text-xs text-zinc-500 dark:text-zinc-400">
        {description}
      </p>
    </Link>
  );
}
