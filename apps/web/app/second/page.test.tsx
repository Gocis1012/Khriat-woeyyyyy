import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";

const { pushMock } = vi.hoisted(() => ({ pushMock: vi.fn() }));
vi.mock("next/navigation", () => ({ useRouter: () => ({ push: pushMock }) }));

import Second from "./page";

function seed() {
  sessionStorage.setItem(
    "translationResult",
    JSON.stringify({ result: "เรียบร้อยแล้วครับ" })
  );
  sessionStorage.setItem("translationInput", "งานช้าอีกแล้ว");
  sessionStorage.setItem("translationLevel", "2");
  sessionStorage.setItem("translationTarget", "boss");
}

beforeEach(() => {
  vi.clearAllMocks();
  sessionStorage.clear();
});

describe("Second (result) page", () => {
  it("renders the result, original input, target and level labels", async () => {
    seed();
    render(<Second />);
    expect(await screen.findByText("เรียบร้อยแล้วครับ")).toBeInTheDocument();
    expect(screen.getByText("งานช้าอีกแล้ว")).toBeInTheDocument();
    expect(screen.getByText("หัวหน้า")).toBeInTheDocument(); // target label
    expect(screen.getByText("พูดดีด้วยละนะ")).toBeInTheDocument(); // level 2 label
  });

  it("copies the result to the clipboard", async () => {
    seed();
    const writeText = vi.fn().mockResolvedValue(undefined);
    Object.assign(navigator, { clipboard: { writeText } });

    render(<Second />);
    await screen.findByText("เรียบร้อยแล้วครับ");
    fireEvent.click(screen.getByText("คัดลอก"));

    await waitFor(() => expect(writeText).toHaveBeenCalledWith("เรียบร้อยแล้วครับ"));
    expect(await screen.findByText(/คัดลอกแล้ว/)).toBeInTheDocument();
  });

  it("shows a not-found state when there is no result", async () => {
    render(<Second />);
    expect(await screen.findByText("ไม่พบผลลัพท์")).toBeInTheDocument();
  });

  it("navigates home from the back button", async () => {
    seed();
    render(<Second />);
    await screen.findByText("เรียบร้อยแล้วครับ");
    fireEvent.click(screen.getByText(/บ่นใหม่/));
    expect(pushMock).toHaveBeenCalledWith("/");
  });
});
