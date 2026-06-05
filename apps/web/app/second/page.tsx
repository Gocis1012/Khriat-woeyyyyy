"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

export default function Second() {
  const [result, setResult] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const stored = sessionStorage.getItem("translationResult");
    if (stored) {
      const parsed = JSON.parse(stored);
      setResult(parsed.result ?? null);
    }
    setLoading(false);
  }, []);

  const handleCopy = async () => {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  };

  return (
    <div className="relative flex-1 flex flex-col items-center justify-center w-full px-4 py-12 overflow-hidden bg-[#111111]">

      {/* Soft green glow — calm after the storm */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse 70% 50% at 50% 55%, rgba(22,163,74,0.15) 0%, transparent 70%)",
        }}
      />

      {loading ? (
        <p className="text-zinc-400 text-lg">กำลังโหลด...</p>
      ) : result ? (
        <div className="relative w-full max-w-2xl fade-up">
          {/* Label */}
          <p className="text-xs uppercase tracking-widest text-zinc-500 mb-3 text-center">
            ข้อความพร้อมส่งแล้ว ✅
          </p>

          {/* Result card */}
          <div className="rounded-3xl border border-zinc-700 bg-zinc-900/90 backdrop-blur-sm p-8 shadow-2xl">
            <p className="text-xl md:text-2xl text-zinc-100 leading-relaxed whitespace-pre-wrap">
              {result}
            </p>
          </div>

          {/* Actions */}
          <div className="flex items-center justify-between mt-6">
            <button
              onClick={() => router.push("/")}
              className="px-5 py-2 rounded-xl border border-zinc-700 text-zinc-400 hover:text-white hover:border-zinc-500 transition-colors text-sm"
            >
              ← บ่นใหม่
            </button>

            <button
              onClick={handleCopy}
              className={`px-6 py-2 rounded-xl text-sm font-semibold transition-all ${
                copied
                  ? "bg-green-600 text-white"
                  : "bg-zinc-800 hover:bg-zinc-700 text-zinc-200"
              }`}
            >
              {copied ? "✓ คัดลอกแล้ว!" : "คัดลอก"}
            </button>
          </div>
        </div>
      ) : (
        <div className="relative text-center fade-up">
          <p className="text-zinc-400 text-lg mb-4">ไม่พบผลลัพท์</p>
          <button
            onClick={() => router.push("/")}
            className="px-6 py-2 rounded-xl bg-zinc-800 hover:bg-zinc-700 text-zinc-200 text-sm transition-colors"
          >
            ← กลับไปบ่น
          </button>
        </div>
      )}
    </div>
  );
}
