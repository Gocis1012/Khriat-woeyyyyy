// apps/web/components/Navbar.tsx
import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="relative w-full h-16 flex items-center border-b px-4">
      {/* Logo */}
      <Link href="/" className="text-lg font-bold text-black hover:opacity-80">
        เครียดโว้ยยยย
      </Link>

      {/* Center Menu */}
      <div className="absolute left-1/2 -translate-x-1/2 flex gap-6">
        <Link href="/" className="hover:text-gray-400">
          หน้าแรก
        </Link>

        <Link href="/second" className="hover:text-gray-400">
          เกี่ยวกับเรา
        </Link>
      </div>

      {/* Right Menu */}
      <div className="ml-auto flex gap-4">
        <Link
          href="/login"
          className="px-4 py-2 border rounded hover:bg-gray-100"
        >
          เข้าสู่ระบบ
        </Link>
        <Link
          href="/signup"
          className="px-4 py-2 border rounded hover:bg-gray-100"
        >
          สมัครสมาชิก
        </Link>
      </div>
    </nav>
  );
}
