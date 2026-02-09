import { useState, useCallback } from "react";
import { Repository } from "@/types";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface ApiResponse {
  success?: boolean;
  count?: number;
  data?: Repository[];
  message?: string;
  error?: string;
}

interface RepoFile {
  name?: string;
  path: string;
  type?: string;
  content?: string;
}

const useGitHub = () => {
  const [repos, setRepos] = useState<Repository[]>([]);
  const [repoFiles, setRepoFiles] = useState<RepoFile[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Remove duplicate repositories based on name and owner
  const removeDuplicateRepos = (repositories: Repository[]): Repository[] => {
    const seen = new Set();
    return repositories.filter((repo) => {
      const key = `${repo.owner}/${repo.name}`;
      if (seen.has(key)) {
        return false;
      }
      seen.add(key);
      return true;
    });
  };

  // Sort repositories by last updated date (newest first)
  const sortRepositories = (repositories: Repository[]): Repository[] => {
    return [...repositories].sort((a, b) => {
      const dateA = (a.lastUpdated || a.updated_at)
        ? new Date(a.lastUpdated || a.updated_at!).getTime()
        : 0;
      const dateB = (b.lastUpdated || b.updated_at)
        ? new Date(b.lastUpdated || b.updated_at!).getTime()
        : 0;
      return dateB - dateA;
    });
  };

  // Map backend response to frontend Repository format
  const mapRepoResponse = (r: Record<string, unknown>): Repository => ({
    id: String(r.id ?? ""),
    name: String(r.name ?? ""),
    owner: String(r.owner ?? ""),
    full_name: r.full_name as string,
    description: r.description as string,
    language: r.language as string,
    url: (r.html_url || r.url) as string,
    html_url: r.html_url as string,
    private: Boolean(r.is_private ?? r.private ?? false),
    is_private: Boolean(r.is_private ?? r.private ?? false),
    lastUpdated: (r.updated_at || r.lastUpdated) as string,
    updated_at: r.updated_at as string,
  });

  const fetchRepositories = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.get<ApiResponse>(
        API_ENDPOINTS.GITHUB.REPOSITORIES
      );
      const rawRepos = response.data?.data ?? response.data ?? [];
      const mappedRepos = (Array.isArray(rawRepos)
        ? rawRepos
        : []
      ).map((r: Record<string, unknown>) => mapRepoResponse(r));

      const uniqueRepos = removeDuplicateRepos(mappedRepos);
      const sortedRepos = sortRepositories(uniqueRepos);

      setRepos(sortedRepos);

      const repoNameList = sortedRepos.map((repo: Repository) => repo.name);
      localStorage.setItem("repos", JSON.stringify(repoNameList));
    } catch (err: unknown) {
      console.error("Failed to fetch repositories:", err);
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const fetchRepositoryContents = useCallback(
    async (owner: string, repo: string) => {
      if (!owner) {
        setError("Repository owner not available");
        return;
      }

      setLoading(true);
      setError(null);
      try {
        const response = await api.get(
          API_ENDPOINTS.GITHUB.REPOSITORY_FILES(owner, repo)
        );
        const rawData = response.data?.data ?? response.data;
        const files = Array.isArray(rawData) ? rawData : [];
        setRepoFiles(files as RepoFile[]);
      } catch (err: unknown) {
        console.error("Failed to fetch repository contents:", err);
        const errorMessage = getErrorMessage(err);
        setError(errorMessage);
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // Clear error state
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    repos,
    repoFiles,
    loading,
    error,
    fetchRepositories,
    fetchRepositoryContents,
    clearError,
  };
};

export default useGitHub;
