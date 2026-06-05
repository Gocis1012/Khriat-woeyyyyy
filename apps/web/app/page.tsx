// apps/web/app/page.tsx
"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import TextareaAutosize from "react-textarea-autosize";
import { useGuest } from "./context/GuestContext";
import { useAuth } from "./context/AuthContext";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

// Floating anger particles
const PARTICLES = [
  { emoji: "💢", left: "5%",  delay: "0s",   duration: "7s"   },
  { emoji: "🔥", left: "15%", delay: "1.5s", duration: "8s"   },
  { emoji: "😤", left: "25%", delay: "0.3s", duration: "9s"   },
  { emoji: "💥", left: "38%", delay: "2.4s", duration: "7.5s" },
  { emoji: "😡", left: "52%", delay: "0.9s", duration: "10s"  },
  { emoji: "🤬", left: "65%", delay: "1.8s", duration: "6.5s" },
  { emoji: "🔥", left: "75%", delay: "3.2s", duration: "8.5s" },
  { emoji: "💢", left: "88%", delay: "0.5s", duration: "7s"   },
  { emoji: "👊", left: "95%", delay: "2s",   duration: "9.5s" },
];

// Rotating taglines
const TAGLINES = [
  "เจ้านายสั่งงานตอนศุกร์บ่ายโมง 🙃",
  'ลูกค้าบอกว่า "แก้นิดเดียว" 😇',
  "Meeting ที่ควรเป็น Email 📧",
  "ปิดงานแล้วยังไลน์มา 💀",
  'ถูกถามว่า "เสร็จยัง?" ตอน 5 ทุ่ม 🫠',
  "เปิดแชทเห็น... พิมพ์แล้วลบ... 😰",
  'ได้ยินคำว่า "Urgent" วันละ 47 รอบ 🚨',
  "ขอ OT แต่โดนบอกว่าเป็น Teamwork 🤡",
];

// 5 Language levels
const LEVELS = [
  { value: 1, emoji: "🙏", label: "สุภาพสุดๆ",  color: "from-green-600  to-green-700" },
  { value: 2, emoji: "😊", label: "สุภาพ",      color: "from-blue-600   to-blue-700"  },
  { value: 3, emoji: "😐", label: "ปกติ",       color: "from-yellow-600 to-yellow-700"},
  { value: 4, emoji: "😈", label: "แรงนิดๆ",    color: "from-orange-600 to-orange-700"},
  { value: 5, emoji: "💀", label: "แรงสุด",     color: "from-red-600    to-red-700"   },
];

