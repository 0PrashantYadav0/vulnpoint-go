import { useState, useCallback } from "react";
import { Repository } from "@/types";
import useAuth from "./useAuth";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface ApiResponse {
  success?: boolean;
  count?: number;
  data?: Repository[];
  message?: string;
  error?: string;
}

const useGitHub = () => {
  const [repos, setRepos] = useState<Repository[]>([]);
  const [repoFiles, setRepoFiles] = useState([]);
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
    return repositories.sort((a, b) => {
      const dateA = a.lastUpdated ? new Date(a.lastUpdated).getTime() : 0;
      const dateB = b.lastUpdated ? new Date(b.lastUpdated).getTime() : 0;
      return dateB - dateA;
    });
  };

  const fetchRepositories = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.get<ApiResponse>(API_ENDPOINTS.GITHUB.REPOSITORIES);
      
      const rawRepos = response.data.data || [];
      
      // Remove duplicates and sort
      const uniqueRepos = removeDuplicateRepos(rawRepos);
      const sortedRepos = sortRepositories(uniqueRepos);
      
      setRepos(sortedRepos);
      
      // Store repository names for caching
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

  const { user } = useAuth();
  
  const fetchRepositoryContents = useCallback(async (repo: string) => {
    const owner = user?.username;

    if (!owner) {
      setError("User information not available");
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await api.get(
        API_ENDPOINTS.GITHUB.REPOSITORY_FILES(owner, repo)
      );
      
      const repoFiles = response.data;
      setRepoFiles(repoFiles);
    } catch (err: unknown) {
      console.error("Failed to fetch repository contents:", err);
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  }, [user]);

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
