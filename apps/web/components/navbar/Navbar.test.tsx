import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";

const { guest, auth, logoutMock } = vi.hoisted(() => ({
  guest: { value: { credit: 6 as number } },
  auth: {
    value: {
      user: null as null | { username: string; avatarUrl: string | null },
      isLoggedIn: false,
      login: vi.fn(),
      logout: vi.fn(),
    },
  },
  logoutMock: vi.fn(),
}));

vi.mock("next/navigation", () => ({ usePathname: () => "/" }));
vi.mock("../../app/context/GuestContext", () => ({ useGuest: () => guest.value }));
vi.mock("../../app/context/AuthContext", () => ({ useAuth: () => auth.value }));

import Navbar from "./์Narbar";

beforeEach(() => {
  guest.value = { credit: 6 };
  auth.value = { user: null, isLoggedIn: false, login: vi.fn(), logout: logoutMock };
  logoutMock.mockClear();
});

describe("Navbar", () => {
  it("shows the credit pill and a login button for guests", () => {
    render(<Navbar />);
    expect(screen.getByText(/6 ครั้ง/)).toBeInTheDocument();
    expect(screen.getByLabelText("เข้าสู่ระบบ")).toBeInTheDocument();
  });

  it("shows 'เครดิตหมด' when credit is 0", () => {
    guest.value = { credit: 0 };
    render(<Navbar />);
    expect(screen.getByText("เครดิตหมด")).toBeInTheDocument();
  });

  it("renders 'เกี่ยวกับเรา' as inert (not a link)", () => {
    render(<Navbar />);
    const about = screen.getByText("เกี่ยวกับเรา");
    expect(about.tagName.toLowerCase()).not.toBe("a");
  });

  it("shows the avatar and opens the logout modal when logged in", () => {
    auth.value = {
      user: { username: "Bob", avatarUrl: null },
      isLoggedIn: true,
      login: vi.fn(),
      logout: logoutMock,
    };
    render(<Navbar />);

    // avatar fallback initial
    const avatar = screen.getByLabelText("โปรไฟล์");
    expect(avatar).toBeInTheDocument();

    fireEvent.click(avatar);
    expect(screen.getByText("ออกจากระบบ?")).toBeInTheDocument();

    fireEvent.click(screen.getByText("ยืนยัน"));
    expect(logoutMock).toHaveBeenCalled();
  });
});
