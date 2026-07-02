// apps/web/app/page.tsx
"use client";

import { useState, useRef } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import TextareaAutosize from "react-textarea-autosize";
import PillBar from "../components/PillBar";
import { TARGETS, LEVELS } from "./lib/translateOptions";
import { useGuest } from "./context/GuestContext";
import { useAuth } from "./context/AuthContext";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export default function Home() {
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [target, setTarget] = useState<"boss" | "client" | "friend">("boss");
  const [level, setLevel] = useState(3);
  const [lang, setLang] = useState<"th" | "en">("th");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const router = useRouter();
  const { credit, refreshCredit } = useGuest();
  const { token } = useAuth();

  const activeLevel = LEVELS.find((l) => l.value === level)!;

  const handleTranslate = async () => {
    const text = textareaRef.current?.value ?? "";
    if (!text.trim()) {
      setErrorMsg("พิมพ์เรื่องที่อยากบ่นก่อนน้า~");
      return;
    }
    setErrorMsg(null);
    setLoading(true);

    try {
      const headers: Record<string, string> = {
        "Content-Type": "application/json",
      };
      if (token) headers["Authorization"] = `Bearer ${token}`;

      const res = await fetch(`${API_BASE}/translate`, {
        method: "POST",
        credentials: "include",
        headers,
        body: JSON.stringify({ text, level, target, lang }),
      });

      if (res.status === 402) {
        setErrorMsg("เครดิตหมดแล้ว กรุณาเข้าสู่ระบบเพื่อรับเครดิตเพิ่ม");
        return;
      }
      if (res.status === 401) {
        await refreshCredit();
        setErrorMsg("เซสชันหมดอายุ กรุณาลองใหม่");
        return;
      }
      if (!res.ok) {
        setErrorMsg("เกิดข้อผิดพลาด กรุณาลองใหม่");
        return;
      }

      const data = await res.json();
      await refreshCredit();

      sessionStorage.setItem("translationInput", text);
      sessionStorage.setItem("translationResult", JSON.stringify(data));
      sessionStorage.setItem("translationLevel", String(level));
      sessionStorage.setItem("translationTarget", target);
      sessionStorage.setItem("translationLang", lang);
      router.push("/second");
    } catch {
      setErrorMsg("เชื่อมต่อเซิร์ฟเวอร์ไม่ได้ กรุณาลองใหม่");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex-1 w-full bg-white flex flex-col items-center px-4 py-8">
      <div className="w-full max-w-5xl flex flex-col gap-5">
        {/* ── ส่งให้ใคร ───────────────────────────────── */}
        <div className="mx-auto w-full max-w-xl flex flex-col gap-2">
          <span className="text-[#ff7b00] text-xl font-bold pl-2">ส่งให้ใคร</span>
          <PillBar options={TARGETS} value={target} onChange={setTarget} />
        </div>

        {/* ── เลือกระดับภาษา ──────────────────────────── */}
        <div className="flex flex-col gap-2">
          <span className="text-[#ff7b00] text-xl font-bold pl-2">
            เลือกระดับภาษา
          </span>
          <PillBar
            options={LEVELS}
            value={level}
            onChange={setLevel}
            textSize="text-xs sm:text-base"
          />
        </div>

        {/* ── ฟีลลิ่ง (description, cross-fades on level switch) ── */}
        <div
          key={level}
          className="soft-fade min-h-[5rem] text-[15px] sm:text-base leading-relaxed text-black/90 px-2"
        >
          <p>
            <span className="text-[#64579f] font-bold">ฟีลลิ่ง</span>
            <span className="text-[#01a021] font-bold">: </span>
            {activeLevel.feeling}
          </p>
          <p className="mt-1">
            <span className="text-[#64579f] font-bold">ตัวอย่างคำบนเว็บ: </span>
            {activeLevel.example}
          </p>
        </div>

        {/* ── พิมพ์หม่องนี้เลยยยย (input) ──────────────── */}
        <div className="flex flex-col gap-2">
          <span className="text-[#ff7b00] text-xl font-bold pl-2">
            พิมพ์ตรงนี้เลยยยย
          </span>
          <div className="relative w-full rounded-[18px] border border-black/80 bg-white p-5 pb-20 shadow-sm">
            <TextareaAutosize
              ref={textareaRef as React.RefObject<HTMLTextAreaElement>}
              defaultValue=""
              placeholder="พิมพ์เรื่องที่อยากบ่น... แล้วเราจะแปลงให้สุภาพ"
              minRows={5}
              maxRows={12}
              className="w-full resize-none border-none bg-transparent text-xl md:text-2xl leading-[1.6] text-black outline-none placeholder-black/30 font-dog"
            />

            {/* Send button (orange ">") */}
            <button
              type="button"
              onClick={handleTranslate}
              disabled={loading || credit === 0}
              aria-label="แปลงข้อความ"
              className="absolute bottom-4 right-4 flex h-[60px] w-[60px] items-center justify-center rounded-[22px] bg-[#ff8055] text-white text-2xl font-bold shadow-md transition-all hover:bg-[#ff6a3c] active:scale-95 disabled:opacity-40 disabled:cursor-not-allowed"
            >
              {loading ? (
                <span className="animate-spin text-xl">↻</span>
              ) : (
                ">"
              )}
            </button>
          </div>
        </div>

        {/* Error */}
        {errorMsg && (
          <p className="text-center text-[#ff7b00] font-bold">{errorMsg}</p>
        )}

        {/* Out of credit — give the user a next step */}
        {credit === 0 &&
          (token ? (
            <div className="text-center">
              <Link
                href="/payment"
                className="inline-block rounded-full bg-[#ff8055] px-6 py-2 text-white font-bold hover:bg-[#ff6a3c] transition-colors"
              >
                เติมเครดิต
              </Link>
            </div>
          ) : (
            <p className="text-center text-[#ff7b00] font-bold">
              เครดิตหมดแล้ว กรุณาเข้าสู่ระบบเพื่อรับเครดิตเพิ่ม
            </p>
          ))}
      </div>
    </div>
  );
}
