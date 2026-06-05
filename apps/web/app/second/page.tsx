"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";

// Per level + target fun messages
const LEVEL_MESSAGES: Record<number, string> = {
  1: "สุภาพจนพระสงฆ์ยังอนุโมทนา 🙏✨",
  2: "ข้อความนี้ส่งให้ CEO ได้เลย 💼",
  3: "สุภาพพอดี ไม่มากไม่น้อย 😌",
  4: "สุภาพ... แต่คนอ่านจะรู้สึกอะไรบางอย่าง 😈",
  5: "ดูดี... แต่แฝงดาบทุกตัวอักษร 💀🔪",
};

const TARGET_LABELS: Record<string, { emoji: string; label: string }> = {
  boss:    { emoji: "👔", label: "เจ้านาย"  },
  client:  { emoji: "💼", label: "ลูกค้า"   },
  teacher: { emoji: "📚", label: "อาจารย์"  },
  friend:  { emoji: "🫂", label: "เพื่อน"   },
};

const GLOW_COLORS: Record<number, string> = {
  1: "rgba(34,197,94,0.15)",
  2: "rgba(59,130,246,0.15)",
  3: "rgba(234,179,8,0.12)",
  4: "rgba(249,115,22,0.15)",
  5: "rgba(239,68,68,0.15)",
};

const BORDER_CLASSES: Record<number, string> = {
  1: "border-green-800",
  2: "border-blue-800",
  3: "border-yellow-800",
  4: "border-orange-800",
  5: "border-red-800",
};

export default function Second() {
  const [result,  setResult]  = useState<string | null>(null);
  const [level,   setLevel]   = useState<number>(3);
  const [target,  setTarget]  = useState<string>("boss");
  const [loading, setLoading] = useState(true);
  const [copied,  setCopied]  = useState(false);
  const [sparkle, setSparkle] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const stored       = sessionStorage.getItem("translationResult");
    const storedLevel  = sessionStorage.getItem("translationLevel");
    const storedTarget = sessionStorage.getItem("translationTarget");
    if (stored)       { const p = JSON.parse(stored); setResult(p.result ?? null); }
    if (storedLevel)  setLevel(parseInt(storedLevel, 10));
    if (storedTarget) setTarget(storedTarget);
    setLoading(false);
  }, []);

  const handleCopy = async () => {
    if (!result) return;
    try {
      await navigator.clipboard.writeText(result);
      setCopied(true);
      setSparkle(true);
      setTimeout(() => setCopied(false), 2500);
      setTimeout(() => setSparkle(false), 1000);
    } catch {}
  };

  const targetInfo = TARGET_LABELS[target] ?? { emoji: "💬", label: target };

  return (
    <div className="relative flex-1 flex flex-col items-center justify-center w-full px-4 py-12 overflow-hidden bg-[#111111]">
      {/* Background glow */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{ background: `radial-gradient(ellipse 70% 50% at 50% 55%, ${GLOW_COLORS[level] ?? GLOW_COLORS[3]} 0%, transparent 70%)` }}
      />

      {loading ? (
        <p className="text-zinc-400 text-lg">กำลังโหลด...</p>
      ) : result ? (
        <div className="relative w-full max-w-2xl fade-up">
          {/* Target + level badge */}
          <div className="flex items-center justify-center gap-3 mb-4">
            <span className="px-3 py-1 rounded-full bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm">
              {targetInfo.emoji} ส่งให้{targetInfo.label}
            </span>
            <span className="text-zinc-500 text-sm">•</span>
            <span className="text-sm text-zinc-400">
              {LEVEL_MESSAGES[level] ?? "พร้อมส่งแล้ว ✅"}
            </span>
          </div>

          {/* Result card */}
          <div
            className={`relative rounded-3xl border ${BORDER_CLASSES[level] ?? "border-zinc-700"} bg-zinc-900/90 backdrop-blur-sm p-8 shadow-2xl`}
          >
            <p className="text-xl md:text-2xl text-zinc-100 leading-relaxed whitespace-pre-wrap">
              {result}
            </p>

            {/* Sparkle overlay on copy */}
            {sparkle && (
              <div className="absolute inset-0 rounded-3xl pointer-events-none overflow-hidden">
                {Array.from({ length: 12 }).map((_, i) => (
                  <span
                    key={i}
                    className="absolute text-yellow-300 animate-ping"
                    style={{
                      left: `${10 + Math.random() * 80}%`,
                      top:  `${10 + Math.random() * 80}%`,
                      fontSize: `${10 + Math.random() * 14}px`,
                      animationDuration: `${0.5 + Math.random() * 0.5}s`,
                      animationDelay:    `${Math.random() * 0.3}s`,
                    }}
                  >
                    ✨
                  </span>
                ))}
              </div>
            )}
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
                  ? "bg-green-600 text-white scale-105"
                  : "bg-zinc-800 hover:bg-zinc-700 text-zinc-200"
              }`}
            >
              {copied ? "✓ คัดลอกแล้ว! ✨" : "คัดลอก 📋"}
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
