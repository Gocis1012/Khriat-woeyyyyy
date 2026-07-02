"use client";

import { useState, useRef, useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "../context/AuthContext";
import { useGuest } from "../context/GuestContext";

const API_BASE =
  process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

const MIN_AMOUNT = 20;
const MAX_AMOUNT = 5000;

type PaymentStatus = "idle" | "creating" | "pending" | "success" | "failed" | "error";

export default function PaymentPage() {
  const router = useRouter();
  const { token, isLoggedIn, loading: authLoading } = useAuth();
  const { refreshCredit } = useGuest();

  const [amount, setAmount] = useState(50);
  const [status, setStatus] = useState<PaymentStatus>("idle");
  const [qrCodeUri, setQrCodeUri] = useState<string | null>(null);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);
  const pollRef = useRef<ReturnType<typeof setInterval> | null>(null);

  // Auth guard — guests can't enter the payment flow.
  useEffect(() => {
    if (!authLoading && !isLoggedIn) {
      router.replace("/");
    }
  }, [authLoading, isLoggedIn, router]);

  const stopPolling = useCallback(() => {
    if (pollRef.current) {
      clearInterval(pollRef.current);
      pollRef.current = null;
    }
  }, []);

  useEffect(() => stopPolling, [stopPolling]);

  const pollStatus = useCallback(
    (paymentId: string) => {
      stopPolling();
      pollRef.current = setInterval(async () => {
        try {
          const res = await fetch(`${API_BASE}/api/v1/payments/${paymentId}/status`, {
            headers: { Authorization: `Bearer ${token}` },
            credentials: "include",
          });
          if (!res.ok) return;
          const data = await res.json();
          if (data.status === "success") {
            setStatus("success");
            stopPolling();
            await refreshCredit();
          } else if (data.status === "failed") {
            setStatus("failed");
            stopPolling();
          }
        } catch {
          // transient network error — keep polling until the next tick
        }
      }, 3000);
    },
    [token, stopPolling, refreshCredit],
  );

  const handleCreateCharge = async () => {
    if (!token) return;
    setErrorMsg(null);
    setStatus("creating");

    try {
      const res = await fetch(`${API_BASE}/api/v1/payments/create`, {
        method: "POST",
        credentials: "include",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ amount }),
      });

      if (!res.ok) {
        setStatus("error");
        setErrorMsg("สร้างรายการชำระเงินไม่สำเร็จ กรุณาลองใหม่");
        return;
      }

      const data = await res.json();
      setQrCodeUri(data.qrCodeUri);
      setStatus("pending");
      pollStatus(data.paymentId);
    } catch {
      setStatus("error");
      setErrorMsg("เชื่อมต่อเซิร์ฟเวอร์ไม่ได้ กรุณาลองใหม่");
    }
  };

  const handleRetry = () => {
    stopPolling();
    setStatus("idle");
    setQrCodeUri(null);
    setErrorMsg(null);
  };

  if (authLoading || !isLoggedIn) {
    return (
      <div className="flex-1 w-full flex items-center justify-center py-20">
        <p className="text-black/50 text-lg">กำลังตรวจสอบสิทธิ์...</p>
      </div>
    );
  }

  return (
    <div className="flex-1 w-full bg-white flex flex-col items-center px-4 py-8">
      <div className="w-full max-w-md flex flex-col gap-6">
        <h1 className="text-2xl font-bold text-[#ff7b00]">เติมเครดิต</h1>

        {status === "success" ? (
          <div className="flex flex-col items-center gap-4 fade-up py-10">
            <div className="text-5xl">✅</div>
            <p className="text-xl font-bold text-[#01a021]">ชำระเงินสำเร็จ!</p>
            <button
              onClick={() => router.push("/")}
              className="rounded-full bg-[#ff8055] px-6 py-2 text-white font-bold hover:bg-[#ff6a3c] transition-colors"
            >
              กลับไปบ่นต่อ
            </button>
          </div>
        ) : status === "failed" ? (
          <div className="flex flex-col items-center gap-4 fade-up py-10">
            <div className="text-5xl">❌</div>
            <p className="text-xl font-bold text-red-500">การชำระเงินล้มเหลว</p>
            <button
              onClick={handleRetry}
              className="rounded-full border border-black/30 px-6 py-2 text-black/70 hover:border-black/60 hover:text-black transition-colors"
            >
              ลองอีกครั้ง
            </button>
          </div>
        ) : status === "pending" && qrCodeUri ? (
          <div className="flex flex-col items-center gap-4 fade-up">
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={qrCodeUri}
              alt="QR Code สำหรับชำระเงิน"
              className="w-64 h-64 border border-black/10 rounded-2xl"
            />
            <p className="text-black/70">สแกน QR เพื่อชำระเงิน {amount} บาท</p>
            <p className="text-black/40 text-sm animate-pulse">
              กำลังรอการชำระเงิน...
            </p>
            <button
              onClick={handleRetry}
              className="text-sm text-black/40 hover:text-black/70 transition-colors"
            >
              ยกเลิก
            </button>
          </div>
        ) : (
          <div className="flex flex-col gap-4">
            <label className="flex flex-col gap-2">
              <span className="text-black/70">
                จำนวนเงิน (บาท) — {MIN_AMOUNT}-{MAX_AMOUNT}
              </span>
              <input
                type="number"
                min={MIN_AMOUNT}
                max={MAX_AMOUNT}
                step={10}
                value={amount}
                onChange={(e) => setAmount(Number(e.target.value))}
                className="rounded-xl border border-black/20 px-4 py-2 text-lg outline-none focus:border-[#ff8055]"
              />
            </label>
            <button
              onClick={handleCreateCharge}
              disabled={
                status === "creating" || amount < MIN_AMOUNT || amount > MAX_AMOUNT
              }
              className="rounded-full bg-[#ff8055] px-6 py-3 text-white font-bold hover:bg-[#ff6a3c] transition-colors disabled:opacity-40 disabled:cursor-not-allowed"
            >
              {status === "creating" ? "กำลังสร้าง QR..." : "สร้าง QR Code"}
            </button>
            {errorMsg && (
              <p className="text-center text-red-500 font-bold">{errorMsg}</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