export default function Home() {
  const [input, setInput] = useState("");
  const [level, setLevel] = useState(3);
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [shaking, setShaking] = useState(false);
  const [tagline, setTagline] = useState(TAGLINES[0]);
  const boxRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const { credit, refreshCredit } = useGuest();
  const { token } = useAuth();

  // Rotate taglines
  useEffect(() => {
    const interval = setInterval(() => {
      setTagline((prev) => {
        const idx = TAGLINES.indexOf(prev);
        return TAGLINES[(idx + 1) % TAGLINES.length];
      });
    }, 3500);
    return () => clearInterval(interval);
  }, []);

  const triggerShake = () => {
    setShaking(true);
    setTimeout(() => setShaking(false), 120);
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
    if (e.target.value.length % 5 === 0 && e.target.value.length > 0) {
      triggerShake();
    }
  };

  const handleTranslate = async () => {
    if (!input.trim()) {
      triggerShake();
      setErrorMsg("กรุณากรอกข้อความก่อน");
      return;
    }

    setErrorMsg(null);
    setLoading(true);

    try {
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };
      if (token) {
        headers["Authorization"] = `Bearer ${token}`;
      }

      const response = await fetch(`${API_BASE}/translate`, {
        method: "POST",
        credentials: "include",
        headers,
        body: JSON.stringify({ text: input, level }),
      });

      if (response.status === 402) {
        setErrorMsg("เครดิตหมดแล้ว กรุณาเข้าสู่ระบบเพื่อรับเครดิตเพิ่ม");
        return;
      }
      if (response.status === 401) {
        await refreshCredit();
        setErrorMsg("เซสชันหมดอายุ กรุณาลองใหม่");
        return;
      }
      if (!response.ok) {
        setErrorMsg("เกิดข้อผิดพลาด กรุณาลองใหม่");
        return;
      }

      const data = await response.json();
      await refreshCredit();

      sessionStorage.setItem("translationInput", input);
      sessionStorage.setItem("translationResult", JSON.stringify(data));
      sessionStorage.setItem("translationLevel", String(level));
      router.push("/second");
    } catch {
      setErrorMsg("ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้ กรุณาลองใหม่");
    } finally {
      setLoading(false);
    }
  };

  // Border glow color based on selected level
  const glowColors = [
    "rgba(34,197,94,0.4)",  // green (level 1)
    "rgba(59,130,246,0.4)", // blue (level 2)
    "rgba(234,179,8,0.4)",  // yellow (level 3)
    "rgba(249,115,22,0.4)", // orange (level 4)
    "rgba(239,68,68,0.5)",  // red (level 5)
  ];

  return (
    <div className="relative flex-1 flex flex-col items-center justify-center w-full px-4 py-12 overflow-hidden bg-[#111111]">
      {/* Radial glow — color shifts with level */}
      <div
        className="absolute inset-0 pointer-events-none transition-all duration-700"
        style={{
          background: `radial-gradient(ellipse 70% 50% at 50% 60%, ${glowColors[level - 1]} 0%, transparent 70%)`,
          animation: "rage-pulse 3s ease-in-out infinite",
        }}
      />

      {/* Floating particles */}
      {PARTICLES.map((p, i) => (
        <span
          key={i}
          className="particle"
          style={{
            left: p.left,
            animationDelay: p.delay,
            animationDuration: p.duration,
          }}
        >
          {p.emoji}
        </span>
      ))}

      {/* Title */}
      <h1
        className="relative text-4xl md:text-5xl font-bold text-white mb-2 text-center tracking-wide"
        style={{ textShadow: "0 0 32px rgba(239,68,68,0.6)" }}
      >
        พิมพ์เรื่องที่อยากบ่นให้เราฟังสิ
      </h1>

      {/* Rotating tagline */}
      <p
        key={tagline}
        className="relative text-sm text-zinc-400 mb-4 text-center fade-up"
      >
        {tagline}
      </p>

      {/* Credit */}
      <div className="relative mb-4 flex items-center gap-2 text-sm">
        {credit === null ? (
          <span className="text-zinc-500">กำลังโหลด...</span>
        ) : credit <= 0 ? (
          <span className="px-3 py-1 rounded-full bg-red-900/60 border border-red-600 text-red-300 font-semibold animate-pulse">
            ⚠️ เครดิตหมด — กรุณาเข้าสู่ระบบ
          </span>
        ) : (
          <span className="px-3 py-1 rounded-full bg-zinc-800 border border-zinc-600 text-zinc-300">
            🔥 เครดิตคงเหลือ{" "}
            <span className="font-bold text-white">{credit}</span> ครั้ง
          </span>
        )}
      </div>

      {/* ── Level Selector (Rage Meter) ────────────────────── */}
      <div className="relative flex items-center gap-1.5 mb-6">
        <span className="text-xs text-zinc-500 mr-2">ระดับความแรง:</span>
        {LEVELS.map((l) => (
          <button
            key={l.value}
            onClick={() => setLevel(l.value)}
            className={`
              flex flex-col items-center px-3 py-2 rounded-xl transition-all text-xs font-medium border
              ${
                level === l.value
                  ? `bg-gradient-to-b ${l.color} text-white border-white/20 scale-110 shadow-lg`
                  : "bg-zinc-800/70 text-zinc-400 border-zinc-700 hover:bg-zinc-700"
              }
            `}
            style={
              level === l.value
                ? { boxShadow: `0 0 16px ${glowColors[l.value - 1]}` }
                : {}
            }
          >
            <span className="text-lg mb-0.5">{l.emoji}</span>
            <span>{l.label}</span>
          </button>
        ))}
      </div>

      {/* ── Input Box ──────────────────────────────────────── */}
      <div
        ref={boxRef}
        className={`relative w-full max-w-3xl rounded-3xl border bg-zinc-900/80 backdrop-blur-sm p-6 mb-6 transition-all duration-500 ${shaking ? "shake" : ""}`}
        style={{
          borderColor: glowColors[level - 1],
          boxShadow: `0 0 20px 3px ${glowColors[level - 1]}`,
          animation: "box-glow 2.4s ease-in-out infinite",
        }}
      >
        <TextareaAutosize
          value={input}
          onChange={handleChange}
          placeholder="อยากบ่นอะไร... พิมพ์มาเลย 🔥"
          minRows={3}
          maxRows={8}
          className="w-full text-2xl md:text-3xl bg-transparent resize-none border-none outline-none text-zinc-100 placeholder-zinc-600 font-dog leading-[1.6] pt-2"
        />
      </div>

      {/* Error */}
      {errorMsg && (
        <p className="relative mb-4 text-red-400 text-base font-medium text-center">
          {errorMsg}
        </p>
      )}

      {/* Submit */}
      <button
        type="button"
        onClick={handleTranslate}
        disabled={loading || credit === 0}
        className="relative px-12 py-3 bg-red-600 hover:bg-red-500 text-white text-2xl font-bold rounded-xl cursor-pointer shadow-lg transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
        style={{ boxShadow: "0 0 20px rgba(239,68,68,0.4)" }}
      >
        {loading
          ? "กำลังแปลง... 🔄"
          : level >= 4
          ? "ปล่อยของ 💀"
          : "เสร็จแล้ว 💢"}
      </button>
    </div>
  );
}
