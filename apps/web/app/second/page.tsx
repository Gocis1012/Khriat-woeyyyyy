"use client";

import { useEffect, useState } from "react";
import Link from "next/link";

export default function Second() {
  const [result, setResult] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [copied, setCopied] = useState(false);

  useEffect(() => {
    const stored = sessionStorage.getItem("translationResult");
    const input = sessionStorage.getItem("translationInput");
    if (stored) {
      setResult(JSON.parse(stored));
    }
    if (input) {
      setResult((prev: any) => ({ ...prev, originalInput: input }));
    }
    setLoading(false);
  }, []);

  const handleCopy = async () => {
    const textToCopy = result.result || "ไม่มีข้อมูล";
    try {
      await navigator.clipboard.writeText(textToCopy);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000); // Reset after 2 seconds
    } catch (err) {
      console.error("Failed to copy:", err);
    }
  };

  return (
    <div className="flex flex-col flex-1 items-center justify-center bg-zinc-50">
      <main className="flex flex-1 w-full max-w-3xl flex-col items-center justify-between py-10 px-16 sm:items-start text-black">
        <div className="w-full">
          <h1 className="max-w-xs text-3xl font-semibold leading-10 tracking-tight text-black mb-8">
            ผลลัพท์การแปล
          </h1>

          {loading ? (
            <p className="text-lg text-slate-600">กำลังโหลด...</p>
          ) : result ? (
            <div className="w-full space-y-6">
              <div className="p-6 bg-blue-50 dark:bg-blue-950 rounded-lg border border-blue-200 dark:border-blue-800">
                <h2 className="text-sm font-semibold text-slate-600 dark:text-slate-400 mb-2">
                  ข้อความต้นฉบับ:
                </h2>
                <p className="text-lg text-slate-800 dark:text-slate-100">
                  {result.originalInput || result.text || result.original || "ไม่มีข้อมูล"}
                </p>
              </div>

              <div className="p-6 bg-green-50 rounded-lg border border-black-200 text-black-800">
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-sm font-semibold text-slate-600 text-black-800">
                    ผลลัพท์:
                  </h2>
                  <button
                    onClick={handleCopy}
                    className={`px-3 py-1 rounded text-sm font-medium transition-all ${
                      copied
                        ? "bg-green-500 text-white"
                        : "bg-slate-300 hover:bg-slate-400 text-slate-800"
                    }`}
                  >
                    {copied ? "✓ Copied!" : "Copy"}
                  </button>
                </div>
                <p className="text-lg text-black-800 dark:text-black-100">
                  {result.result || "ไม่มีข้อมูล"}
                </p>
              </div>

      
            </div>
          ) : (
            <p className="text-lg text-black-600">ไม่พบผลลัพท์</p>
          )}

          <Link
            href="/"
            className="inline-block mt-8 px-6 py-2 bg-slate-200 hover:bg-slate-300 text-slate-800 font-semibold rounded-lg transition-colors"
          >
            ← กลับไปแปลใหม่
          </Link>
        </div>
      </main>
    </div>
  );
}