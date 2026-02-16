import { Github } from "lucide-react";

export function Footer() {
  return (
    <footer className="border-t border-zinc-200 bg-white dark:border-zinc-800 dark:bg-zinc-950">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-4 sm:px-6 lg:px-8">
        <p className="text-xs text-zinc-500 dark:text-zinc-500">
          &copy; {new Date().getFullYear()} GitScore. Built with Next.js, Go &amp; FastAPI.
        </p>
        <a
          href="https://github.com/vanshg101/GithubScore-Eval"
          target="_blank"
          rel="noopener noreferrer"
          className="text-zinc-400 hover:text-zinc-600 dark:hover:text-zinc-300"
        >
          <Github className="h-5 w-5" />
        </a>
      </div>
    </footer>
  );
}
