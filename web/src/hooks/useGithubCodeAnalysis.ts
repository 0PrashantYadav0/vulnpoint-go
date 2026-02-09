import { useState, useCallback } from "react";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface RepoFile {
  path: string;
  content: string;
  isBinary: boolean;
}

// Matches backend AnalyzeCodeResponse
export interface CodeAnalysisResponse {
  analysis: string;
  vulnerabilities: string[];
  security_score: number;
  recommendations: string;
  vulnerability_count: number;
}

interface UseGithubCodeAnalysisReturn {
  repoFiles: RepoFile[];
  analysis: CodeAnalysisResponse | null;
  loadingRepo: boolean;
  loadingAnalysis: boolean;
  error: string | null;
  fetchRepositoryContents: () => Promise<RepoFile[]>;
  analyzeCode: (question: string) => Promise<CodeAnalysisResponse | undefined>;
}

export const useGithubCodeAnalysis = (
  owner: string,
  repo: string
): UseGithubCodeAnalysisReturn => {
  const [repoFiles, setRepoFiles] = useState<RepoFile[]>([]);
  const [analysis, setAnalysis] = useState<CodeAnalysisResponse | null>(null);
  const [loadingRepo, setLoadingRepo] = useState<boolean>(false);
  const [loadingAnalysis, setLoadingAnalysis] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const fetchRepositoryContents = useCallback(async () => {
    if (!owner || !repo) return [];

    setLoadingRepo(true);
    setError(null);
    try {
      const response = await api.get(
        API_ENDPOINTS.GITHUB.REPOSITORY_FILES(owner, repo)
      );
      const rawData = response.data?.data ?? response.data;
      const files = Array.isArray(rawData) ? rawData : [];
      setRepoFiles(files);
      return files;
    } catch (err: unknown) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoadingRepo(false);
    }
  }, [owner, repo]);

  const analyzeCode = useCallback(
    async (question: string) => {
      if (!repoFiles.length) {
        setError("No repository files loaded");
        return;
      }

      setLoadingAnalysis(true);
      try {
        // Note: The backend doesn't have a /api/code/query endpoint
        // Using /api/code/analyze instead with the repository files
        const response = await api.post<{ success: boolean; data: CodeAnalysisResponse }>(
          API_ENDPOINTS.CODE.ANALYZE,
          {
            code: JSON.stringify(repoFiles),
            language: "multiple",
            filename: question,
          }
        );

        const data = response.data?.data ?? response.data;
        setAnalysis(data as CodeAnalysisResponse);
        return data;
      } catch (err: unknown) {
        const errorMessage = getErrorMessage(err);
        setError(errorMessage);
        throw new Error(errorMessage);
      } finally {
        setLoadingAnalysis(false);
      }
    },
    [repoFiles]
  );

  return {
    repoFiles,
    analysis,
    loadingRepo,
    loadingAnalysis,
    error,
    fetchRepositoryContents,
    analyzeCode,
  };
};
