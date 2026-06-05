"use client";

import Link from "next/link";
import { useEffect, useRef, useCallback } from "react";
import { useGuest } from "../../app/context/GuestContext";
import { useAuth } from "../../app/context/AuthContext";

declare global {
  interface Window {
    google?: {
      accounts: {
        id: {
          initialize: (config: Record<string, unknown>) => void;
          renderButton: (
            el: HTMLElement,
            config: Record<string, unknown>
          ) => void;
          prompt: () => void;
        };
      };
    };
  }
}

const GOOGLE_CLIENT_ID =
  process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID ?? "";

export default function Navbar() {
  const { credit } = useGuest();
  const { user, isLoggedIn, login, logout } = useAuth();
  const googleBtnRef = useRef<HTMLDivElement>(null);

  const handleCredentialResponse = useCallback(
    async (response: { credential: string }) => {
      try {
        await login(response.credential);
      } catch {
        // login failed — ignored silently
      }
    },
    [login]
  );

  useEffect(() => {
    if (isLoggedIn || !googleBtnRef.current) return;

    const initGoogle = () => {
      if (!window.google) return;
      window.google.accounts.id.initialize({
        client_id: GOOGLE_CLIENT_ID,
        callback: handleCredentialResponse,
      });
      window.google.accounts.id.renderButton(googleBtnRef.current!, {
        theme: "filled_black",
        size: "medium",
        shape: "pill",
        text: "signin_with",
      });
    };

    // Google script might not be loaded yet
    if (window.google) {
      initGoogle();
    } else {
      const interval = setInterval(() => {
        if (window.google) {
          initGoogle();
          clearInterval(interval);
        }
      }, 100);
      return () => clearInterval(interval);
    }
  }, [isLoggedIn, handleCredentialResponse]);

  return (
    <nav className="relative w-full h-16 flex items-center border-b border-zinc-800 bg-[#111111] px-6 z-50">
      {/* Logo */}
      <Link
        href="/"
        className="text-lg font-bold text-white hover:text-red-400 transition-colors"
        style={{ textShadow: "0 0 12px rgba(239,68,68,0.4)" }}
      >
        เครียดโว้ยยยย 😤
      </Link>

      {/* Right side */}
      <div className="ml-auto flex items-center gap-3">
        {/* Credit pill */}
        {credit !== null && (
          credit <= 0 ? (
            <span className="text-xs px-3 py-1 rounded-full bg-red-900/50 border border-red-700 text-red-400 font-semibold">
              เครดิตหมด
            </span>
          ) : (
            <span className="text-xs px-3 py-1 rounded-full bg-zinc-800 border border-zinc-700 text-zinc-300">
              🔥 {credit} ครั้ง
            </span>
          )
        )}

        {isLoggedIn && user ? (
          /* Logged-in state */
          <div className="flex items-center gap-3">
            {user.avatarUrl && (
              <img
                src={user.avatarUrl}
                alt={user.username}
                className="w-8 h-8 rounded-full border border-zinc-600"
              />
            )}
            <span className="text-sm text-zinc-300 hidden sm:inline">
              {user.username}
            </span>
            <button
              onClick={logout}
              className="px-3 py-1.5 border border-zinc-700 rounded-lg text-zinc-400 hover:text-white hover:border-zinc-500 transition-colors text-xs"
            >
              ออกจากระบบ
            </button>
          </div>
        ) : (
          /* Guest state — Google login button */
          <div ref={googleBtnRef} />
        )}
      </div>
    </nav>
  );
}
