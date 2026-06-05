"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { useAuth } from "./AuthContext";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

interface GuestState {
  credit: number | null;
  isLoggedIn: boolean;
  refreshCredit: () => Promise<void>;
}

const GuestContext = createContext<GuestState>({
  credit: null,
  isLoggedIn: false,
  refreshCredit: async () => {},
});

export function GuestProvider({ children }: { children: ReactNode }) {
  const [credit, setCredit] = useState<number | null>(null);
  const { user, token, isLoggedIn, loading: authLoading } = useAuth();

  const refreshCredit = useCallback(async () => {
    try {
      if (isLoggedIn && token) {
        // Logged-in user — get credit from /api/v1/user/me
        const res = await fetch(`${API_BASE}/api/v1/user/me`, {
          headers: { Authorization: `Bearer ${token}` },
          credentials: "include",
        });
        if (res.ok) {
          const data = await res.json();
          setCredit(data.credit ?? null);
        }
      } else {
        // Guest — get credit from /guest/status
        const res = await fetch(`${API_BASE}/guest/status`, {
          credentials: "include",
        });
        if (res.ok) {
          const data = await res.json();
          setCredit(data.credit ?? null);
        }
      }
    } catch {
      // silently ignore network errors
    }
  }, [isLoggedIn, token]);

  // Update credit when user data changes
  useEffect(() => {
    if (user) {
      setCredit(user.credit);
    }
  }, [user]);

  // Fetch credit on first mount and when auth state resolves
  useEffect(() => {
    if (!authLoading) {
      refreshCredit();
    }
  }, [authLoading, refreshCredit]);

  return (
    <GuestContext.Provider value={{ credit, isLoggedIn, refreshCredit }}>
      {children}
    </GuestContext.Provider>
  );
}

export function useGuest() {
  return useContext(GuestContext);
}
