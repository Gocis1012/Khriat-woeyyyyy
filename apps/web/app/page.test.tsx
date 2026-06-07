import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

const { pushMock, refreshCreditMock, authState } = vi.hoisted(() => ({
  pushMock: vi.fn(),
  refreshCreditMock: vi.fn().mockResolvedValue(undefined),
  authState: { token: null as string | null },
}));

vi.mock("next/navigation", () => ({ useRouter: () => ({ push: pushMock }) }));
vi.mock("./context/GuestContext", () => ({
  useGuest: () => ({ credit: 6, refreshCredit: refreshCreditMock }),
}));
vi.mock("./context/AuthContext", () => ({ useAuth: () => authState }));

import Home from "./page";

const sendBtn = () => screen.getByRole("button", { name: "แปลงข้อความ" });
const textarea = () => screen.getByPlaceholderText(/พิมพ์เรื่องที่อยากบ่น/);

beforeEach(() => {
  vi.clearAllMocks();
  sessionStorage.clear();
  authState.token = null;
});

describe("Home page", () => {
  it("renders target and level selectors", () => {
    render(<Home />);
    expect(screen.getByText("หัวหน้า")).toBeInTheDocument();
    expect(screen.getByText("มนุษย์ปกติ")).toBeInTheDocument();
  });

  it("shows a validation error when input is empty", async () => {
    render(<Home />);
    fireEvent.click(sendBtn());
    expect(await screen.findByText(/พิมพ์เรื่องที่อยากบ่นก่อน/)).toBeInTheDocument();
  });

  it("translates and navigates to /second on success", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      status: 200,
      ok: true,
      json: async () => ({ result: "polite text", level: 3, target: "boss" }),
    });
    vi.stubGlobal("fetch", fetchMock);

    render(<Home />);
    fireEvent.change(textarea(), { target: { value: "angry rant" } });
    fireEvent.click(sendBtn());

    await waitFor(() => expect(pushMock).toHaveBeenCalledWith("/second"));
    expect(sessionStorage.getItem("translationResult")).toContain("polite text");
    expect(sessionStorage.getItem("translationInput")).toBe("angry rant");
    expect(refreshCreditMock).toHaveBeenCalled();
  });

  it("sends the Authorization header when logged in", async () => {
    authState.token = "jwt-abc";
    const fetchMock = vi.fn().mockResolvedValue({
      status: 200,
      ok: true,
      json: async () => ({ result: "ok" }),
    });
    vi.stubGlobal("fetch", fetchMock);

    render(<Home />);
    fireEvent.change(textarea(), { target: { value: "hi" } });
    fireEvent.click(sendBtn());

    await waitFor(() => expect(fetchMock).toHaveBeenCalled());
    const opts = fetchMock.mock.calls[0][1];
    expect(opts.headers.Authorization).toBe("Bearer jwt-abc");
  });

  it("shows an out-of-credit message on 402", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({ status: 402, ok: false, json: async () => ({}) })
    );
    render(<Home />);
    fireEvent.change(textarea(), { target: { value: "hi" } });
    fireEvent.click(sendBtn());
    expect(await screen.findByText(/เครดิตหมดแล้ว/)).toBeInTheDocument();
  });

  it("shows a generic error on a 500", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue({ status: 500, ok: false, json: async () => ({}) })
    );
    render(<Home />);
    fireEvent.change(textarea(), { target: { value: "hi" } });
    fireEvent.click(sendBtn());
    expect(await screen.findByText(/เกิดข้อผิดพลาด/)).toBeInTheDocument();
  });

  it("updates the feeling text when a different level is chosen", () => {
    render(<Home />);
    // pick level 5 by its label
    fireEvent.click(screen.getByText("ตัวแม่จะแคร์เพื่อ"));
    expect(screen.getByText(/Passive-Aggressive ขั้นเทพ/)).toBeInTheDocument();
  });
});
