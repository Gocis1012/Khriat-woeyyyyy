import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import PillBar from "./PillBar";

const opts = [
  { value: 1, label: "Alpha" },
  { value: 2, label: "Beta" },
  { value: 3, label: "Gamma" },
];

describe("PillBar", () => {
  it("renders all option labels", () => {
    render(<PillBar options={opts} value={1} onChange={() => {}} />);
    expect(screen.getByText("Alpha")).toBeInTheDocument();
    expect(screen.getByText("Beta")).toBeInTheDocument();
    expect(screen.getByText("Gamma")).toBeInTheDocument();
  });

  it("gives the selected option white text and others black", () => {
    render(<PillBar options={opts} value={2} onChange={() => {}} />);
    expect(screen.getByText("Beta").className).toContain("text-white");
    expect(screen.getByText("Alpha").className).toContain("text-black");
  });

  it("fires onChange with the clicked value", () => {
    const onChange = vi.fn();
    render(<PillBar options={opts} value={1} onChange={onChange} />);
    fireEvent.click(screen.getByText("Gamma"));
    expect(onChange).toHaveBeenCalledWith(3);
  });

  it("applies a custom text size class", () => {
    render(
      <PillBar options={opts} value={1} onChange={() => {}} textSize="text-xs" />
    );
    expect(screen.getByText("Alpha").className).toContain("text-xs");
  });
});
