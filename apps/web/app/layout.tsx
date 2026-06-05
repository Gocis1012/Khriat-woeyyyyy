import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import Script from "next/script";
import "./globals.css";
import Navbar from "../components/navbar/์Narbar";
import { iannnnnDog } from "../fonts";
import { AuthProvider } from "./context/AuthContext";
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
  description: "แปลงคำบ่นเป็นภาษาคนมีการศึกษา 😤",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="th" className={iannnnnDog.variable}>
      <head>
        <Script
          src="https://accounts.google.com/gsi/client"
          strategy="beforeInteractive"
        />
      </head>
      <body className="flex flex-col min-h-screen bg-[#111111] antialiased">
        <AuthProvider>
          <GuestProvider>
            <Navbar />
            <main className="flex-1 flex flex-col">{children}</main>
          </GuestProvider>
        </AuthProvider>
      </body>
    </html>
  );
}
