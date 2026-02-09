import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { setAuthToken } from "@/lib/api";
import api from "@/lib/api";
import { API_ENDPOINTS } from "@/lib/apiEndpoints";

/**
 * GitHub OAuth Callback Component
 *
 * Handles the OAuth callback from GitHub, stores the token, fetches user,
 * and redirects to the dashboard.
 */
export default function AuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = searchParams.get("token");

    if (!token) {
      navigate("/");
      return;
    }

    setAuthToken(token);

    const completeAuth = async () => {
      try {
        const response = await api.get(API_ENDPOINTS.AUTH.USER);
        const userData = response.data?.data;
        if (userData) {
          localStorage.setItem("user", JSON.stringify(userData));
        }
        navigate("/dashboard", { replace: true });
      } catch (err) {
        console.error("Failed to fetch user after auth:", err);
        navigate("/dashboard", { replace: true });
      }
    };

    completeAuth();
  }, [searchParams, navigate]);

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center text-red-600 dark:text-red-400">
          <p>{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-center min-h-screen bg-white dark:bg-zinc-950">
      <div className="text-center">
        <h2 className="text-2xl font-semibold mb-2 text-gray-900 dark:text-white">
          Completing authentication...
        </h2>
        <p className="text-gray-600 dark:text-gray-400">
          Please wait while we log you in.
        </p>
        <div className="mt-4 flex justify-center">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-blue-600 border-t-transparent" />
        </div>
      </div>
    </div>
  );
}
