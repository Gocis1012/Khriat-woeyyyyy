// apps/web/app/page.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import TextareaAutosize from "react-textarea-autosize";
import { useGuest } from "./context/GuestContext";

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export default function Home() {
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const router = useRouter();
  const { refreshCredit } = useGuest();

  const handleTranslate = async () => {
    if (!input.trim()) {
      setErrorMsg("กรุณากรอกข้อความก่อน");
      return;
    }

    setErrorMsg(null);
    setLoading(true);

    try {
      const response = await fetch(`${API_BASE}/translate`, {
        method: "POST",
        credentials: "include", // send guest_id cookie cross-origin
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: input }),
      });

      if (response.status === 402) {
        // Out of credits
        setErrorMsg("เครดิตหมดแล้ว กรุณาเข้าสู่ระบบเพื่อรับเครดิตเพิ่ม");
        return;
      }

      if (response.status === 401) {
        // Session expired — refresh credit will create a new session
        await refreshCredit();
        setErrorMsg("เซสชันหมดอายุ กรุณาลองใหม่");
        return;
      }

      if (!response.ok) {
        setErrorMsg("เกิดข้อผิดพลาด กรุณาลองใหม่");
        return;
      }

      const data = await response.json();

      // Refresh credit count in Navbar after successful translation
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
    <div className="flex-1 flex flex-col items-center justify-center w-full px-4 py-12">
      <h1 className="text-4xl md:text-5xl font-bold text-[#2D2D2D] mb-6 text-center tracking-wide">
        พิมพ์เรื่องที่อยากบ่นให้เราฟังสิ
      </h1>

      <div className="w-full max-w-3xl bg-white rounded-3xl border border-slate-300 shadow-sm p-6 mb-6">
        <TextareaAutosize
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="พิมพ์ในนี้นะ...."
          minRows={3}
          maxRows={8}
          className="w-full text-2xl md:text-3xl bg-transparent resize-none border-none outline-none text-slate-800 placeholder-slate-400 font-dog leading-[1.6] pt-2"
        />
      </div>

      {/* Error message */}
      {errorMsg && (
        <p className="mb-4 text-red-500 text-base font-medium text-center">
          {errorMsg}
        </p>
      )}

      <button
        type="button"
        onClick={handleTranslate}
        disabled={loading}
        className="px-12 py-2.5 bg-[#FCFF91] hover:bg-[#F4F776] border border-slate-400 text-2xl font-bold rounded-xl cursor-pointer shadow-sm transition-all active:scale-95 text-slate-800 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? "กำลังประมวลผล..." : "เสร็จแล้ว"}
      </button>

      <Link
        href="/second"
        className="text-sm text-slate-400 hover:underline mt-8"
      >
        ดูหน้าสอง (Second Page) →
      </Link>
    </div>
  );
}
