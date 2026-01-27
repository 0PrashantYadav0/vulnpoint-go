import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { setAuthToken } from '@/lib/api';

/**
 * GitHub OAuth Callback Component
 * 
 * This component handles the OAuth callback from GitHub.
 * It extracts the token from the URL, stores it, and redirects to the dashboard.
 */
export default function AuthCallback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  useEffect(() => {
    // Get token from query parameters
    const token = searchParams.get('token');
    const userJson = searchParams.get('user');
    
    if (token) {
      // Store the JWT token
      setAuthToken(token);
      
      // Store user data if provided
      if (userJson) {
        try {
          const user = JSON.parse(decodeURIComponent(userJson));
          localStorage.setItem('user', JSON.stringify(user));
        } catch (e) {
          console.error('Failed to parse user data:', e);
        }
      }
      
      // Redirect to dashboard
      navigate('/dashboard');
    } else {
      // No token, redirect to home page
      navigate('/');
    }
  }, [searchParams, navigate]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <h2 className="text-2xl font-semibold mb-2">Completing authentication...</h2>
        <p className="text-gray-600">Please wait while we log you in.</p>
      </div>
    </div>
  );
}
