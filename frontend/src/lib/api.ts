const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

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
    return this.request<{
      id: string;
      username: string;
      display_name: string;
      avatar_url: string;
      email: string;
    }>("/auth/me");
  }

  async logout() {
    return this.request("/auth/logout", { method: "POST" });
  }

  // ---- Developers ----
  async fetchDeveloper(username: string) {
    return this.request<Record<string, unknown>>(
      `/api/developers/${username}/fetch`,
      { method: "POST" }
    );
  }

  async getDeveloper(username: string) {
    return this.request<Record<string, unknown>>(
      `/api/developers/${username}`
    );
  }

  async listDevelopers() {
    return this.request<Record<string, unknown>[]>("/api/developers");
  }

  // ---- Scoring ----
  async computeScore(username: string) {
    return this.request<Record<string, unknown>>(
      `/api/developers/${username}/score`,
      { method: "POST" }
    );
  }

  async getScore(username: string) {
    return this.request<Record<string, unknown>>(
      `/api/developers/${username}/score`
    );
  }

  // ---- Ranking ----
  async compareDevelopers(usernames: string[]) {
    return this.request<Record<string, unknown>>("/api/compare", {
      method: "POST",
      body: { usernames },
    });
  }

  async evaluateOrg(org: string) {
    return this.request<Record<string, unknown>>(
      `/api/orgs/${org}/evaluate`
    );
  }

  async getRankings(page = 1, pageSize = 20) {
    return this.request<Record<string, unknown>>(
      `/api/rankings?page=${page}&page_size=${pageSize}`
    );
  }
}

export const api = new ApiClient(API_BASE_URL);
export default api;
