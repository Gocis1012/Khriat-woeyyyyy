// apps/web/app/page.tsx
"use client";

import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import TextareaAutosize from "react-textarea-autosize";
import { useGuest } from "./context/GuestContext";
import { useAuth } from "./context/AuthContext";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

// ── Floating particles ────────────────────────────────────────────────────────
const PARTICLES = [
  { emoji: "💢", left: "4%",  delay: "0s",   duration: "7s"   },
  { emoji: "🔥", left: "14%", delay: "1.5s", duration: "8s"   },
  { emoji: "😤", left: "24%", delay: "0.3s", duration: "9s"   },
  { emoji: "💥", left: "37%", delay: "2.4s", duration: "7.5s" },
  { emoji: "😡", left: "51%", delay: "0.9s", duration: "10s"  },
  { emoji: "🤬", left: "64%", delay: "1.8s", duration: "6.5s" },
  { emoji: "🔥", left: "76%", delay: "3.2s", duration: "8.5s" },
  { emoji: "💢", left: "87%", delay: "0.5s", duration: "7s"   },
  { emoji: "👊", left: "95%", delay: "2s",   duration: "9.5s" },
];

// ── Target contexts ───────────────────────────────────────────────────────────
const TARGETS = [
  {
    value: "boss",
    emoji: "👔",
    label: "เจ้านาย",
    desc: "อำนาจสูงกว่า กระทบอนาคต",
    placeholder: "เจ้านายสั่งงานตอนศุกร์บ่ายโมงอีกแล้ว...",
    taglines: [
      "เจ้านายสั่งงานตอนศุกร์บ่ายโมง 🙃",
      'โดนถามว่า "เสร็จยัง?" ตอน 5 ทุ่ม 🫠',
      '"แก้นิดเดียว" แต่ต้องทำใหม่ทั้งหมด 🫃',
      "Meeting ที่ควรเป็น Email 📧",
    ],
  },
  {
    value: "client",
    emoji: "💼",
    label: "ลูกค้า",
    desc: "แหล่งรายได้ รักษาภาพลักษณ์",
    placeholder: "ลูกค้าบอกว่าแก้นิดเดียว แต่...",
    taglines: [
      'ลูกค้าบอกว่า "แก้นิดเดียว" 😇',
      "เปลี่ยน Brief ครั้งที่ 7 ในอาทิตย์เดียว 💀",
      '"ทำให้มันป็อปขึ้นหน่อย" 🤡',
      "Deadline พรุ่งนี้เช้า แต่เพิ่ง Confirm วันนี้เย็น 🚨",
    ],
  },
  {
    value: "teacher",
    emoji: "📚",
    label: "อาจารย์",
    desc: "ผู้ให้ความรู้ มีอำนาจเกรด",
    placeholder: "อาจารย์ให้งานด่วนอีกแล้ว...",
    taglines: [
      "อาจารย์บอกส่งพรุ่งนี้ ตอน 4 โมงเย็น 📚",
      "เกรดออกมาแล้วแต่ไม่อธิบายเหตุผล 😤",
      '"ดูจากตัวอย่างในหนังสือ" แต่ไม่มีตัวอย่างจริง 🙃',
      "Assignment 5 ข้อ ภายใน 2 วัน 🫠",
    ],
  },
  {
    value: "friend",
    emoji: "🫂",
    label: "เพื่อน",
    desc: "เท่าเทียมกัน แต่ต้องสื่อให้ชัด",
    placeholder: "เพื่อนทำอะไรบางอย่างที่หงุดหงิดมาก...",
    taglines: [
      "เพื่อนบอกจะมา แล้วก็หาย 🫠",
      'เปิดแชทเห็น "พิมพ์แล้วลบ" 😰',
      "ขอยืมเงินแล้วหายไป 3 อาทิตย์ 💸",
      "บอก Split bill แต่ดันสั่งอาหารแพงสุด 🤌",
    ],
  },
];

