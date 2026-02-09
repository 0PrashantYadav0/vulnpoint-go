import { User } from "@/types";
import { useState, useEffect, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import api, { removeAuthToken, getAuthToken } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

interface AuthResponse {
  success: boolean;
  data?: User;
}

interface AuthURLResponse {
  success: boolean;
  data?: { url: string; state: string };
}

const useAuth = () => {
  const [user, setUser] = useState<User | null>(() => {
    const storedUser = localStorage.getItem("user");
    return storedUser ? JSON.parse(storedUser) : null;
  });
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const fetchUser = useCallback(async () => {
    if (!getAuthToken()) {
      setUser(null);
      setLoading(false);
      return;
    }
    try {
      const response = await api.get<AuthResponse>(API_ENDPOINTS.AUTH.USER);
      const userData = response.data?.data;
      if (userData) {
        setUser(userData);
        localStorage.setItem("user", JSON.stringify(userData));
      } else {
        setUser(null);
        localStorage.removeItem("user");
        removeAuthToken();
      }
    } catch {
      setUser(null);
      localStorage.removeItem("user");
      removeAuthToken();
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  const loginWithGithub = async () => {
    try {
      const response = await api.get<AuthURLResponse>(API_ENDPOINTS.AUTH.GITHUB_URL);
      const authUrl = response.data?.data?.url;
      if (authUrl) {
        window.location.href = authUrl;
      } else {
        console.error("No auth URL received - ensure GitHub OAuth is configured");
      }
    } catch (error) {
      console.error("Failed to initiate GitHub login:", error);
      throw error;
    }
  };

  const logout = useCallback(async () => {
    try {
      await api.post(API_ENDPOINTS.AUTH.LOGOUT);
    } catch {
      // Ignore logout API errors
    } finally {
      setUser(null);
      localStorage.removeItem("user");
      removeAuthToken();
      navigate("/");
    }
  }, [navigate]);

  return { user, loading, loginWithGithub, logout, refetchUser: fetchUser };
};

export default useAuth;
