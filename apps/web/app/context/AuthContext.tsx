"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

interface User {
  id: string;
  email: string;
  username: string;
  avatarUrl: string | null;
  credit: number;
  memberType: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  isLoggedIn: boolean;
  loading: boolean;
  login: (idToken: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthState>({
  user: null,
  token: null,
  isLoggedIn: false,
  loading: true,
  login: async () => {},
  logout: () => {},
  refreshUser: async () => {},
});

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  // Refresh user data from /api/v1/user/me
  const refreshUser = useCallback(async () => {
    const storedToken = token ?? localStorage.getItem("jwt");
    if (!storedToken) {
      setLoading(false);
      return;
    }
    try {
      const res = await fetch(`${API_BASE}/api/v1/user/me`, {
        headers: { Authorization: `Bearer ${storedToken}` },
        credentials: "include",
      });
      if (!res.ok) {
        // Token invalid or expired — clear everything
        localStorage.removeItem("jwt");
        setToken(null);
        setUser(null);
        setLoading(false);
        return;
      }
      const data = await res.json();
      setUser(data);
      setToken(storedToken);
    } catch {
      // Network error — keep existing state
    } finally {
      setLoading(false);
    }
  }, [token]);

  // Login with Google ID token
  const login = useCallback(async (idToken: string) => {
    const res = await fetch(`${API_BASE}/api/v1/auth/google`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({ idToken }),
    });

    if (!res.ok) {
      throw new Error("Login failed");
    }

    const data = await res.json();
    localStorage.setItem("jwt", data.token);
    setToken(data.token);
    setUser(data.user);
  }, []);

  const logout = useCallback(() => {
    localStorage.removeItem("jwt");
    setToken(null);
    setUser(null);
  }, []);

  // On mount, check for existing JWT
  useEffect(() => {
    const stored = localStorage.getItem("jwt");
    if (stored) {
      setToken(stored);
    } else {
      setLoading(false);
    }
  }, []);

  // When token changes, refresh user
  useEffect(() => {
    if (token) {
      refreshUser();
    }
  }, [token, refreshUser]);

  return (
    <AuthContext.Provider
      value={{
        user,
        token,
        isLoggedIn: !!user,
        loading,
        login,
        logout,
        refreshUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  return useContext(AuthContext);
}
