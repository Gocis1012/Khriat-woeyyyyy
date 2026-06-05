"use client";

import Link from "next/link";
import { useGuest } from "../../app/context/GuestContext";

export default function Navbar() {
  const { credit } = useGuest();

  return (
    <nav className="relative w-full h-16 flex items-center border-b border-zinc-800 bg-[#111111] px-6">
      {/* Logo */}
      <Link
        href="/"
        className="text-lg font-bold text-white hover:text-red-400 transition-colors"
        style={{ textShadow: "0 0 12px rgba(239,68,68,0.4)" }}
      >
        เครียดโว้ยยยย 😤
      </Link>

      {/* Right side */}
      <div className="ml-auto flex items-center gap-4">
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

        <Link
          href="/login"
          className="px-4 py-1.5 border border-zinc-700 rounded-lg text-zinc-300 hover:border-zinc-500 hover:text-white transition-colors text-sm"
        >
          เข้าสู่ระบบ
        </Link>
        <Link
          href="/signup"
          className="px-4 py-1.5 rounded-lg bg-red-700 hover:bg-red-600 text-white transition-colors text-sm font-semibold"
        >
          สมัครสมาชิก
        </Link>
      </div>
    </nav>
  );
}
