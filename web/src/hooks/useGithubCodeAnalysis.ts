import { useState, useCallback } from "react";
import useAuth from "@/hooks/useAuth";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface RepoFile {
  path: string;
  content: string;
  isBinary: boolean;
}

interface CodeAnalysisResponse {
  response: {
    summary: string;
    key_features: string[];
    potential_issues: string[];
    best_practices: string[];
  };
}

interface UseGithubCodeAnalysisReturn {
  repoFiles: RepoFile[];
  analysis: CodeAnalysisResponse | null;
  loadingRepo: boolean;
  loadingAnalysis: boolean;
  error: string | null;
  fetchRepositoryContents: (repo: string) => Promise<RepoFile[]>;
  analyzeCode: (question: string) => Promise<CodeAnalysisResponse | undefined>;
}

export const useGithubCodeAnalysis = (): UseGithubCodeAnalysisReturn => {
  const [repoFiles, setRepoFiles] = useState<RepoFile[]>([]);
  const [analysis, setAnalysis] = useState<CodeAnalysisResponse | null>(null);
  const [loadingRepo, setLoadingRepo] = useState<boolean>(false);
  const [loadingAnalysis, setLoadingAnalysis] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const { user } = useAuth();

  const fetchRepositoryContents = useCallback(
    async (repo: string) => {
      const owner = user?.username;

      setLoadingRepo(true);
      setError(null);
      try {
        const response = await api.get<RepoFile[]>(
          API_ENDPOINTS.GITHUB.REPOSITORY_FILES(owner!, repo)
        );

        const repoFiles = response.data;
        setRepoFiles(repoFiles);
        return repoFiles;
      } catch (err: unknown) {
        const errorMessage = getErrorMessage(err);
        setError(errorMessage);
        throw new Error(errorMessage);
      } finally {
        setLoadingRepo(false);
      }
    },
    [user]
  );

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
        const response = await api.post<CodeAnalysisResponse>(
          API_ENDPOINTS.CODE.ANALYZE,
          {
            code: JSON.stringify(repoFiles),
            language: "multiple",
            filename: question, // Using filename field for the question/context
          }
        );

        const data = response.data;
        setAnalysis(data);
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
