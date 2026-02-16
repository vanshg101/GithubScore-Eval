"use client";

import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

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

        {/* Quick stats cards */}
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <DashCard
            title="Search Developer"
            description="Look up any GitHub user by username to view their contribution score."
            href="/search"
            action="Search →"
          />
          <DashCard
            title="Leaderboard"
            description="View the ranked leaderboard of all scored developers."
            href="/leaderboard"
            action="View Rankings →"
          />
          <DashCard
            title="Compare"
            description="Compare 2–10 developers side by side with scoring breakdown."
            href="/compare"
            action="Compare →"
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
  action,
}: {
  title: string;
  description: string;
  href: string;
  action: string;
}) {
  return (
    <a
      href={href}
      className="group rounded-xl border border-zinc-200 bg-white p-6 transition-all hover:border-zinc-300 hover:shadow-md dark:border-zinc-800 dark:bg-zinc-900 dark:hover:border-zinc-700"
    >
      <h2 className="text-lg font-semibold text-zinc-900 dark:text-white">
        {title}
      </h2>
      <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
        {description}
      </p>
      <span className="mt-4 inline-block text-sm font-medium text-zinc-900 group-hover:underline dark:text-white">
        {action}
      </span>
    </a>
  );
}
