"use client";

import { useEffect, useLayoutEffect, useRef, useState } from "react";

export interface PillOption {
  value: string | number;
  label: string;
}

interface PillBarProps {
  options: PillOption[];
  value: string | number;
  onChange: (value: never) => void;
  /** Tailwind text size for labels */
  textSize?: string;
}

// useLayoutEffect on client, useEffect on server (avoids SSR warning)
const useIsoLayoutEffect =
  typeof window !== "undefined" ? useLayoutEffect : useEffect;

/**
 * A rounded pill selector with a single orange highlight that
 * smoothly slides between options when the selection changes.
 */
export default function PillBar({
  options,
  value,
  onChange,
  textSize = "text-base sm:text-xl",
}: PillBarProps) {
  const btnRefs = useRef<Array<HTMLButtonElement | null>>([]);
  const [indicator, setIndicator] = useState({ left: 0, width: 0, ready: false });

  const recalc = () => {
    const idx = options.findIndex((o) => o.value === value);
    const btn = btnRefs.current[idx];
    if (btn) {
      setIndicator({ left: btn.offsetLeft, width: btn.offsetWidth, ready: true });
    }
  };

  // Measure before the browser paints so the highlight starts in the
  // correct place (no animate-from-zero flash on first render).
  useIsoLayoutEffect(() => {
    recalc();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value, options]);

  // Re-measure on resize (font load / responsive width changes)
  useEffect(() => {
    window.addEventListener("resize", recalc);
    return () => window.removeEventListener("resize", recalc);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [value, options]);

  return (
    <div className="relative flex items-stretch w-full rounded-[24px] border border-black/80 bg-white p-1.5 shadow-sm">
      {/* Sliding orange highlight */}
      <div
        className="pointer-events-none absolute top-1.5 bottom-1.5 rounded-[18px] bg-[#ff8055] transition-all duration-300 ease-out"
        style={{
          left: indicator.left,
          width: indicator.width,
          opacity: indicator.ready ? 1 : 0,
        }}
      />
      {options.map((o, i) => (
        <button
          key={o.value}
          type="button"
          ref={(el) => {
            btnRefs.current[i] = el;
          }}
          onClick={() => onChange(o.value as never)}
          className={`relative z-10 min-w-0 flex-1 rounded-[18px] px-1 py-2.5 text-center font-bold leading-tight transition-colors duration-300 ${textSize} ${
            value === o.value ? "text-white" : "text-black hover:text-[#ff7b00]"
          }`}
        >
          {o.label}
        </button>
      ))}
    </div>
  );
}
