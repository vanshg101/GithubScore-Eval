"use client";

import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { SearchBar } from "@/components/SearchBar";
import { Search } from "lucide-react";

export default function SearchPage() {
  const { isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, loading, router]);

  if (loading || !isAuthenticated) return null;

  return (
    <div className="flex-1 p-6 lg:p-8">
      <div className="mx-auto max-w-2xl">
        <div className="mb-8 text-center">
          <div className="mx-auto mb-4 flex h-14 w-14 items-center justify-center rounded-2xl bg-zinc-100 dark:bg-zinc-800">
            <Search className="h-6 w-6 text-zinc-600 dark:text-zinc-400" />
          </div>
          <h1 className="text-2xl font-bold text-zinc-900 dark:text-white">
            Search Developer
          </h1>
          <p className="mt-2 text-sm text-zinc-500 dark:text-zinc-400">
            Enter a GitHub username to view their contribution score and analytics.
          </p>
        </div>

        <SearchBar placeholder="Enter GitHub username (e.g., torvalds)" />

        <div className="mt-8 rounded-xl border border-zinc-200 bg-zinc-50 p-6 dark:border-zinc-800 dark:bg-zinc-900/50">
          <h3 className="text-sm font-semibold text-zinc-700 dark:text-zinc-300">
            How it works
          </h3>
          <ul className="mt-3 space-y-2 text-xs text-zinc-500 dark:text-zinc-400">
            <li className="flex items-start gap-2">
              <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-zinc-200 text-[10px] font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                1
              </span>
              Enter a GitHub username and hit Search
            </li>
            <li className="flex items-start gap-2">
              <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-zinc-200 text-[10px] font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                2
              </span>
              We fetch their public GitHub data (commits, PRs, issues, repos)
            </li>
            <li className="flex items-start gap-2">
              <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-zinc-200 text-[10px] font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                3
              </span>
              12 weighted indicators compute a contribution score
            </li>
            <li className="flex items-start gap-2">
              <span className="mt-0.5 flex h-4 w-4 shrink-0 items-center justify-center rounded-full bg-zinc-200 text-[10px] font-medium text-zinc-600 dark:bg-zinc-800 dark:text-zinc-400">
                4
              </span>
              ML model predicts an overall impact score
            </li>
          </ul>
        </div>
      </div>
    </div>
  );
}
