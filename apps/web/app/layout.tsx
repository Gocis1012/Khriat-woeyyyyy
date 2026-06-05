import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import Navbar from "../components/navbar/์Narbar";
import { iannnnnDog } from "../fonts";
import { GuestProvider } from "./context/GuestContext";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "เครียดโว้ยยยย",
  description: "Next.js + Tailwind v4",
};

// apps/web/app/layout.tsx
export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={iannnnnDog.variable}>
      {/* ปรับ body ให้เป็น flex คอลัมน์ เพื่อแยกสัดส่วนระหว่าง Navbar กับเนื้อหาให้ชัดเจน */}
      <body className="flex flex-col min-h-screen bg-[#FAF8F5] antialiased">
        <GuestProvider>
          <Navbar />
          <main className="flex-1 flex flex-col">
            {children}
          </main>
        </GuestProvider>
      </body>
    </html>
  );
}