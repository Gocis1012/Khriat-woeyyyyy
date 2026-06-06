"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { targetLabel, levelLabel } from "../lib/translateOptions";

export default function Second() {
  const [result, setResult] = useState<string | null>(null);
  const [original, setOriginal] = useState<string>("");
  const [level, setLevel] = useState<number>(3);
  const [target, setTarget] = useState<string>("boss");
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);
  const router = useRouter();

  useEffect(() => {
    const stored = sessionStorage.getItem("translationResult");
    const input = sessionStorage.getItem("translationInput");
    const lvl = sessionStorage.getItem("translationLevel");
    const tgt = sessionStorage.getItem("translationTarget");
    if (stored) setResult(JSON.parse(stored).result ?? null);
    if (input) setOriginal(input);
    if (lvl) setLevel(parseInt(lvl, 10));
    if (tgt) setTarget(tgt);
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
    <div className="flex-1 w-full bg-white flex flex-col items-center px-4 py-8">
      <div className="w-full max-w-5xl flex flex-col gap-5">
        {loading ? (
          <p className="text-black/50 text-lg">กำลังโหลด...</p>
        ) : result ? (
          <div className="flex flex-col gap-5 fade-up">
            {/* ส่งให้ใคร */}
            <div className="flex flex-col gap-1">
              <span className="text-[#ff7b00] text-xl font-bold pl-2">
                ส่งให้ใคร
              </span>
              <span className="text-black/80 text-xl pl-3">
                {targetLabel(target)}
              </span>
            </div>

            {/* เลือกระดับภาษา */}
            <div className="flex flex-col gap-1">
              <span className="text-[#ff7b00] text-xl font-bold pl-2">
                เลือกระดับภาษา
              </span>
              <span className="text-black/80 text-xl pl-3">
                {levelLabel(level)}
              </span>
            </div>

            {/* ของเจ้า (original) */}
            <div className="flex flex-col gap-1">
              <span className="text-[#ff7b00] text-xl font-bold pl-2">
                ของเจ้า
              </span>
              <p className="text-black/80 text-lg pl-3 whitespace-pre-wrap break-words">
                {original || "—"}
              </p>
            </div>

            {/* ของข่อย (result box) */}
            <div className="flex flex-col gap-2">
              <span className="text-[#ff7b00] text-xl font-bold pl-2">
                ของข่อย
              </span>
              <div className="relative w-full rounded-[18px] border border-black/80 bg-white p-6 pb-20 shadow-sm">
                <p className="text-xl md:text-2xl leading-[1.7] text-black whitespace-pre-wrap break-words">
                  {result}
                </p>

                {/* Copy button */}
                <button
                  type="button"
                  onClick={handleCopy}
                  className={`absolute bottom-4 right-4 flex h-[60px] items-center justify-center rounded-[22px] px-5 text-lg font-bold text-white shadow-md transition-all active:scale-95 ${
                    copied ? "bg-[#01a021]" : "bg-[#ff8055] hover:bg-[#ff6a3c]"
                  }`}
                >
                  {copied ? "คัดลอกแล้ว ✓" : "คัดลอก"}
                </button>
              </div>
            </div>

            {/* Back */}
            <button
              onClick={() => router.push("/")}
              className="self-start rounded-full border border-black/30 px-6 py-2 text-base text-black/70 transition-colors hover:border-black/60 hover:text-black"
            >
              ← บ่นใหม่
            </button>
          </div>
        ) : (
          <div className="text-center fade-up py-20">
            <p className="text-black/60 text-lg mb-4">ไม่พบผลลัพท์</p>
            <button
              onClick={() => router.push("/")}
              className="rounded-full bg-[#ff8055] px-6 py-2 text-white font-bold hover:bg-[#ff6a3c] transition-colors"
            >
              ← กลับไปบ่น
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
