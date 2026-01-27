import { User } from "@/types";
import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import api, { removeAuthToken } from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

const useAuth = () => {
  const [user, setUser] = useState<User | null>(() => {
    const storedUser = localStorage.getItem("user");
    return storedUser ? JSON.parse(storedUser) : null;
  });
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  const fetchUser = async () => {
    try {
      const response = await api.get(API_ENDPOINTS.AUTH.USER);
      const userData = response.data;
      setUser(userData);
      localStorage.setItem("user", JSON.stringify(userData));
    } catch (error) {
      localStorage.removeItem("user");
      removeAuthToken();
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUser();
  }, []);

  const loginWithGithub = async () => {
    if (!user) {
      try {
        // Fetch the GitHub OAuth URL from backend
        const response = await api.get<{ success: boolean; data: { url: string; state: string } }>(
          API_ENDPOINTS.AUTH.GITHUB_URL
        );
        // Backend wraps response in { success: true, data: { url, state } }
        const authUrl = response.data.data?.url;
        
        if (authUrl) {
          // Redirect to GitHub OAuth
          window.location.href = authUrl;
        } else {
          console.error("No auth URL received from backend");
        }
      } catch (error) {
        console.error("Failed to initiate GitHub login:", error);
      }
    }
  };

  const logout = async () => {
    try {
      await api.post(API_ENDPOINTS.AUTH.LOGOUT);
      setUser(null);
      localStorage.removeItem("user");
      removeAuthToken();
      navigate("/");
    } catch (error) {
      console.error("Logout failed", error);
    }
  };

  return { user, loading, loginWithGithub, logout };
};

export default useAuth;