// ── Level definitions ─────────────────────────────────────────────────────────
const LEVELS = [
  { value: 1, emoji: "🙏", label: "สุภาพสุดๆ", glowColor: "rgba(34,197,94,0.45)",  btnClass: "from-green-600  to-green-700"  },
  { value: 2, emoji: "😊", label: "สุภาพ",      glowColor: "rgba(59,130,246,0.45)", btnClass: "from-blue-600   to-blue-700"   },
  { value: 3, emoji: "😐", label: "ปกติ",       glowColor: "rgba(234,179,8,0.45)",  btnClass: "from-yellow-600 to-yellow-700" },
  { value: 4, emoji: "😈", label: "แรงนิดๆ",    glowColor: "rgba(249,115,22,0.45)", btnClass: "from-orange-600 to-orange-700" },
  { value: 5, emoji: "💀", label: "แรงสุด",     glowColor: "rgba(239,68,68,0.55)",  btnClass: "from-red-600    to-red-700"    },
];

export default function Home() {
  const [input, setInput]       = useState("");
  const [targetIdx, setTargetIdx] = useState(0);       // index into TARGETS
  const [level, setLevel]       = useState(3);
  const [loading, setLoading]   = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [shaking, setShaking]   = useState(false);
  const [taglineIdx, setTaglineIdx] = useState(0);
  const boxRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const { credit, refreshCredit } = useGuest();
  const { token } = useAuth();

  const target = TARGETS[targetIdx];
  const currentGlow = LEVELS[level - 1].glowColor;

  // Rotate taglines per target
  useEffect(() => {
    setTaglineIdx(0);
    const iv = setInterval(() => {
      setTaglineIdx((i) => (i + 1) % target.taglines.length);
    }, 3500);
    return () => clearInterval(iv);
  }, [targetIdx, target.taglines.length]);

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
      const headers: Record<string, string> = { "Content-Type": "application/json" };
      if (token) headers["Authorization"] = `Bearer ${token}`;

      const res = await fetch(`${API_BASE}/translate`, {
        method: "POST",
        credentials: "include",
        headers,
        body: JSON.stringify({ text: input, level, target: target.value }),
      });

      if (res.status === 402) { setErrorMsg("เครดิตหมดแล้ว กรุณาเข้าสู่ระบบเพื่อรับเครดิตเพิ่ม"); return; }
      if (res.status === 401) { await refreshCredit(); setErrorMsg("เซสชันหมดอายุ กรุณาลองใหม่"); return; }
      if (!res.ok)            { setErrorMsg("เกิดข้อผิดพลาด กรุณาลองใหม่"); return; }

      const data = await res.json();
      await refreshCredit();

      sessionStorage.setItem("translationInput",  input);
      sessionStorage.setItem("translationResult", JSON.stringify(data));
      sessionStorage.setItem("translationLevel",  String(level));
      sessionStorage.setItem("translationTarget", target.value);
      router.push("/second");
    } catch {
      setErrorMsg("ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้ กรุณาลองใหม่");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative flex-1 flex flex-col items-center justify-center w-full px-4 py-10 overflow-hidden bg-[#111111]">

      {/* Radial glow — shifts with level */}
      <div
        className="absolute inset-0 pointer-events-none transition-all duration-700"
        style={{
          background: `radial-gradient(ellipse 70% 50% at 50% 60%, ${currentGlow} 0%, transparent 70%)`,
          animation: "rage-pulse 3s ease-in-out infinite",
        }}
      />

      {/* Particles */}
      {PARTICLES.map((p, i) => (
        <span key={i} className="particle" style={{ left: p.left, animationDelay: p.delay, animationDuration: p.duration }}>
          {p.emoji}
        </span>
      ))}

      {/* Title */}
      <h1
        className="relative text-4xl md:text-5xl font-bold text-white mb-1 text-center tracking-wide"
        style={{ textShadow: "0 0 32px rgba(239,68,68,0.6)" }}
      >
        พิมพ์เรื่องที่อยากบ่นให้เราฟังสิ
      </h1>

      {/* Rotating tagline (target-specific) */}
      <p key={`${targetIdx}-${taglineIdx}`} className="relative text-sm text-zinc-400 mb-5 text-center fade-up">
        {target.taglines[taglineIdx]}
      </p>

      {/* ── STEP 1: Target selector ──────────────────────────────── */}
      <div className="relative w-full max-w-3xl mb-4">
        <p className="text-xs text-zinc-500 mb-2 text-center tracking-widest uppercase">
          ส่งให้ใคร?
        </p>
        <div className="grid grid-cols-4 gap-2">
          {TARGETS.map((t, i) => (
            <button
              key={t.value}
              onClick={() => setTargetIdx(i)}
              className={`
                flex flex-col items-center gap-1 py-3 px-2 rounded-2xl border transition-all text-center
                ${targetIdx === i
                  ? "bg-zinc-700 border-white/30 scale-105 shadow-lg"
                  : "bg-zinc-900/70 border-zinc-700 hover:bg-zinc-800"}
              `}
            >
              <span className="text-2xl">{t.emoji}</span>
              <span className={`text-sm font-bold ${targetIdx === i ? "text-white" : "text-zinc-300"}`}>
                {t.label}
              </span>
              <span className="text-[10px] text-zinc-500 leading-tight">{t.desc}</span>
            </button>
          ))}
        </div>
      </div>

      {/* ── STEP 2: Level (Rage Meter) ───────────────────────────── */}
      <div className="relative w-full max-w-3xl mb-4">
        <p className="text-xs text-zinc-500 mb-2 text-center tracking-widest uppercase">
          ระดับความแรง
        </p>
        <div className="flex gap-2 justify-center">
          {LEVELS.map((l) => (
            <button
              key={l.value}
              onClick={() => setLevel(l.value)}
              className={`
                flex flex-col items-center px-4 py-2 rounded-xl transition-all text-xs font-medium border flex-1 max-w-[80px]
                ${level === l.value
                  ? `bg-gradient-to-b ${l.btnClass} text-white border-white/20 scale-110 shadow-lg`
                  : "bg-zinc-800/70 text-zinc-400 border-zinc-700 hover:bg-zinc-700"}
              `}
              style={level === l.value ? { boxShadow: `0 0 16px ${l.glowColor}` } : {}}
            >
              <span className="text-lg mb-0.5">{l.emoji}</span>
              <span className="leading-tight text-center">{l.label}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Credit */}
      <div className="relative mb-3 flex items-center gap-2 text-sm">
        {credit === null ? (
          <span className="text-zinc-600">กำลังโหลด...</span>
        ) : credit <= 0 ? (
          <span className="px-3 py-1 rounded-full bg-red-900/60 border border-red-600 text-red-300 font-semibold animate-pulse">
            ⚠️ เครดิตหมด — กรุณาเข้าสู่ระบบ
          </span>
        ) : (
          <span className="px-3 py-1 rounded-full bg-zinc-800 border border-zinc-600 text-zinc-300">
            🔥 เครดิตคงเหลือ <span className="font-bold text-white">{credit}</span> ครั้ง
          </span>
        )}
      </div>

      {/* ── STEP 3: Input Box ────────────────────────────────────── */}
      <div
        ref={boxRef}
        className={`relative w-full max-w-3xl rounded-3xl border bg-zinc-900/80 backdrop-blur-sm p-6 mb-5 transition-all duration-500 ${shaking ? "shake" : ""}`}
        style={{
          borderColor: currentGlow,
          boxShadow: `0 0 20px 3px ${currentGlow}`,
          animation: "box-glow 2.4s ease-in-out infinite",
        }}
      >
        <TextareaAutosize
          value={input}
          onChange={handleChange}
          placeholder={target.placeholder}
          minRows={3}
          maxRows={8}
          className="w-full text-2xl md:text-3xl bg-transparent resize-none border-none outline-none text-zinc-100 placeholder-zinc-600 font-dog leading-[1.6] pt-2"
        />
      </div>

      {errorMsg && (
        <p className="relative mb-3 text-red-400 text-base font-medium text-center">{errorMsg}</p>
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
          ? `ปล่อยของใส่ ${target.emoji} 💀`
          : `แปลงให้ ${target.emoji} 💢`}
      </button>
    </div>
  );
}
