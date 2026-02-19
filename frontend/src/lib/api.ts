const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

// ---- Types ----

export interface DeveloperProfile {
  name: string;
  bio: string;
  public_repos: number;
  followers: number;
  avatar_url: string;
}

export interface DeveloperMetrics {
  total_commits: number;
  total_prs: number;
  merged_prs: number;
  total_issues_opened: number;
  total_issues_closed: number;
  review_comments: number;
  active_weeks: number;
  repos_contributed: number;
  total_stars: number;
  total_forks: number;
  avg_pr_lines_changed: number;
  avg_issue_response_hours: number;
  commit_trend: string;
  languages: string[];
}

export interface Developer {
  username: string;
  profile: DeveloperProfile;
  metrics: DeveloperMetrics;
  fetched_at: string;
}

export interface IndicatorScore {
  raw: number;
  normalized: number;
  weighted: number;
}

export interface Score {
  username: string;
  weighted_score: number;
  ml_impact_score: number;
  indicator_scores: Record<string, IndicatorScore>;
  percentile: number;
  computed_at: string;
}

export interface RankEntry {
  rank: number;
  username: string;
  score: number;
  ml_score: number;
}

export interface Ranking {
  snapshot_date: string;
  rankings: RankEntry[];
  total_developers: number;
  created_at: string;
}

export interface AuthUser {
  id: string;
  username: string;
  display_name: string;
  avatar_url: string;
  email: string;
}

interface RequestOptions {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  token?: string | null;
}

interface ApiError {
  error: string;
}

class ApiClient {
  private baseURL: string;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
  }

  private getToken(): string | null {
    if (typeof window === "undefined") return null;
    return localStorage.getItem("token");
  }

  private async request<T>(
    endpoint: string,
    options: RequestOptions = {}
  ): Promise<T> {
    const { method = "GET", body, headers = {}, token } = options;
    const authToken = token ?? this.getToken();

    const config: RequestInit = {
      method,
      headers: {
        "Content-Type": "application/json",
        ...(authToken ? { Authorization: `Bearer ${authToken}` } : {}),
        ...headers,
      },
    };

    if (body) {
      config.body = JSON.stringify(body);
    }

    const response = await fetch(`${this.baseURL}${endpoint}`, config);

    if (!response.ok) {
      const errorData: ApiError = await response.json().catch(() => ({
        error: `HTTP ${response.status}`,
      }));
      throw new Error(errorData.error || `Request failed: ${response.status}`);
    }

    return response.json();
  }

  // ---- Auth ----
  getLoginURL(): string {
    return `${this.baseURL}/auth/github/login`;
  }

  async getCurrentUser() {
    return this.request<AuthUser>("/auth/me");
  }

  async logout() {
    return this.request("/auth/logout", { method: "POST" });
  }

  // ---- Developers ----
  async fetchDeveloper(username: string) {
    return this.request<Developer>(
      `/api/developers/${username}/fetch`,
      { method: "POST" }
    );
  }

  async getDeveloper(username: string) {
    return this.request<Developer>(
      `/api/developers/${username}`
    );
  }

  async listDevelopers() {
    return this.request<Developer[]>("/api/developers");
  }

  // ---- Scoring ----
  async computeScore(username: string) {
    return this.request<Score>(
      `/api/developers/${username}/score`,
      { method: "POST" }
    );
  }

  async getScore(username: string) {
    return this.request<Score>(
      `/api/developers/${username}/score`
    );
  }

  // ---- Ranking ----
  async compareDevelopers(usernames: string[]) {
    return this.request<Ranking>("/api/compare", {
      method: "POST",
      body: { usernames },
    });
  }

  async evaluateOrg(org: string) {
    return this.request<Ranking>(
      `/api/orgs/${org}/evaluate`
    );
  }

  async getRankings(page = 1, pageSize = 20) {
    return this.request<Ranking>(
      `/api/rankings?page=${page}&page_size=${pageSize}`
    );
  }
}

export const api = new ApiClient(API_BASE_URL);
export default api;
