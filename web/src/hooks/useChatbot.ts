import { useState, useCallback } from "react";
import api, { getErrorMessage } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

interface ChatResponse {
  response: string;
  conversation_id?: string;
}

interface ExplainVulnerabilityRequest {
  vulnerability_type: string;
  context?: string;
}

interface RemediationRequest {
  vulnerability_type: string;
  code_snippet?: string;
  language?: string;
}

interface SecurityQuestionRequest {
  question: string;
  category?: string;
}

export interface ChatbotHook {
  loading: boolean;
  error: string | null;
  chat: (message: string, conversationHistory?: ChatMessage[]) => Promise<ChatResponse>;
  explainVulnerability: (request: ExplainVulnerabilityRequest) => Promise<any>;
  getRemediation: (request: RemediationRequest) => Promise<any>;
  askQuestion: (request: SecurityQuestionRequest) => Promise<any>;
  clearError: () => void;
}

export const useChatbot = (): ChatbotHook => {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Chat with AI
  const chat = useCallback(async (message: string, conversationHistory?: ChatMessage[]) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post<ChatResponse>(API_ENDPOINTS.CHATBOT.CHAT, {
        message,
        conversation_history: conversationHistory || [],
      });
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Explain vulnerability
  const explainVulnerability = useCallback(async (request: ExplainVulnerabilityRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post(API_ENDPOINTS.CHATBOT.EXPLAIN, request);
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Get remediation steps
  const getRemediation = useCallback(async (request: RemediationRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post(API_ENDPOINTS.CHATBOT.REMEDIATE, request);
      return response.data;
    } catch (err) {
      const errorMessage = getErrorMessage(err);
      setError(errorMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Ask security question
  const askQuestion = useCallback(async (request: SecurityQuestionRequest) => {
    setLoading(true);
    setError(null);
    try {
      const response = await api.post(API_ENDPOINTS.CHATBOT.ASK, request);
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
    chat,
    explainVulnerability,
    getRemediation,
    askQuestion,
    clearError,
  };
};
