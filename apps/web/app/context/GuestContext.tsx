"use client";

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from "react";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

interface GuestState {
  credit: number | null;   // null = not loaded yet
  refreshCredit: () => Promise<void>;
}

const GuestContext = createContext<GuestState>({
  credit: null,
  refreshCredit: async () => {},
});

export function GuestProvider({ children }: { children: ReactNode }) {
  const [credit, setCredit] = useState<number | null>(null);

  const refreshCredit = useCallback(async () => {
    try {
      const res = await fetch(`${API_BASE}/guest/status`, {
        credentials: "include", // send & receive the guest_id cookie
      });
      if (!res.ok) return;
      const data = await res.json();
      setCredit(data.credit ?? null);
    } catch {
      // network error — silently ignore, credit stays null
    }
  }, []);

  // Fetch on first mount (triggers cookie creation via middleware)
  useEffect(() => {
    refreshCredit();
  }, [refreshCredit]);

  return (
    <GuestContext.Provider value={{ credit, refreshCredit }}>
      {children}
    </GuestContext.Provider>
  );
}

export function useGuest() {
  return useContext(GuestContext);
}
