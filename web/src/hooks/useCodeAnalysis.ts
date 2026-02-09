import { useState, useCallback } from "react";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface AnalyzeCodeRequest {
  code: string;
  language: string;
  filename?: string;
}

interface AnalyzeCodeResponse {
  analysis: string;
  vulnerabilities: string[];
  security_score: number;
  recommendations: string;
  vulnerability_count: number;
}

interface QuickScanResponse {
  vulnerabilities: string[];
  vulnerability_count: number;
  security_score: number;
  scan_type: string;
}

interface CompareCodeRequest {
  code1: string;
  code2: string;
  language1: string;
  language2: string;
}

interface CompareCodeResponse {
  similarity: number;
  similarity_percent: number;
  is_duplicate: boolean;
  common_keywords: string[];
}

export interface CodeAnalysisHook {
  loading: boolean;
  error: string | null;
  analyzeCode: (request: AnalyzeCodeRequest) => Promise<AnalyzeCodeResponse>;
  quickScan: (request: AnalyzeCodeRequest) => Promise<QuickScanResponse>;
  compareCode: (request: CompareCodeRequest) => Promise<CompareCodeResponse>;
  clearError: () => void;
}

export const useCodeAnalysis = (): CodeAnalysisHook => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Analyze code with AI
  const analyzeCode = useCallback(async (request: AnalyzeCodeRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response =await api.post<AnalyzeCodeResponse>(
        API_ENDPOINTS.CODE.ANALYZE,
        request
      );
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Quick vulnerability scan
  const quickScan = useCallback(async (request: AnalyzeCodeRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post<QuickScanResponse>(
        API_ENDPOINTS.CODE.QUICK_SCAN,
        request
      );
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Compare two code snippets
  const compareCode = useCallback(async (request: CompareCodeRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post<CompareCodeResponse>(
        API_ENDPOINTS.CODE.COMPARE,
        request
      );
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Clear error
  const clearError = useCallback(() => {
    setError(null);
  }, []);

  return {
    loading,
    error,
    analyzeCode,
    quickScan,
    compareCode,
    clearError,
  };
};
