"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useRef, useCallback, useState } from "react";
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

const GOOGLE_CLIENT_ID = process.env.NEXT_PUBLIC_GOOGLE_CLIENT_ID ?? "";

export default function Navbar() {
  const pathname = usePathname();
  const { credit } = useGuest();
  const { user, isLoggedIn, login, logout } = useAuth();
  const [showLogout, setShowLogout] = useState(false);
  const initialized = useRef(false);

  const handleCredential = useCallback(
    async (response: { credential: string }) => {
      try {
        await login(response.credential);
      } catch {
        /* ignore */
      }
    },
    [login]
  );

  // Initialize Google Identity once (One Tap on click)
  useEffect(() => {
    if (initialized.current || isLoggedIn || !GOOGLE_CLIENT_ID) return;

    const init = () => {
      if (!window.google) return false;
      window.google.accounts.id.initialize({
        client_id: GOOGLE_CLIENT_ID,
        callback: handleCredential,
      });
      initialized.current = true;
      return true;
    };

    if (!init()) {
      const iv = setInterval(() => {
        if (init()) clearInterval(iv);
      }, 150);
      return () => clearInterval(iv);
    }
  }, [isLoggedIn, handleCredential]);

  const handleLoginClick = () => {
    if (!GOOGLE_CLIENT_ID) {
      console.warn("[Navbar] NEXT_PUBLIC_GOOGLE_CLIENT_ID is not set");
      return;
    }
    window.google?.accounts.id.prompt();
  };

  return (
    <>
      <nav className="relative w-full h-20 flex items-center border-b border-black/10 bg-white px-6 md:px-10">
        {/* Logo */}
        <Link
          href="/"
          className="text-3xl md:text-[40px] font-bold text-black hover:opacity-80 transition-opacity"
        >
          เครียดโว้ยยยย
        </Link>

        {/* Center nav */}
        <div className="absolute left-1/2 -translate-x-1/2 hidden sm:flex items-center gap-10">
          <Link
            href="/"
            className={`text-xl font-bold transition-colors ${
              pathname === "/" ? "text-[#64579f]" : "text-black hover:text-[#64579f]"
            }`}
          >
            หน้าแรก
          </Link>
          {/* Placeholder — no About page yet, so this is inert */}
          <span
            className="text-xl font-bold text-black/30 cursor-default select-none"
            aria-disabled="true"
          >
            เกี่ยวกับเรา
          </span>
        </div>

        {/* Right side */}
        <div className="ml-auto flex items-center gap-3">
          {/* Credit pill */}
          {credit !== null &&
            (credit <= 0 ? (
              <span className="text-sm px-3 py-1 rounded-full bg-red-50 border border-red-300 text-red-500 font-bold">
                เครดิตหมด
              </span>
            ) : (
              <span className="text-sm px-3 py-1 rounded-full bg-orange-50 border border-[#ff8055]/40 text-[#ff7b00] font-bold">
                🔥 {credit} ครั้ง
              </span>
            ))}

          {/* Orange circle: avatar (logged in) or login (guest) */}
          {isLoggedIn && user ? (
            <button
              onClick={() => setShowLogout(true)}
              aria-label="โปรไฟล์"
              className="h-12 w-12 rounded-full bg-[#ff8055] overflow-hidden border-2 border-[#ff8055] hover:opacity-90 transition-opacity flex items-center justify-center"
            >
              {user.avatarUrl ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img
                  src={user.avatarUrl}
                  alt={user.username}
                  referrerPolicy="no-referrer"
                  className="h-full w-full object-cover"
                />
              ) : (
                <span className="text-white text-lg font-bold">
                  {user.username.charAt(0).toUpperCase()}
                </span>
              )}
            </button>
          ) : (
            <button
              onClick={handleLoginClick}
              aria-label="เข้าสู่ระบบ"
              title="เข้าสู่ระบบด้วย Google"
              className="h-12 w-12 rounded-full bg-[#ff8055] hover:bg-[#ff6a3c] transition-colors flex items-center justify-center shadow-sm"
            >
              <svg
                width="22"
                height="22"
                viewBox="0 0 24 24"
                fill="none"
                stroke="white"
                strokeWidth="2.2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                <circle cx="12" cy="7" r="4" />
              </svg>
            </button>
          )}
        </div>
      </nav>

      {/* ── Logout modal ───────────────────────────────── */}
      {showLogout && (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 backdrop-blur-sm"
          onClick={() => setShowLogout(false)}
        >
          <div
            className="bg-white border border-black/10 rounded-3xl p-6 w-80 flex flex-col gap-4 shadow-2xl"
            onClick={(e) => e.stopPropagation()}
          >
            <h2 className="text-black font-bold text-lg">ออกจากระบบ?</h2>
            <p className="text-black/60">คุณต้องการออกจากระบบใช่ไหม?</p>
            <div className="flex gap-3 mt-2">
              <button
                onClick={() => setShowLogout(false)}
                className="flex-1 px-4 py-2 rounded-full border border-black/30 text-black/70 hover:border-black/60 hover:text-black transition-colors"
              >
                ยกเลิก
              </button>
              <button
                onClick={() => {
                  logout();
                  setShowLogout(false);
                }}
                className="flex-1 px-4 py-2 rounded-full bg-[#ff8055] hover:bg-[#ff6a3c] text-white font-bold transition-colors"
              >
                ยืนยัน
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
