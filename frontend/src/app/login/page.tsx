"use client";

import { useAuth } from "@/context/AuthContext";
import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { Github, BarChart3, Users, Zap } from "lucide-react";

export default function LoginPage() {
  const { login, isAuthenticated, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && isAuthenticated) {
      router.push("/dashboard");
    }
  }, [isAuthenticated, loading, router]);

  if (loading) {
    return (
      <div className="flex min-h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-zinc-300 border-t-zinc-900 dark:border-zinc-700 dark:border-t-white" />
      </div>
    );
  }

  return (
    <div className="flex min-h-screen flex-col items-center justify-center px-4">
      <div className="w-full max-w-md space-y-8">
        {/* Header */}
        <div className="text-center">
          <div className="mx-auto mb-4 flex h-16 w-16 items-center justify-center rounded-2xl bg-zinc-900 dark:bg-white">
            <Github className="h-8 w-8 text-white dark:text-zinc-900" />
          </div>
          <h1 className="text-3xl font-bold text-zinc-900 dark:text-white">
            GitScore
          </h1>
          <p className="mt-2 text-zinc-500 dark:text-zinc-400">
            GitHub Contribution Scoring &amp; Developer Evaluation
          </p>
        </div>

        {/* Features */}
        <div className="space-y-3">
          <FeatureItem
            icon={<BarChart3 className="h-5 w-5" />}
            title="Contribution Scoring"
            description="12 weighted indicators for comprehensive evaluation"
          />
          <FeatureItem
            icon={<Zap className="h-5 w-5" />}
            title="ML Impact Prediction"
            description="Ridge regression model predicts developer impact"
          />
          <FeatureItem
            icon={<Users className="h-5 w-5" />}
            title="Rankings & Comparison"
            description="Compare developers and evaluate organizations"
          />
        </div>

        {/* Login button */}
        <button
          onClick={login}
          className="flex w-full items-center justify-center gap-3 rounded-xl bg-zinc-900 px-6 py-3.5 text-base font-semibold text-white shadow-lg transition-all hover:bg-zinc-700 hover:shadow-xl dark:bg-white dark:text-zinc-900 dark:hover:bg-zinc-200"
        >
          <Github className="h-5 w-5" />
          Sign in with GitHub
        </button>

        <p className="text-center text-xs text-zinc-400 dark:text-zinc-600">
          We only request read access to your public GitHub profile.
        </p>
      </div>
    </div>
  );
}

function FeatureItem({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="flex items-start gap-3 rounded-lg border border-zinc-200 p-4 dark:border-zinc-800">
      <div className="mt-0.5 text-zinc-600 dark:text-zinc-400">{icon}</div>
      <div>
        <h3 className="text-sm font-semibold text-zinc-900 dark:text-white">
          {title}
        </h3>
        <p className="text-xs text-zinc-500 dark:text-zinc-400">
          {description}
        </p>
      </div>
    </div>
  );
}
