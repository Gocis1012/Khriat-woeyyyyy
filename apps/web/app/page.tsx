// apps/web/app/page.tsx
"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import TextareaAutosize from "react-textarea-autosize";

export default function Home() {
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleTranslate = async () => {
    if (!input.trim()) {
      alert("กรุณากรอกข้อความก่อน");
      return;
    }

    setLoading(true);
    try {
      const response = await fetch("http://localhost:8080/translate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ text: input }),
      });

      if (!response.ok) {
        throw new Error("API call failed");
      }

      const data = await response.json();
      // Store both original input and API result
      sessionStorage.setItem("translationInput", input);
      sessionStorage.setItem("translationResult", JSON.stringify(data));
      router.push("/second");
    } catch (error) {
      console.error("Error:", error);
      alert("เกิดข้อผิดพลาด กรุณาลองใหม่");
    } finally {
      setLoading(false);
    }
  };

  return (
    // เปลี่ยนมาใช้ py-12 เพื่อให้มีระยะห่างจากขอบบน (Navbar) และขอบล่างที่สมดุลกันพอดี
    <div className="flex-1 flex flex-col items-center justify-center w-full px-4 py-12">
      {/* 1. หัวข้อด้านบน ปรับลด margin-bottom (mb-6) ให้กระชับขึ้น */}
      <h1 className="text-4xl md:text-5xl font-bold text-[#2D2D2D] mb-6 text-center tracking-wide">
        พิมพ์เรื่องที่อยากบ่นให้เราฟังสิ
      </h1>

      {/* ⭐️ 2. แก้ไขกล่องหุ้มภายนอก (ลบ overflow-hidden ออก เพื่อให้กล่องงอกลงมาได้ และเอา p-6 ออกเพื่อความเนียน) */}
      <div className="w-full max-w-3xl bg-white rounded-3xl border border-slate-300 shadow-sm p-6 mb-6">
        <TextareaAutosize
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="พิมพ์ในนี้นะ...."
          // ⭐️ 3. ตั้งค่าการยืดหยุ่นความสูง
          minRows={3} // ความสูงเริ่มต้น (ประมาณ 3 บรรทัด)
          maxRows={8} // ขยายได้สูงสุด 8 บรรทัด พอเกินจากนี้จะเปิดระบบ Scroll mouse อัตโนมัติ
          /* ⭐️ 4. ปรับ Class Tailwind นิดหน่อย:
      - เอา min-h ออก (เพราะใช้ minRows คุมแทนแล้ว)
      - ใส่ pt-2 หรือ pt-4 ดันไม้จัตวาเหมือนเดิม
      - ใส่ custom scrollbar (ถ้าต้องการ) หรือปล่อยให้ระบบขึ้น scroll อัตโนมัติเมื่อเกิน maxRows
    */
          className="w-full text-2xl md:text-3xl bg-transparent resize-none border-none outline-none text-slate-800 placeholder-slate-400 font-dog leading-[1.6] pt-2"
        />
      </div>

      {/* 3. ปุ่ม "เสร็จแล้ว" สีเหลือง มนพอดีๆ */}
      <button
        type="button"
        onClick={handleTranslate}
        disabled={loading}
        className="px-12 py-2.5 bg-[#FCFF91] hover:bg-[#F4F776] border border-slate-400 text-2xl font-bold rounded-xl cursor-pointer shadow-sm transition-all active:scale-95 text-slate-800 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {loading ? "กำลังประมวลผล..." : "เสร็จแล้ว"}
      </button>

      {/* ลิงก์ไปหน้าสอง */}
      <Link
        href="/second"
        className="text-sm text-slate-400 hover:underline mt-8"
      >
        ดูหน้าสอง (Second Page) →
      </Link>
    </div>
  );
}
