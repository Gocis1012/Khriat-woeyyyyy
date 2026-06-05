// apps/web/app/page.tsx
"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import TextareaAutosize from "react-textarea-autosize";
import { useGuest } from "./context/GuestContext";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

// Anger particles — emoji + position + speed
const PARTICLES = [
  { emoji: "💢", left: "8%",  delay: "0s",    duration: "6s"  },
  { emoji: "🔥", left: "18%", delay: "1.2s",  duration: "7s"  },
  { emoji: "😤", left: "30%", delay: "0.4s",  duration: "8.5s"},
  { emoji: "💥", left: "45%", delay: "2.1s",  duration: "6.5s"},
  { emoji: "😡", left: "58%", delay: "0.8s",  duration: "9s"  },
  { emoji: "🔥", left: "70%", delay: "1.7s",  duration: "7.5s"},
  { emoji: "💢", left: "82%", delay: "3s",    duration: "6s"  },
  { emoji: "💥", left: "92%", delay: "0.2s",  duration: "8s"  },
];

export default function Home() {
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const [shaking, setShaking] = useState(false);
  const boxRef = useRef<HTMLDivElement>(null);
  const router = useRouter();
  const { credit, refreshCredit } = useGuest();

  const triggerShake = () => {
    setShaking(true);
    setTimeout(() => setShaking(false), 120);
  };

  const handleChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInput(e.target.value);
    // subtle shake every ~5 chars typed
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
      const response = await fetch(`${API_BASE}/translate`, {
        method: "POST",
        credentials: "include",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ text: input }),
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
      router.push("/second");
    } catch {
      setErrorMsg("ไม่สามารถเชื่อมต่อกับเซิร์ฟเวอร์ได้ กรุณาลองใหม่");
    } finally {
      setLoading(false);
    }
  };

  return (
    /* Dark angry background */
    <div className="relative flex-1 flex flex-col items-center justify-center w-full px-4 py-12 overflow-hidden bg-[#111111]">

      {/* Radial red glow behind content */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background: "radial-gradient(ellipse 70% 50% at 50% 60%, rgba(185,28,28,0.22) 0%, transparent 70%)",
          animation: "rage-pulse 3s ease-in-out infinite",
        }}
      />

      {/* Floating anger particles */}
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
        className="relative text-4xl md:text-5xl font-bold text-white mb-3 text-center tracking-wide"
        style={{ textShadow: "0 0 32px rgba(239,68,68,0.6)" }}
      >
        พิมพ์เรื่องที่อยากบ่นให้เราฟังสิ
      </h1>

      {/* Credit indicator */}
      <div className="relative mb-6 flex items-center gap-2 text-sm">
        {credit === null ? (
          <span className="text-zinc-500">กำลังโหลด...</span>
        ) : credit <= 0 ? (
          <span className="px-3 py-1 rounded-full bg-red-900/60 border border-red-600 text-red-300 font-semibold">
            ⚠️ เครดิตหมด — กรุณาเข้าสู่ระบบ
          </span>
        ) : (
          <span className="px-3 py-1 rounded-full bg-zinc-800 border border-zinc-600 text-zinc-300">
            🔥 เครดิตคงเหลือ{" "}
            <span className="font-bold text-white">{credit}</span> ครั้ง
          </span>
        )}
      </div>

      {/* Input box */}
      <div
        ref={boxRef}
        className={`relative w-full max-w-3xl rounded-3xl border border-red-800/60 bg-zinc-900/80 backdrop-blur-sm p-6 mb-6 rage-box ${shaking ? "shake" : ""}`}
      >
        <TextareaAutosize
          value={input}
          onChange={handleChange}
          placeholder="พิมพ์ในนี้นะ...."
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

      {/* Submit button */}
      <button
        type="button"
        onClick={handleTranslate}
        disabled={loading || credit === 0}
        className="relative px-12 py-3 bg-red-600 hover:bg-red-500 text-white text-2xl font-bold rounded-xl cursor-pointer shadow-lg transition-all active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
        style={{ boxShadow: "0 0 20px rgba(239,68,68,0.4)" }}
      >
        {loading ? "กำลังประมวลผล..." : "เสร็จแล้ว 💢"}
      </button>
    </div>
  );
}
